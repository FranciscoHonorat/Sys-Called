package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type fakeEmployeeRepository struct {
	employees   []employee.Employee
	err         error
	registerErr error
}

func (r *fakeEmployeeRepository) List(_ context.Context) ([]employee.Employee, error) {
	return r.employees, r.err
}

func (r *fakeEmployeeRepository) Register(_ context.Context, e employee.Employee) error {
	if r.registerErr != nil {
		return r.registerErr
	}
	r.employees = append(r.employees, e)
	return nil
}

func TestListEmployeesUseCase(t *testing.T) {
	t.Run("should map employees to output", func(t *testing.T) {
		repo := &fakeEmployeeRepository{employees: []employee.Employee{
			employee.NewEmployee("agent-1", "Ana Souza"),
			employee.NewEmployee("agent-2", "Bruno Lima"),
		}}
		uc := application.NewListEmployeesUseCase(repo)

		output, err := uc.Execute(context.Background())

		require.NoError(t, err)
		require.Len(t, output, 2)
		assert.Equal(t, "agent-1", output[0].ID)
		assert.Equal(t, "Ana Souza", output[0].Name)
	})

	t.Run("should propagate repository errors", func(t *testing.T) {
		repo := &fakeEmployeeRepository{err: assert.AnError}
		uc := application.NewListEmployeesUseCase(repo)

		_, err := uc.Execute(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
