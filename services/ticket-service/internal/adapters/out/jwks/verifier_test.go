package jwks_test

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/jwks"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

type signingKey struct {
	id      string
	private ed25519.PrivateKey
	public  ed25519.PublicKey
}

func newSigningKey(t *testing.T, id string) signingKey {
	t.Helper()
	public, private, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	return signingKey{id: id, private: private, public: public}
}

type issuerServer struct {
	*httptest.Server
	keys    []signingKey
	fetches int
	failing bool
	hold    chan struct{}
	held    chan struct{}
}

func newIssuerServer(t *testing.T, keys ...signingKey) *issuerServer {
	t.Helper()
	s := &issuerServer{keys: keys}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		s.fetches++
		if s.hold != nil {
			close(s.held)
			<-s.hold
		}
		if s.failing {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal error"})
			return
		}
		published := []map[string]string{}
		for _, k := range s.keys {
			published = append(published, map[string]string{
				"kid": k.id, "kty": "OKP", "crv": "Ed25519", "alg": "EdDSA", "use": "sig",
				"x": base64.RawURLEncoding.EncodeToString(k.public),
			})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": published})
	}))
	t.Cleanup(s.Close)
	return s
}

func (s *issuerServer) jwksURL() string {
	return s.URL + "/.well-known/jwks.json"
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"sub":  "admin-1",
		"name": "Administradora",
		"role": "admin",
		"iss":  "employees-service",
		"aud":  []string{"sys-called"},
		"exp":  time.Now().Add(time.Minute).Unix(),
	}
}

func sign(t *testing.T, key signingKey, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = key.id
	signed, err := token.SignedString(key.private)
	require.NoError(t, err)
	return signed
}

func TestVerifier(t *testing.T) {
	t.Run("should turn a valid token into the actor it identifies", func(t *testing.T) {
		key := newSigningKey(t, "key-1")
		server := newIssuerServer(t, key)
		verifier := jwks.NewVerifier(server.jwksURL(), http.DefaultClient)

		got, err := verifier.Verify(context.Background(), sign(t, key, validClaims()))

		require.NoError(t, err)
		assert.Equal(t, "admin-1", got.ID())
		assert.Equal(t, "Administradora", got.Name())
		assert.Equal(t, actor.RoleAdmin, got.Role())
	})

	with := func(key string, value any) jwt.MapClaims {
		c := validClaims()
		if value == nil {
			delete(c, key)
		} else {
			c[key] = value
		}
		return c
	}

	rejected := []struct {
		name  string
		token func(t *testing.T, published signingKey) string
	}{
		{"a token signed by another key with the same kid", func(t *testing.T, published signingKey) string {
			forger := newSigningKey(t, published.id)
			return sign(t, forger, validClaims())
		}},
		{"a token with an unknown kid", func(t *testing.T, _ signingKey) string {
			return sign(t, newSigningKey(t, "unknown"), validClaims())
		}},
		{"an expired token", func(t *testing.T, published signingKey) string {
			return sign(t, published, with("exp", time.Now().Add(-time.Minute).Unix()))
		}},
		{"a token without expiration", func(t *testing.T, published signingKey) string {
			return sign(t, published, with("exp", nil))
		}},
		{"a token from another issuer", func(t *testing.T, published signingKey) string {
			return sign(t, published, with("iss", "someone-else"))
		}},
		{"a token for another audience", func(t *testing.T, published signingKey) string {
			return sign(t, published, with("aud", []string{"another-system"}))
		}},
		{"a token with an unknown role", func(t *testing.T, published signingKey) string {
			return sign(t, published, with("role", "root"))
		}},
		{"an HMAC token trying to use the public key as secret", func(t *testing.T, published signingKey) string {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims())
			token.Header["kid"] = published.id
			signed, err := token.SignedString([]byte(published.public))
			require.NoError(t, err)
			return signed
		}},
	}
	for _, tc := range rejected {
		t.Run("should reject "+tc.name, func(t *testing.T) {
			key := newSigningKey(t, "key-1")
			server := newIssuerServer(t, key)
			verifier := jwks.NewVerifier(server.jwksURL(), http.DefaultClient)

			_, err := verifier.Verify(context.Background(), tc.token(t, key))

			assert.Error(t, err)
		})
	}
}

