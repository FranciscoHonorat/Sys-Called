package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/session"
)

func TestRefreshTokenStore(t *testing.T) {
	store := postgres.NewRefreshTokenStore(testPool(t))
	ctx := context.Background()
	expiresAt := time.Now().Add(time.Hour).UTC().Truncate(time.Microsecond)

	t.Run("should save and find a token by its hash", func(t *testing.T) {
		hash, employeeID := uniqueID(), uniqueID()
		require.NoError(t, store.Save(ctx, session.NewRefreshToken(hash, employeeID, expiresAt)))

		found, ok, err := store.FindByHash(ctx, hash)

		require.NoError(t, err)
		require.True(t, ok)
		assert.Equal(t, employeeID, found.EmployeeID())
		assert.False(t, found.IsRevoked())
		assert.False(t, found.IsExpired(expiresAt.Add(-time.Second)))
		assert.True(t, found.IsExpired(expiresAt))
	})

	t.Run("should report an unknown hash", func(t *testing.T) {
		_, ok, err := store.FindByHash(ctx, uniqueID())

		require.NoError(t, err)
		assert.False(t, ok)
	})

	t.Run("should revoke a single token", func(t *testing.T) {
		hash := uniqueID()
		require.NoError(t, store.Save(ctx, session.NewRefreshToken(hash, uniqueID(), expiresAt)))

		revoked, err := store.Revoke(ctx, hash)
		require.NoError(t, err)
		assert.True(t, revoked)

		found, _, err := store.FindByHash(ctx, hash)
		require.NoError(t, err)
		assert.True(t, found.IsRevoked())
	})

	t.Run("should report that a token was already revoked", func(t *testing.T) {
		hash := uniqueID()
		require.NoError(t, store.Save(ctx, session.NewRefreshToken(hash, uniqueID(), expiresAt)))
		_, err := store.Revoke(ctx, hash)
		require.NoError(t, err)

		revoked, err := store.Revoke(ctx, hash)

		require.NoError(t, err)
		assert.False(t, revoked)
	})

	t.Run("should revoke every token of one employee only", func(t *testing.T) {
		employeeID, otherEmployeeID := uniqueID(), uniqueID()
		first, second, other := uniqueID(), uniqueID(), uniqueID()
		require.NoError(t, store.Save(ctx, session.NewRefreshToken(first, employeeID, expiresAt)))
		require.NoError(t, store.Save(ctx, session.NewRefreshToken(second, employeeID, expiresAt)))
		require.NoError(t, store.Save(ctx, session.NewRefreshToken(other, otherEmployeeID, expiresAt)))

		require.NoError(t, store.RevokeAllForEmployee(ctx, employeeID))

		for hash, wantRevoked := range map[string]bool{first: true, second: true, other: false} {
			found, _, err := store.FindByHash(ctx, hash)
			require.NoError(t, err)
			assert.Equal(t, wantRevoked, found.IsRevoked(), hash)
		}
	})
}
