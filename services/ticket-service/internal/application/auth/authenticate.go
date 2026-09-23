package auth

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type AuthenticateUseCase struct {
	verifier out.TokenVerifier
}

func NewAuthenticateUseCase(verifier out.TokenVerifier) *AuthenticateUseCase {
	return &AuthenticateUseCase{verifier: verifier}
}

func (uc *AuthenticateUseCase) Execute(ctx context.Context, token string) (actor.Actor, error) {
	if token == "" {
		return actor.Actor{}, domainErr.ErrUnauthenticated
	}

	a, err := uc.verifier.Verify(ctx, token)
	if err != nil {
		return actor.Actor{}, domainErr.ErrUnauthenticated
	}
	return a, nil
}
