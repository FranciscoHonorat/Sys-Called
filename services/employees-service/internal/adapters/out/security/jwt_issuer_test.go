package security_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/security"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type testClaims struct {
	Name               string `json:"name"`
	Role               string `json:"role"`
	MustChangePassword bool   `json:"must_change_password"`
	jwt.RegisteredClaims
}

func newKeyPair(t *testing.T) (ed25519.PublicKey, ed25519.PrivateKey) {
	t.Helper()
	public, private, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	return public, private
}

func parse(token string, key ed25519.PublicKey) (*testClaims, error) {
	claims := &testClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) { return key, nil },
		jwt.WithValidMethods([]string{"EdDSA"}),
		jwt.WithIssuer("employees-service"),
		jwt.WithAudience("sys-called"),
	)
	return claims, err
}

func TestJWTIssuer(t *testing.T) {
	admin := employee.NewEmployee("admin-1", "Administradora", "admin", employee.RoleAdmin, "hash")

	t.Run("should issue a token with the employee identity and role", func(t *testing.T) {
		public, private := newKeyPair(t)
		issuer := security.NewJWTIssuer(private, 15*time.Minute)

		token, expiresIn, err := issuer.Issue(admin)
		require.NoError(t, err)

		claims, err := parse(token, public)
		require.NoError(t, err)
		assert.Equal(t, "admin-1", claims.Subject)
		assert.Equal(t, "Administradora", claims.Name)
		assert.Equal(t, "admin", claims.Role)
		assert.Equal(t, 15*time.Minute, expiresIn)
		assert.WithinDuration(t, time.Now().Add(15*time.Minute), claims.ExpiresAt.Time, 2*time.Second)
	})

	t.Run("should tell when the employee must choose a new password", func(t *testing.T) {
		public, private := newKeyPair(t)
		issuer := security.NewJWTIssuer(private, 15*time.Minute)
		withTemporaryPassword := employee.NewEmployee("user-1", "Usuário Padrão", "usuario", employee.RoleUser, "hash")
		withTemporaryPassword.SetTemporaryPassword("temporary-hash")

		token, _, err := issuer.Issue(withTemporaryPassword)
		require.NoError(t, err)
		claims, err := parse(token, public)
		require.NoError(t, err)
		assert.True(t, claims.MustChangePassword)

		token, _, err = issuer.Issue(admin)
		require.NoError(t, err)
		claims, err = parse(token, public)
		require.NoError(t, err)
		assert.False(t, claims.MustChangePassword)
	})

	t.Run("should publish its public key and reference it in every token", func(t *testing.T) {
		public, private := newKeyPair(t)
		issuer := security.NewJWTIssuer(private, 15*time.Minute)

		keys := issuer.PublicKeys()
		require.Len(t, keys, 1)
		assert.Equal(t, "OKP", keys[0].KeyType)
		assert.Equal(t, "Ed25519", keys[0].Curve)
		assert.Equal(t, "EdDSA", keys[0].Algorithm)
		x, err := base64.RawURLEncoding.DecodeString(keys[0].X)
		require.NoError(t, err)
		assert.Equal(t, []byte(public), x)
		assert.NotEmpty(t, keys[0].ID)

		token, _, err := issuer.Issue(admin)
		require.NoError(t, err)
		parsed, _, err := jwt.NewParser().ParseUnverified(token, &testClaims{})
		require.NoError(t, err)
		assert.Equal(t, keys[0].ID, parsed.Header["kid"])
	})

	t.Run("should not be accepted by a different key", func(t *testing.T) {
		_, private := newKeyPair(t)
		otherPublic, _ := newKeyPair(t)
		issuer := security.NewJWTIssuer(private, 15*time.Minute)

		token, _, err := issuer.Issue(admin)
		require.NoError(t, err)

		_, err = parse(token, otherPublic)
		assert.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
	})
}