func TestVerifierKeyCache(t *testing.T) {
	t.Run("should download the key set once and reuse it", func(t *testing.T) {
		key := newSigningKey(t, "key-1")
		server := newIssuerServer(t, key)
		verifier := jwks.NewVerifier(server.jwksURL(), http.DefaultClient)

		for range 3 {
			_, err := verifier.Verify(context.Background(), sign(t, key, validClaims()))
			require.NoError(t, err)
		}

		assert.Equal(t, 1, server.fetches)
	})

	t.Run("should download the key set again when the issuer rotates its key", func(t *testing.T) {
		oldKey, newKey := newSigningKey(t, "key-1"), newSigningKey(t, "key-2")
		server := newIssuerServer(t, oldKey)
		now := time.Now()
		verifier := jwks.NewVerifier(server.jwksURL(), http.DefaultClient, jwks.WithClock(func() time.Time { return now }))
		_, err := verifier.Verify(context.Background(), sign(t, oldKey, validClaims()))
		require.NoError(t, err)

		server.keys = []signingKey{newKey}
		now = now.Add(jwks.MinRefreshInterval)
		got, err := verifier.Verify(context.Background(), sign(t, newKey, validClaims()))

		require.NoError(t, err)
		assert.Equal(t, "admin-1", got.ID())
		assert.Equal(t, 2, server.fetches)
	})

	t.Run("should not let unknown kids trigger a download on every request", func(t *testing.T) {
		key := newSigningKey(t, "key-1")
		server := newIssuerServer(t, key)
		verifier := jwks.NewVerifier(server.jwksURL(), http.DefaultClient)
		_, err := verifier.Verify(context.Background(), sign(t, key, validClaims()))
		require.NoError(t, err)

		for range 5 {
			_, err := verifier.Verify(context.Background(), sign(t, newSigningKey(t, "random"), validClaims()))
			assert.Error(t, err)
		}

		assert.Equal(t, 1, server.fetches)
	})

	t.Run("should keep the known keys when the issuer answers with an error", func(t *testing.T) {
		key := newSigningKey(t, "key-1")
		server := newIssuerServer(t, key)
		now := time.Now()
		verifier := jwks.NewVerifier(server.jwksURL(), http.DefaultClient, jwks.WithClock(func() time.Time { return now }))
		_, err := verifier.Verify(context.Background(), sign(t, key, validClaims()))
		require.NoError(t, err)

		server.failing = true
		now = now.Add(jwks.MinRefreshInterval)
		_, err = verifier.Verify(context.Background(), sign(t, newSigningKey(t, "random"), validClaims()))
		require.Error(t, err)

		_, err = verifier.Verify(context.Background(), sign(t, key, validClaims()))
		assert.NoError(t, err)
	})

	t.Run("should verify tokens with a known key while a download is in progress", func(t *testing.T) {
		key := newSigningKey(t, "key-1")
		server := newIssuerServer(t, key)
		now := time.Now()
		verifier := jwks.NewVerifier(server.jwksURL(), http.DefaultClient, jwks.WithClock(func() time.Time { return now }))
		_, err := verifier.Verify(context.Background(), sign(t, key, validClaims()))
		require.NoError(t, err)

		server.hold, server.held = make(chan struct{}), make(chan struct{})
		now = now.Add(jwks.MinRefreshInterval)
		downloading := make(chan struct{})
		go func() {
			defer close(downloading)
			_, _ = verifier.Verify(context.Background(), sign(t, newSigningKey(t, "random"), validClaims()))
		}()
		<-server.held

		verified := make(chan error, 1)
		go func() {
			_, err := verifier.Verify(context.Background(), sign(t, key, validClaims()))
			verified <- err
		}()

		select {
		case err := <-verified:
			assert.NoError(t, err)
		case <-time.After(time.Second):
			t.Error("verification waited for the key download")
		}
		close(server.hold)
		<-downloading
	})
}
