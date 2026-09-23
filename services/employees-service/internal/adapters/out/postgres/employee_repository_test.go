package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
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

func uniqueID() string {
	return "test-" + uuid.NewString()
}

func outboxEventCount(t *testing.T, pool *pgxpool.Pool, employeeID string) int {
	t.Helper()

	var count int
	err := pool.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM outbox_events WHERE event_type = 'EmployeeRegistered' AND payload->>'id' = $1`,
		employeeID,
	).Scan(&count)
	require.NoError(t, err)

	return count
}

func findEmployee(t *testing.T, repo *postgres.EmployeeRepository, id string) []employee.Employee {
	t.Helper()

	employees, err := repo.List(context.Background())
	require.NoError(t, err)

	var found []employee.Employee
	for _, e := range employees {
		if e.GetID() == id {
			found = append(found, e)
		}
	}
	return found
}

func TestEmployeeRepository(t *testing.T) {
	pool := testPool(t)
	repo := postgres.NewEmployeeRepository(pool)

	t.Run("should register an employee and list it back with username, role and password hash", func(t *testing.T) {
		id := uniqueID()
		e := employee.NewEmployee(id, "Test Agent", id+"-username", employee.RoleSupport, "hash")

		require.NoError(t, repo.Register(context.Background(), e))

		found := findEmployee(t, repo, id)
		require.Len(t, found, 1)
		assert.Equal(t, "Test Agent", found[0].GetName())
		assert.Equal(t, id+"-username", found[0].GetUsername())
		assert.Equal(t, employee.RoleSupport, found[0].GetRole())
		assert.Equal(t, "hash", found[0].GetPasswordHash())
	})

	t.Run("should find an employee by username", func(t *testing.T) {
		id := uniqueID()
		require.NoError(t, repo.Register(context.Background(),
			employee.NewEmployee(id, "Test Agent", id+"-username", employee.RoleAdmin, "hash")))

		found, ok, err := repo.FindByUsername(context.Background(), id+"-username")

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, id, found.GetID())
		assert.Equal(t, employee.RoleAdmin, found.GetRole())
		assert.Equal(t, "hash", found.GetPasswordHash())
	})

	t.Run("should find an employee by ID", func(t *testing.T) {
		id := uniqueID()
		require.NoError(t, repo.Register(context.Background(),
			employee.NewEmployee(id, "Test Agent", id+"-username", employee.RoleSupport, "hash")))

		found, ok, err := repo.FindByID(context.Background(), id)

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, id+"-username", found.GetUsername())
		assert.Equal(t, employee.RoleSupport, found.GetRole())
	})

	t.Run("should report when no employee has the username", func(t *testing.T) {
		_, ok, err := repo.FindByUsername(context.Background(), uniqueID())

		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("should write an outbox event when registering a new employee", func(t *testing.T) {
		id := uniqueID()
		e, err := employee.Register(id, "Test Agent", id+"-username", employee.RoleSupport, "hash")
		require.NoError(t, err)

		require.NoError(t, repo.Register(context.Background(), e))

		assert.Equal(t, 1, outboxEventCount(t, pool, id))
	})

	t.Run("should not duplicate the employee or the outbox event on repeated registration", func(t *testing.T) {
		id := uniqueID()
		e, err := employee.Register(id, "Test Agent", id+"-username", employee.RoleSupport, "hash")
		require.NoError(t, err)

		require.NoError(t, repo.Register(context.Background(), e))
		require.NoError(t, repo.Register(context.Background(), e))

		assert.Len(t, findEmployee(t, repo, id), 1)
		assert.Equal(t, 1, outboxEventCount(t, pool, id))
	})

	t.Run("should give credentials to an employee registered before accounts existed", func(t *testing.T) {
		id := uniqueID()
		_, err := pool.Exec(context.Background(),
			`INSERT INTO employees (id, name, registered_at) VALUES ($1, 'Legacy Agent', now())`, id)
		require.NoError(t, err)
		e, err := employee.Register(id, "Legacy Agent", id+"-username", employee.RoleSupport, "hash")
		require.NoError(t, err)

		require.NoError(t, repo.Register(context.Background(), e))

		found, ok, err := repo.FindByUsername(context.Background(), id+"-username")
		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, employee.RoleSupport, found.GetRole())
		assert.Equal(t, "hash", found.GetPasswordHash())
		assert.Equal(t, 1, outboxEventCount(t, pool, id))
	})
}
