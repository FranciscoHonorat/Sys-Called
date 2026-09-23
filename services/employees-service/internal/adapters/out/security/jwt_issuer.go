package security

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

const (
	TokenIssuer   = "employees-service"
	TokenAudience = "sys-called"
)

var (
	_ out.TokenIssuer = (*JWTIssuer)(nil)
	_ out.SigningKeys = (*JWTIssuer)(nil)
)

type claims struct {
	Name               string `json:"name"`
	Role               string `json:"role"`
	MustChangePassword bool   `json:"must_change_password"`
	jwt.RegisteredClaims
}

type JWTIssuer struct {
	key       ed25519.PrivateKey
	ttl       time.Duration
	publicKey out.PublicKey
}

func NewJWTIssuer(key ed25519.PrivateKey, ttl time.Duration) *JWTIssuer {
	return &JWTIssuer{key: key, ttl: ttl, publicKey: publicJWK(key.Public().(ed25519.PublicKey))}
}

func publicJWK(key ed25519.PublicKey) out.PublicKey {
	x := base64.RawURLEncoding.EncodeToString(key)
	thumbprint := sha256.Sum256([]byte(`{"crv":"Ed25519","kty":"OKP","x":"` + x + `"}`))
	return out.PublicKey{
		ID:        base64.RawURLEncoding.EncodeToString(thumbprint[:]),
		KeyType:   "OKP",
		Curve:     "Ed25519",
		Algorithm: "EdDSA",
		X:         x,
	}
}

func (i *JWTIssuer) PublicKeys() []out.PublicKey {
	return []out.PublicKey{i.publicKey}
}

func (i *JWTIssuer) Issue(e employee.Employee) (string, time.Duration, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims{
		Name:               e.GetName(),
		Role:               e.GetRole().String(),
		MustChangePassword: e.MustChangePassword(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   e.GetID(),
			Issuer:    TokenIssuer,
			Audience:  jwt.ClaimStrings{TokenAudience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(i.ttl)),
		},
	})
	token.Header["kid"] = i.publicKey.ID

	signed, err := token.SignedString(i.key)
	if err != nil {
		return "", 0, err
	}
	return signed, i.ttl, nil
}
