package ticket_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

func ticketWith(t *testing.T, requester, assignee string, status valueobjects.Status) *ticket.Ticket {
	t.Helper()
	id, title, description, _, _, _ := validTicketParts(t)
	tk, err := ticket.NewTicket(id, title, description, valueobjects.TicketStatusOpen, nil, nil, requester)
	require.NoError(t, err)
	if assignee != "" {
		a, err := valueobjects.NewAssigneeID(assignee)
		require.NoError(t, err)
		require.NoError(t, tk.AssignTo(&a))
	}
	switch status {
	case valueobjects.TicketStatusInProgress:
		require.NoError(t, tk.MoveToInProgress())
	case valueobjects.TicketStatusClosed:
		require.NoError(t, tk.Close("Resolvido"))
	}
	return tk
}

func actorOf(t *testing.T, id, role string) actor.Actor {
	t.Helper()
	a, err := actor.New(id, id, role)
	require.NoError(t, err)
	return a
}

func TestTicketVisibility(t *testing.T) {
	open, inProgress, closed := valueobjects.TicketStatusOpen, valueobjects.TicketStatusInProgress, valueobjects.TicketStatusClosed

	cases := []struct {
		name    string
		viewer  actor.Actor
		ticket  *ticket.Ticket
		visible bool
	}{
		{"admin sees a ticket from anyone", actorOf(t, "admin-1", "admin"), ticketWith(t, "user-1", "agent-2", closed), true},
		{"user sees a ticket they opened", actorOf(t, "user-1", "user"), ticketWith(t, "user-1", "agent-1", inProgress), true},
		{"user does not see a ticket opened by someone else", actorOf(t, "user-1", "user"), ticketWith(t, "user-2", "", open), false},
		{"support sees an open unassigned ticket", actorOf(t, "agent-1", "support"), ticketWith(t, "user-1", "", open), true},
		{"support sees an open ticket assigned to another agent", actorOf(t, "agent-1", "support"), ticketWith(t, "user-1", "agent-2", open), true},
		{"support sees a ticket they are working on", actorOf(t, "agent-1", "support"), ticketWith(t, "user-1", "agent-1", inProgress), true},
		{"support sees a ticket they closed", actorOf(t, "agent-1", "support"), ticketWith(t, "user-1", "agent-1", closed), true},
		{"support does not see another agent's ticket in progress", actorOf(t, "agent-1", "support"), ticketWith(t, "user-1", "agent-2", inProgress), false},
		{"support does not see another agent's closed ticket", actorOf(t, "agent-1", "support"), ticketWith(t, "user-1", "agent-2", closed), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.visible, tc.ticket.IsVisibleTo(tc.viewer))
		})
	}
}
