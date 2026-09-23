package security

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

var _ out.PasswordHasher = BcryptHasher{}

type BcryptHasher struct{}

func NewBcryptHasher() BcryptHasher {
	return BcryptHasher{}
}

func (BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (BcryptHasher) Matches(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
