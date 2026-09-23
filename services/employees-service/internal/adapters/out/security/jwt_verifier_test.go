package security_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/security"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestJWTVerifier(t *testing.T) {
	admin := employee.NewEmployee("admin-1", "Administradora", "admin", employee.RoleAdmin, "hash")

	t.Run("accepts tokens this service issued and tells who is calling", func(t *testing.T) {
		public, private := newKeyPair(t)
		token, _, err := security.NewJWTIssuer(private, time.Minute).Issue(admin)
		require.NoError(t, err)

		caller, err := security.NewJWTVerifier(public).Verify(token)

		require.NoError(t, err)
		assert.Equal(t, "admin-1", caller.ID)
		assert.Equal(t, employee.RoleAdmin, caller.Role)
	})

	t.Run("rejects a tampered token", func(t *testing.T) {
		public, private := newKeyPair(t)
		token, _, err := security.NewJWTIssuer(private, time.Minute).Issue(admin)
		require.NoError(t, err)

		_, err = security.NewJWTVerifier(public).Verify(token[:len(token)-2] + "xx")

		assert.Error(t, err)
	})

	t.Run("rejects a token signed by another key", func(t *testing.T) {
		public, _ := newKeyPair(t)
		_, otherPrivate := newKeyPair(t)
		token, _, err := security.NewJWTIssuer(otherPrivate, time.Minute).Issue(admin)
		require.NoError(t, err)

		_, err = security.NewJWTVerifier(public).Verify(token)

		assert.Error(t, err)
	})

	t.Run("rejects an expired token", func(t *testing.T) {
		public, private := newKeyPair(t)
		token, _, err := security.NewJWTIssuer(private, -time.Minute).Issue(admin)
		require.NoError(t, err)

		_, err = security.NewJWTVerifier(public).Verify(token)

		assert.Error(t, err)
	})
}
