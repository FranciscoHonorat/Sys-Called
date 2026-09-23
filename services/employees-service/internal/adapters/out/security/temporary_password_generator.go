package security

import (
	"crypto/rand"
	"math/big"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

const (
	temporaryPasswordLength   = 10
	temporaryPasswordAlphabet = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
)

var _ out.PasswordGenerator = TemporaryPasswordGenerator{}

type TemporaryPasswordGenerator struct{}

func NewTemporaryPasswordGenerator() TemporaryPasswordGenerator {
	return TemporaryPasswordGenerator{}
}

func (TemporaryPasswordGenerator) Generate() (string, error) {
	limit := big.NewInt(int64(len(temporaryPasswordAlphabet)))
	password := make([]byte, temporaryPasswordLength)
	for i := range password {
		n, err := rand.Int(rand.Reader, limit)
		if err != nil {
			return "", err
		}
		password[i] = temporaryPasswordAlphabet[n.Int64()]
	}
	return string(password), nil
}
