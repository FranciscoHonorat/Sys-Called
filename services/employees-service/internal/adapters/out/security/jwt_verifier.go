package security

import (
	"crypto/ed25519"

	"github.com/golang-jwt/jwt/v5"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

var _ out.TokenVerifier = (*JWTVerifier)(nil)

type JWTVerifier struct {
	key    ed25519.PublicKey
	parser *jwt.Parser
}

func NewJWTVerifier(key ed25519.PublicKey) *JWTVerifier {
	return &JWTVerifier{
		key: key,
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
			jwt.WithIssuer(TokenIssuer),
			jwt.WithAudience(TokenAudience),
			jwt.WithExpirationRequired(),
		),
	}
}

func (v *JWTVerifier) Verify(token string) (out.Caller, error) {
	parsed := &claims{}
	if _, err := v.parser.ParseWithClaims(token, parsed, func(*jwt.Token) (any, error) { return v.key, nil }); err != nil {
		return out.Caller{}, err
	}
	role, err := employee.NewRole(parsed.Role)
	if err != nil {
		return out.Caller{}, err
	}
	return out.Caller{ID: parsed.Subject, Role: role}, nil
}
