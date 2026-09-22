package application

import (
	"context"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/repository"
)

type RegisterEmployeeInput struct {
	ID   string
	Name string
}

type RegisterEmployeeUseCase struct {
	repo repository.EmployeeRepository
}

func NewRegisterEmployeeUseCase(repo repository.EmployeeRepository) *RegisterEmployeeUseCase {
	return &RegisterEmployeeUseCase{repo: repo}
}

func (uc *RegisterEmployeeUseCase) Execute(ctx context.Context, input RegisterEmployeeInput) error {
	if input.ID == "" {
		return domainErr.ErrInvalidEmployeeID
	}
	if input.Name == "" {
		return domainErr.ErrInvalidEmployeeName
	}

	return uc.repo.Register(ctx, employee.NewEmployee(input.ID, input.Name))
}
