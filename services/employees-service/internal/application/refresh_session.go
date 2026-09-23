package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
)

type RefreshSessionUseCase struct {
	repo     out.EmployeeRepository
	sessions *Sessions
}

func NewRefreshSessionUseCase(repo out.EmployeeRepository, sessions *Sessions) *RefreshSessionUseCase {
	return &RefreshSessionUseCase{repo: repo, sessions: sessions}
}

func (uc *RefreshSessionUseCase) Execute(ctx context.Context, refreshToken string) (SessionOutput, error) {
	employeeID, err := uc.sessions.Consume(ctx, refreshToken)
	if err != nil {
		return SessionOutput{}, err
	}

	e, found, err := uc.repo.FindByID(ctx, employeeID)
	if err != nil {
		return SessionOutput{}, err
	}
	if !found {
		return SessionOutput{}, domainErr.ErrInvalidRefreshToken
	}

	return uc.sessions.Start(ctx, e)
}
