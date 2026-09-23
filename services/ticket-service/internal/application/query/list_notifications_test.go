package query_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
)

var base = time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)

func sent(audience notification.Audience, message string, minutes int, actorID string) notification.Notification {
	return notification.Notification{
		ID: uuid.New(), Audience: audience, Message: message, TicketID: uuid.New(), ActorID: actorID,
		CreatedAt: base.Add(time.Duration(minutes) * time.Minute),
	}
}

func TestListNotificationsUseCase(t *testing.T) {
	agent := outtest.Actor("agent-1", "support")

	t.Run("lists what reached the viewer, newest first, telling which are unread", func(t *testing.T) {
		store := &outtest.NotificationStore{Added: []notification.Notification{
			sent(notification.ToRole("support"), "Novo chamado: Impressora", 0, "user-1"),
			sent(notification.ToUser("agent-2"), "Chamado atribuído a você: Outro", 1, "admin-1"),
			sent(notification.ToUser("agent-1"), "Chamado atribuído a você: Impressora", 2, "admin-1"),
		}}
		require.NoError(t, store.MarkSeen(context.Background(), "agent-1", base))

		output, err := query.NewListNotificationsUseCase(store).Execute(context.Background(), agent)

		require.NoError(t, err)
		assert.Equal(t, 1, output.Unread)
		require.Len(t, output.Items, 2)
		assert.Equal(t, "Chamado atribuído a você: Impressora", output.Items[0].Message)
		assert.True(t, output.Items[0].Unread)
		assert.False(t, output.Items[1].Unread)
	})

	t.Run("marks everything the viewer has seen as read", func(t *testing.T) {
		store := &outtest.NotificationStore{Added: []notification.Notification{
			sent(notification.ToUser("agent-1"), "Chamado atribuído a você: Impressora", 2, "admin-1"),
		}}
		now := base.Add(time.Hour)

		require.NoError(t, command.NewMarkNotificationsReadUseCase(store, func() time.Time { return now }).Execute(context.Background(), agent))

		output, err := query.NewListNotificationsUseCase(store).Execute(context.Background(), agent)
		require.NoError(t, err)
		assert.Zero(t, output.Unread)
	})

	t.Run("leaves out the ticket of notifications about accounts", func(t *testing.T) {
		admin := outtest.Actor("admin-1", "admin")
		store := &outtest.NotificationStore{Added: []notification.Notification{
			notification.ForNewAccount("user-9", "Maria Lima", base),
		}}

		output, err := query.NewListNotificationsUseCase(store).Execute(context.Background(), admin)

		require.NoError(t, err)
		require.Len(t, output.Items, 1)
		assert.Empty(t, output.Items[0].TicketID)
	})
}
