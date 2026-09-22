package employee_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestEmployee(t *testing.T) {
	t.Run("should record a Registered event when registering", func(t *testing.T) {
		e := employee.Register("agent-1", "Ana Souza")

		assert.Equal(t, []employee.Event{employee.Registered{ID: "agent-1", Name: "Ana Souza"}}, e.Events())
		assert.Equal(t, "EmployeeRegistered", e.Events()[0].EventType())
	})

	t.Run("should not record events when rebuilding an existing employee", func(t *testing.T) {
		e := employee.NewEmployee("agent-1", "Ana Souza")

		assert.Empty(t, e.Events())
	})
}
