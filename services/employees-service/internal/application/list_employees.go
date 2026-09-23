package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type EmployeeOutput struct {
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	Username               string `json:"username"`
	Role                   string `json:"role"`
	Status                 string `json:"status"`
	PasswordResetRequested bool   `json:"password_reset_requested"`
}

type ListEmployeesUseCase struct {
	repo out.EmployeeRepository
}

func NewListEmployeesUseCase(repo out.EmployeeRepository) *ListEmployeesUseCase {
	return &ListEmployeesUseCase{repo: repo}
}

func (uc *ListEmployeesUseCase) Execute(ctx context.Context, caller out.Caller) ([]EmployeeOutput, error) {
	if caller.Role != employee.RoleAdmin {
		return nil, domainErr.ErrForbidden
	}

	employees, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]EmployeeOutput, 0, len(employees))
	for _, e := range employees {
		outputs = append(outputs, EmployeeOutput{
			ID:                     e.GetID(),
			Name:                   e.GetName(),
			Username:               e.GetUsername(),
			Role:                   e.GetRole().String(),
			Status:                 string(e.GetStatus()),
			PasswordResetRequested: e.HasRequestedPasswordReset(),
		})
	}

	return outputs, nil
}
