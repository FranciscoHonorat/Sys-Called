package employee_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestRole(t *testing.T) {
	t.Run("should accept the known roles", func(t *testing.T) {
		for _, raw := range []string{"user", "support", "admin"} {
			role, err := employee.NewRole(raw)

			require.NoError(t, err)
			assert.Equal(t, raw, role.String())
		}
	})

	t.Run("should reject an unknown role", func(t *testing.T) {
		_, err := employee.NewRole("root")

		assert.ErrorIs(t, err, domainErr.ErrInvalidRole)
	})
}
