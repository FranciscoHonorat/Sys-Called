package postgres_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/postgres"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestAccountPersistence(t *testing.T) {
	pool := testPool(t)
	repo := postgres.NewEmployeeRepository(pool)
	ctx := context.Background()
	countEvents := func(employeeID, eventType string) int {
		var count int
		require.NoError(t, pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM outbox_events WHERE event_type = $1 AND payload->>'id' = $2`, eventType, employeeID,
		).Scan(&count))
		return count
	}

	t.Run("keeps a sign up waiting for approval and announces it", func(t *testing.T) {
		id := uniqueID()
		e, err := employee.SignUp(id, "Maria Lima", id+"-username", "hash")
		require.NoError(t, err)

		require.NoError(t, repo.Register(ctx, e))

		found, ok, err := repo.FindByID(ctx, id)
		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, employee.StatusPending, found.GetStatus())
		assert.Equal(t, 1, countEvents(id, "EmployeeSignedUp"))
	})

	t.Run("refuses a username that is already taken", func(t *testing.T) {
		id := uniqueID()
		require.NoError(t, repo.Register(ctx, employee.NewEmployee(id, "Primeira", id+"-username", employee.RoleUser, "hash")))
		second, err := employee.SignUp(uniqueID(), "Segunda", id+"-username", "hash")
		require.NoError(t, err)

		assert.ErrorIs(t, repo.Register(ctx, second), domainErr.ErrUsernameTaken)
	})

	t.Run("saves the account changes and their events", func(t *testing.T) {
		id := uniqueID()
		e, err := employee.SignUp(id, "Maria Lima", id+"-username", "hash")
		require.NoError(t, err)
		require.NoError(t, repo.Register(ctx, e))
		stored, _, err := repo.FindByID(ctx, id)
		require.NoError(t, err)

		stored.Approve()
		stored.RequestPasswordReset()
		require.NoError(t, repo.Save(ctx, stored))

		found, _, err := repo.FindByUsername(ctx, id+"-username")
		require.NoError(t, err)
		assert.Equal(t, employee.StatusActive, found.GetStatus())
		assert.True(t, found.HasRequestedPasswordReset())
		assert.Equal(t, 1, countEvents(id, "PasswordResetRequested"))

		found.SetTemporaryPassword("temporary-hash")
		require.NoError(t, repo.Save(ctx, found))
		reloaded, _, err := repo.FindByID(ctx, id)
		require.NoError(t, err)
		assert.True(t, reloaded.MustChangePassword())
		assert.False(t, reloaded.HasRequestedPasswordReset())
		assert.Equal(t, "temporary-hash", reloaded.GetPasswordHash())
	})
}
