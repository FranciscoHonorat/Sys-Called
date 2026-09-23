package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type RegisterEmployeeInput struct {
	ID       string
	Name     string
	Username string
	Password string
	Role     string
}

type RegisterEmployeeUseCase struct {
	repo   out.EmployeeRepository
	hasher out.PasswordHasher
}

func NewRegisterEmployeeUseCase(repo out.EmployeeRepository, hasher out.PasswordHasher) *RegisterEmployeeUseCase {
	return &RegisterEmployeeUseCase{repo: repo, hasher: hasher}
}

func (uc *RegisterEmployeeUseCase) Execute(ctx context.Context, input RegisterEmployeeInput) error {
	if input.Password == "" {
		return domainErr.ErrInvalidPassword
	}
	role, err := employee.NewRole(input.Role)
	if err != nil {
		return err
	}

	passwordHash, err := uc.hasher.Hash(input.Password)
	if err != nil {
		return err
	}

	e, err := employee.Register(input.ID, input.Name, input.Username, role, passwordHash)
	if err != nil {
		return err
	}

	return uc.repo.Register(ctx, e)
}
