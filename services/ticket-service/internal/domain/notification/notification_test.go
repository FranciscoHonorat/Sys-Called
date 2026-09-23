package notification_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

func openTicket(t *testing.T) *ticket.Ticket {
	t.Helper()
	title, err := valueobjects.NewTitle("Impressora")
	require.NoError(t, err)
	description, err := valueobjects.NewDescription("Não imprime")
	require.NoError(t, err)
	tk, err := ticket.NewTicket(valueobjects.NewID(uuid.New()), title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")
	require.NoError(t, err)
	return tk
}

func assign(t *testing.T, tk *ticket.Ticket, agent string) {
	t.Helper()
	a, err := valueobjects.NewAssigneeID(agent)
	require.NoError(t, err)
	require.NoError(t, tk.AssignTo(&a))
}

func lastEvent(tk *ticket.Ticket) event.Event {
	events := tk.GetUncommittedEvents()
	return events[len(events)-1]
}

func audiences(notifications []notification.Notification) []notification.Audience {
	result := make([]notification.Audience, 0, len(notifications))
	for _, n := range notifications {
		result = append(result, n.Audience)
	}
	return result
}

func TestForTicketEvent(t *testing.T) {
	t.Run("a new ticket reaches the whole support team and the admins", func(t *testing.T) {
		tk := openTicket(t)

		got := notification.ForTicketEvent(tk, lastEvent(tk), "user-1")

		assert.Equal(t, []notification.Audience{notification.ToRole("support"), notification.ToRole("admin")}, audiences(got))
		assert.Equal(t, "Novo chamado: Impressora", got[0].Message)
		assert.Equal(t, tk.GetID().GetID(), got[0].TicketID)
		assert.Equal(t, "user-1", got[0].ActorID)
		assert.Equal(t, lastEvent(tk).OccurredAt(), got[0].CreatedAt)
	})

	t.Run("an assignment reaches the agent who got the ticket", func(t *testing.T) {
		tk := openTicket(t)
		assign(t, tk, "agent-1")

		got := notification.ForTicketEvent(tk, lastEvent(tk), "admin-1")

		assert.Equal(t, []notification.Audience{notification.ToUser("agent-1")}, audiences(got))
		assert.Equal(t, "Chamado atribuído a você: Impressora", got[0].Message)
	})

	t.Run("an agent taking a ticket for themselves is not notified about it", func(t *testing.T) {
		tk := openTicket(t)
		assign(t, tk, "agent-1")

		assert.Empty(t, notification.ForTicketEvent(tk, lastEvent(tk), "agent-1"))
	})

	t.Run("a message reaches the requester and the responsible, except whoever wrote it", func(t *testing.T) {
		tk := openTicket(t)
		assign(t, tk, "agent-1")
		author, err := valueobjects.NewAuthorID("admin-1")
		require.NoError(t, err)
		content, err := valueobjects.NewContent("Pode verificar hoje?")
		require.NoError(t, err)
		r, err := response.NewResponse(valueobjects.NewID(uuid.Nil), tk.GetID(), &author, content)
		require.NoError(t, err)
		require.NoError(t, tk.AddResponse(r))

		got := notification.ForTicketEvent(tk, lastEvent(tk), "admin-1")

		assert.Equal(t, []notification.Audience{notification.ToUser("user-1"), notification.ToUser("agent-1")}, audiences(got))
		assert.Equal(t, "Nova mensagem em: Impressora", got[0].Message)
		assert.Equal(t, []notification.Audience{notification.ToUser("agent-1")}, audiences(notification.ForTicketEvent(tk, lastEvent(tk), "user-1")))
	})

	t.Run("the requester hears when the work starts and when the ticket is closed", func(t *testing.T) {
		tk := openTicket(t)
		assign(t, tk, "agent-1")
		require.NoError(t, tk.MoveToInProgress())
		started := notification.ForTicketEvent(tk, lastEvent(tk), "agent-1")
		require.NoError(t, tk.Close("Troquei o toner"))
		closed := notification.ForTicketEvent(tk, lastEvent(tk), "agent-1")

		assert.Equal(t, "Seu chamado está em atendimento: Impressora", started[0].Message)
		assert.Equal(t, "Seu chamado foi fechado: Impressora", closed[0].Message)
		assert.Equal(t, []notification.Audience{notification.ToUser("user-1")}, audiences(closed))
	})

	t.Run("changes nobody needs to hear about produce nothing", func(t *testing.T) {
		tk := openTicket(t)
		priority, err := valueobjects.NewPriority("High")
		require.NoError(t, err)
		require.NoError(t, tk.ChangePriority(&priority))

		assert.Empty(t, notification.ForTicketEvent(tk, lastEvent(tk), "admin-1"))
	})
}
