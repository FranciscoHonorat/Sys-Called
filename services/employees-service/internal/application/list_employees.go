package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/repository"
)

type EmployeeOutput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListEmployeesUseCase struct {
	repo repository.EmployeeRepository
}

func NewListEmployeesUseCase(repo repository.EmployeeRepository) *ListEmployeesUseCase {
	return &ListEmployeesUseCase{repo: repo}
}

func (uc *ListEmployeesUseCase) Execute(ctx context.Context) ([]EmployeeOutput, error) {
	employees, err := uc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]EmployeeOutput, 0, len(employees))
	for _, e := range employees {
		outputs = append(outputs, EmployeeOutput{ID: e.GetID(), Name: e.GetName()})
	}

	return outputs, nil
}
