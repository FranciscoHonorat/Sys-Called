package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

const refreshTokenBytes = 32

var _ out.RefreshTokenGenerator = RefreshTokenGenerator{}

type RefreshTokenGenerator struct{}

func NewRefreshTokenGenerator() RefreshTokenGenerator {
	return RefreshTokenGenerator{}
}

func (RefreshTokenGenerator) Generate() (string, error) {
	raw := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func (RefreshTokenGenerator) Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
