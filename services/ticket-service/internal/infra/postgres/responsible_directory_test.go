package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/postgres"
)

func TestResponsibleDirectory(t *testing.T) {
	pool := testPool(t)
	directory := postgres.NewResponsibleDirectory(pool)

	t.Run("should upsert and list a responsible", func(t *testing.T) {
		id := "test-responsible-" + uuid.New().String()

		require.NoError(t, directory.Upsert(context.Background(), id, "Test Responsible"))

		ids, err := directory.List(context.Background())
		require.NoError(t, err)
		assert.Contains(t, ids, id)
	})

	t.Run("should update the name without duplicating the entry on repeated upsert", func(t *testing.T) {
		id := "test-responsible-" + uuid.New().String()

		require.NoError(t, directory.Upsert(context.Background(), id, "First Name"))
		require.NoError(t, directory.Upsert(context.Background(), id, "Second Name"))

		ids, err := directory.List(context.Background())
		require.NoError(t, err)

		count := 0
		for _, got := range ids {
			if got == id {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})
}
