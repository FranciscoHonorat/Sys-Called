package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/infra/postgres"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("EMPLOYEES_SERVICE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("EMPLOYEES_SERVICE_TEST_DATABASE_URL not set; skipping Postgres integration tests")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool
}

func TestEmployeeRepository(t *testing.T) {
	pool := testPool(t)
	repo := postgres.NewEmployeeRepository(pool)
	outboxStore := postgres.NewOutboxStore(pool)

	t.Run("should register an employee and list it back", func(t *testing.T) {
		e := employee.NewEmployee("test-agent-1", "Test Agent One")

		err := repo.Register(context.Background(), e)
		require.NoError(t, err)

		employees, err := repo.List(context.Background())
		require.NoError(t, err)

		found := false
		for _, got := range employees {
			if got.GetID() == "test-agent-1" {
				found = true
				assert.Equal(t, "Test Agent One", got.GetName())
			}
		}
		assert.True(t, found)
	})

	t.Run("should write an outbox event when registering a new employee", func(t *testing.T) {
		e := employee.NewEmployee("test-agent-2", "Test Agent Two")

		require.NoError(t, repo.Register(context.Background(), e))

		pending, err := outboxStore.FetchPending(context.Background())
		require.NoError(t, err)

		found := false
		for _, evt := range pending {
			if evt.EventType == "EmployeeRegistered" {
				found = true
			}
		}
		assert.True(t, found)
	})

	t.Run("should not duplicate the employee or the outbox event on repeated registration", func(t *testing.T) {
		e := employee.NewEmployee("test-agent-3", "Test Agent Three")

		require.NoError(t, repo.Register(context.Background(), e))
		require.NoError(t, repo.Register(context.Background(), e))

		employees, err := repo.List(context.Background())
		require.NoError(t, err)

		count := 0
		for _, got := range employees {
			if got.GetID() == "test-agent-3" {
				count++
			}
		}
		assert.Equal(t, 1, count)
	})
}
