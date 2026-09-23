package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

var adminCaller = out.Caller{ID: "admin-1", Role: employee.RoleAdmin}

func TestListEmployeesUseCase(t *testing.T) {
	t.Run("gives admins every employee with username and role", func(t *testing.T) {
		repo := &outtest.EmployeeRepository{Employees: []employee.Employee{
			employee.NewEmployee("agent-1", "Ana Souza", "ana", employee.RoleSupport, ""),
			employee.NewEmployee("user-1", "Usuário Padrão", "usuario", employee.RoleUser, ""),
		}}
		uc := application.NewListEmployeesUseCase(repo)

		output, err := uc.Execute(context.Background(), adminCaller)

		require.NoError(t, err)
		assert.Equal(t, []application.EmployeeOutput{
			{ID: "agent-1", Name: "Ana Souza", Username: "ana", Role: "support", Status: "active"},
			{ID: "user-1", Name: "Usuário Padrão", Username: "usuario", Role: "user", Status: "active"},
		}, output)
	})

	t.Run("is only for admins", func(t *testing.T) {
		uc := application.NewListEmployeesUseCase(&outtest.EmployeeRepository{})

		_, err := uc.Execute(context.Background(), out.Caller{ID: "agent-1", Role: employee.RoleSupport})

		assert.ErrorIs(t, err, domainErr.ErrForbidden)
	})

	t.Run("should propagate repository errors", func(t *testing.T) {
		uc := application.NewListEmployeesUseCase(&outtest.EmployeeRepository{Err: assert.AnError})

		_, err := uc.Execute(context.Background(), adminCaller)

		assert.ErrorIs(t, err, assert.AnError)
	})
}
