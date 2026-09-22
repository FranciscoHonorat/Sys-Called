package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type RegisterEmployeeInput struct {
	ID   string
	Name string
}

type RegisterEmployeeUseCase struct {
	repo out.EmployeeRepository
}

func NewRegisterEmployeeUseCase(repo out.EmployeeRepository) *RegisterEmployeeUseCase {
	return &RegisterEmployeeUseCase{repo: repo}
}

func (uc *RegisterEmployeeUseCase) Execute(ctx context.Context, input RegisterEmployeeInput) error {
	if input.ID == "" {
		return domainErr.ErrInvalidEmployeeID
	}
	if input.Name == "" {
		return domainErr.ErrInvalidEmployeeName
	}

	return uc.repo.Register(ctx, employee.Register(input.ID, input.Name))
}
