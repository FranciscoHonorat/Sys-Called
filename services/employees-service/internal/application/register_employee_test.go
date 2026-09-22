package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestRegisterEmployeeUseCase(t *testing.T) {
	t.Run("should register a new employee", func(t *testing.T) {
		repo := &fakeEmployeeRepository{}
		uc := application.NewRegisterEmployeeUseCase(repo)

		err := uc.Execute(context.Background(), application.RegisterEmployeeInput{
			ID:   "agent-1",
			Name: "Ana Souza",
		})

		require.NoError(t, err)
		require.Len(t, repo.employees, 1)
		assert.Equal(t, "agent-1", repo.employees[0].GetID())
		assert.Equal(t, "Ana Souza", repo.employees[0].GetName())
		assert.Equal(t, []employee.Event{employee.Registered{ID: "agent-1", Name: "Ana Souza"}}, repo.employees[0].Events())
	})

	t.Run("should return an error for an empty ID", func(t *testing.T) {
		repo := &fakeEmployeeRepository{}
		uc := application.NewRegisterEmployeeUseCase(repo)

		err := uc.Execute(context.Background(), application.RegisterEmployeeInput{
			ID:   "",
			Name: "Ana Souza",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidEmployeeID)
	})

	t.Run("should return an error for an empty name", func(t *testing.T) {
		repo := &fakeEmployeeRepository{}
		uc := application.NewRegisterEmployeeUseCase(repo)

		err := uc.Execute(context.Background(), application.RegisterEmployeeInput{
			ID:   "agent-1",
			Name: "",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidEmployeeName)
	})

	t.Run("should propagate repository errors", func(t *testing.T) {
		repo := &fakeEmployeeRepository{registerErr: assert.AnError}
		uc := application.NewRegisterEmployeeUseCase(repo)

		err := uc.Execute(context.Background(), application.RegisterEmployeeInput{
			ID:   "agent-1",
			Name: "Ana Souza",
		})

		assert.ErrorIs(t, err, assert.AnError)
	})
}
