package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func validRegisterInput() application.RegisterEmployeeInput {
	return application.RegisterEmployeeInput{
		ID:       "agent-1",
		Name:     "Ana Souza",
		Username: "ana",
		Password: "secret",
		Role:     "support",
	}
}

func TestRegisterEmployeeUseCase(t *testing.T) {
	t.Run("should register a new employee with a hashed password", func(t *testing.T) {
		repo := &outtest.EmployeeRepository{}
		uc := application.NewRegisterEmployeeUseCase(repo, &outtest.PasswordHasher{})

		err := uc.Execute(context.Background(), validRegisterInput())

		require.NoError(t, err)
		require.Len(t, repo.Employees, 1)
		registered := repo.Employees[0]
		assert.Equal(t, "agent-1", registered.GetID())
		assert.Equal(t, "Ana Souza", registered.GetName())
		assert.Equal(t, "ana", registered.GetUsername())
		assert.Equal(t, employee.RoleSupport, registered.GetRole())
		assert.Equal(t, "hashed:secret", registered.GetPasswordHash())
		assert.Equal(t, []employee.Event{employee.Registered{ID: "agent-1", Name: "Ana Souza", Role: "support"}}, registered.Events())
	})

	invalidInputs := []struct {
		name    string
		mutate  func(*application.RegisterEmployeeInput)
		wantErr error
	}{
		{"empty ID", func(in *application.RegisterEmployeeInput) { in.ID = "" }, domainErr.ErrInvalidEmployeeID},
		{"empty name", func(in *application.RegisterEmployeeInput) { in.Name = "" }, domainErr.ErrInvalidEmployeeName},
		{"empty username", func(in *application.RegisterEmployeeInput) { in.Username = "" }, domainErr.ErrInvalidUsername},
		{"empty password", func(in *application.RegisterEmployeeInput) { in.Password = "" }, domainErr.ErrInvalidPassword},
		{"unknown role", func(in *application.RegisterEmployeeInput) { in.Role = "root" }, domainErr.ErrInvalidRole},
	}
	for _, tc := range invalidInputs {
		t.Run("should reject an "+tc.name, func(t *testing.T) {
			repo := &outtest.EmployeeRepository{}
			uc := application.NewRegisterEmployeeUseCase(repo, &outtest.PasswordHasher{})
			input := validRegisterInput()
			tc.mutate(&input)

			err := uc.Execute(context.Background(), input)

			assert.ErrorIs(t, err, tc.wantErr)
			assert.Empty(t, repo.Employees)
		})
	}

	t.Run("should propagate repository errors", func(t *testing.T) {
		repo := &outtest.EmployeeRepository{RegisterErr: assert.AnError}
		uc := application.NewRegisterEmployeeUseCase(repo, &outtest.PasswordHasher{})

		err := uc.Execute(context.Background(), validRegisterInput())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
