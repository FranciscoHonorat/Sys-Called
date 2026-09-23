package application

import (
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
)

type AuthenticateUseCase struct {
	verifier out.TokenVerifier
}

func NewAuthenticateUseCase(verifier out.TokenVerifier) *AuthenticateUseCase {
	return &AuthenticateUseCase{verifier: verifier}
}

func (uc *AuthenticateUseCase) Execute(token string) (out.Caller, error) {
	if token == "" {
		return out.Caller{}, domainErr.ErrUnauthenticated
	}
	caller, err := uc.verifier.Verify(token)
	if err != nil {
		return out.Caller{}, domainErr.ErrUnauthenticated
	}
	return caller, nil
}
