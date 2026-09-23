package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
)

type LoginInput struct {
	Username string
	Password string
}

type LoginUseCase struct {
	repo      out.EmployeeRepository
	hasher    out.PasswordHasher
	sessions  *Sessions
	dummyHash string
}

func NewLoginUseCase(repo out.EmployeeRepository, hasher out.PasswordHasher, sessions *Sessions) *LoginUseCase {
	dummyHash, _ := hasher.Hash("dummy-password-for-unknown-users")
	return &LoginUseCase{repo: repo, hasher: hasher, sessions: sessions, dummyHash: dummyHash}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (SessionOutput, error) {
	e, found, err := uc.repo.FindByUsername(ctx, input.Username)
	if err != nil {
		return SessionOutput{}, err
	}
	hash := uc.dummyHash
	if found {
		hash = e.GetPasswordHash()
	}
	if matches := uc.hasher.Matches(hash, input.Password); !found || !matches {
		return SessionOutput{}, domainErr.ErrInvalidCredentials
	}
	if !e.CanLogIn() {
		return SessionOutput{}, domainErr.ErrPendingApproval
	}

	return uc.sessions.Start(ctx, e)
}
