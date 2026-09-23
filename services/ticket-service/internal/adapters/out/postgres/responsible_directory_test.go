package postgres_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
)

func newTestResponsibleID(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	id := "test-responsible-" + uuid.New().String()
	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), `DELETE FROM responsibles WHERE id = $1`, id)
	})
	return id
}

func TestResponsibleDirectory(t *testing.T) {
	pool := testPool(t)
	directory := postgres.NewResponsibleDirectory(pool)

	t.Run("should list responsibles with their current names", func(t *testing.T) {
		id := newTestResponsibleID(t, pool)
		require.NoError(t, directory.Upsert(context.Background(), id, "First Name"))
		require.NoError(t, directory.Upsert(context.Background(), id, "Second Name"))

		responsibles, err := directory.ListWithNames(context.Background())

		require.NoError(t, err)
		assert.Contains(t, responsibles, out.Responsible{ID: id, Name: "Second Name"})
	})

	t.Run("should upsert and list a responsible", func(t *testing.T) {
		id := newTestResponsibleID(t, pool)

		require.NoError(t, directory.Upsert(context.Background(), id, "Test Responsible"))

		ids, err := directory.List(context.Background())
		require.NoError(t, err)
		assert.Contains(t, ids, id)
	})

	t.Run("should update the name without duplicating the entry on repeated upsert", func(t *testing.T) {
		id := newTestResponsibleID(t, pool)

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
