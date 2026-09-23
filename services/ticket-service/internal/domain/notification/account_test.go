package notification_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
)

func TestAccountNotifications(t *testing.T) {
	at := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)

	t.Run("a new account asks the admins for approval", func(t *testing.T) {
		n := notification.ForNewAccount("user-9", "Maria Lima", at)

		assert.Equal(t, notification.ToRole(actor.RoleAdmin), n.Audience)
		assert.Equal(t, "Nova conta aguardando aprovação: Maria Lima", n.Message)
		assert.Equal(t, uuid.Nil, n.TicketID)
		assert.Equal(t, "user-9", n.ActorID)
		assert.Equal(t, at, n.CreatedAt)
		assert.NotEqual(t, uuid.Nil, n.ID)
	})

	t.Run("a forgotten password asks the admins for a temporary one", func(t *testing.T) {
		n := notification.ForPasswordResetRequest("user-1", "Usuário Padrão", at)

		assert.Equal(t, notification.ToRole(actor.RoleAdmin), n.Audience)
		assert.Equal(t, "Pedido de nova senha: Usuário Padrão", n.Message)
		assert.Equal(t, uuid.Nil, n.TicketID)
		assert.False(t, n.IsAboutTicket())
	})
}
