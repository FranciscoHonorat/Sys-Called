package ticket_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

func TestTicketPermissions(t *testing.T) {
	actors := map[string]actor.Actor{
		"requester":      actorOf(t, "user-1", "user"),
		"another user":   actorOf(t, "user-2", "user"),
		"assigned agent": actorOf(t, "agent-1", "support"),
		"another agent":  actorOf(t, "agent-2", "support"),
		"admin":          actorOf(t, "admin-1", "admin"),
	}

	cases := []struct {
		action     string
		permission ticket.Permission
		allowed    []string
		denied     []string
	}{
		{"edit", ticket.CanEdit, []string{"requester", "admin"}, []string{"another user", "assigned agent", "another agent"}},
		{"respond", ticket.CanRespond, []string{"requester", "assigned agent", "another agent", "admin"}, []string{"another user"}},
		{"manage", ticket.CanManage, []string{"assigned agent", "another agent", "admin"}, []string{"requester", "another user"}},
		{"work", ticket.CanWork, []string{"assigned agent"}, []string{"requester", "another user", "another agent", "admin"}},
	}
	for _, tc := range cases {
		tk := ticketWith(t, "user-1", "agent-1", valueobjects.TicketStatusOpen)
		for _, who := range tc.allowed {
			t.Run(who+" can "+tc.action, func(t *testing.T) {
				assert.True(t, tc.permission(tk, actors[who]))
			})
		}
		for _, who := range tc.denied {
			t.Run(who+" cannot "+tc.action, func(t *testing.T) {
				assert.False(t, tc.permission(tk, actors[who]))
			})
		}
	}

	t.Run("the requester can no longer edit once someone started working on the ticket", func(t *testing.T) {
		inProgress := ticketWith(t, "user-1", "agent-1", valueobjects.TicketStatusInProgress)

		assert.False(t, ticket.CanEdit(inProgress, actors["requester"]))
		assert.True(t, ticket.CanEdit(inProgress, actors["admin"]))
	})

	t.Run("a support agent works only on tickets assigned to them but can still take an unassigned one", func(t *testing.T) {
		unassigned := ticketWith(t, "user-1", "", valueobjects.TicketStatusOpen)

		assert.False(t, ticket.CanWork(unassigned, actors["another agent"]))
		assert.True(t, ticket.CanManage(unassigned, actors["another agent"]))
	})
}
