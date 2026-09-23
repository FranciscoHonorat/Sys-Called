package actor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func TestActor(t *testing.T) {
	t.Run("should identify who is acting and with which role", func(t *testing.T) {
		a, err := actor.New("admin-1", "Administradora", "admin")

		require.NoError(t, err)
		assert.Equal(t, "admin-1", a.ID())
		assert.Equal(t, "Administradora", a.Name())
		assert.Equal(t, actor.RoleAdmin, a.Role())
	})

	t.Run("should require an ID", func(t *testing.T) {
		_, err := actor.New("", "Someone", "user")

		assert.ErrorIs(t, err, domainErr.ErrInvalidActor)
	})

	t.Run("should require a known role", func(t *testing.T) {
		_, err := actor.New("someone", "Someone", "root")

		assert.ErrorIs(t, err, domainErr.ErrInvalidActor)
	})
}
