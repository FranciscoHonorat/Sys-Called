package session_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/session"
)

func TestRefreshToken(t *testing.T) {
	now := time.Date(2026, 9, 22, 10, 0, 0, 0, time.UTC)

	t.Run("should belong to an employee and start active", func(t *testing.T) {
		token := session.NewRefreshToken("hash", "admin-1", now.Add(time.Hour))

		assert.Equal(t, "hash", token.Hash())
		assert.Equal(t, "admin-1", token.EmployeeID())
		assert.False(t, token.IsRevoked())
		assert.False(t, token.IsExpired(now))
	})

	t.Run("should expire once the expiration time is reached", func(t *testing.T) {
		token := session.NewRefreshToken("hash", "admin-1", now.Add(time.Hour))

		assert.False(t, token.IsExpired(now.Add(time.Hour-time.Second)))
		assert.True(t, token.IsExpired(now.Add(time.Hour)))
	})

	t.Run("should be revocable", func(t *testing.T) {
		token := session.NewRefreshToken("hash", "admin-1", now.Add(time.Hour))

		token.Revoke()

		assert.True(t, token.IsRevoked())
	})
}
