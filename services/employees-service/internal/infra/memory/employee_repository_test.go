package memory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/infra/memory"
)

func TestEmployeeRepository(t *testing.T) {
	t.Run("should list at least 3 seeded employees", func(t *testing.T) {
		repo := memory.NewEmployeeRepository()

		employees, err := repo.List(context.Background())

		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(employees), 3)
	})

	t.Run("should return employees with non-empty ID and name", func(t *testing.T) {
		repo := memory.NewEmployeeRepository()

		employees, err := repo.List(context.Background())
		require.NoError(t, err)

		for _, e := range employees {
			assert.NotEmpty(t, e.GetID())
			assert.NotEmpty(t, e.GetName())
		}
	})
}
