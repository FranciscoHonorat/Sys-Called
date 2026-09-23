package employee_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestEmployee(t *testing.T) {
	t.Run("should keep username, role and password hash", func(t *testing.T) {
		e, err := employee.Register("agent-1", "Ana Souza", "ana", employee.RoleSupport, "hashed")
		require.NoError(t, err)

		assert.Equal(t, "ana", e.GetUsername())
		assert.Equal(t, employee.RoleSupport, e.GetRole())
		assert.Equal(t, "hashed", e.GetPasswordHash())
	})

	t.Run("should record a Registered event with the role but never the password", func(t *testing.T) {
		e, err := employee.Register("agent-1", "Ana Souza", "ana", employee.RoleSupport, "hashed")
		require.NoError(t, err)

		assert.Equal(t, []employee.Event{employee.Registered{ID: "agent-1", Name: "Ana Souza", Role: "support"}}, e.Events())
		assert.Equal(t, "EmployeeRegistered", e.Events()[0].EventType())
	})

	t.Run("should not record events when rebuilding an existing employee", func(t *testing.T) {
		e := employee.NewEmployee("agent-1", "Ana Souza", "ana", employee.RoleSupport, "hashed")

		assert.Empty(t, e.Events())
	})
}
