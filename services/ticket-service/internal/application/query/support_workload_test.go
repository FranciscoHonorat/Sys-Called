package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func ticketHandledBy(t *testing.T, store *outtest.EventStore, agentID, stage string) {
	t.Helper()
	ctx, cache := context.Background(), newTestCache()
	id := openTestTicket(t, store)
	agent := outtest.Actor(agentID, "support")
	require.NoError(t, command.NewAssignTicketUseCase(store, cache).Execute(ctx, command.AssignTicketInput{Actor: testAdmin, TicketID: id, AssigneeID: agentID}))
	if stage == "In Progress" || stage == "Closed" {
		require.NoError(t, command.NewMoveTicketToInProgressUseCase(store, cache).Execute(ctx, command.MoveTicketToInProgressInput{Actor: agent, TicketID: id}))
	}
	if stage == "Closed" {
		require.NoError(t, command.NewCloseTicketUseCase(store, cache).Execute(ctx, command.CloseTicketInput{Actor: agent, TicketID: id, Resolution: "Feito"}))
	}
}

func TestSupportWorkloadUseCase(t *testing.T) {
	directory := &outtest.ResponsibleDirectory{
		Responsibles: []string{"agent-1", "agent-2", "agent-3"},
		Names:        map[string]string{"agent-1": "Ana Souza", "agent-2": "Bruno Lima", "agent-3": "Carla Melo"},
	}

	t.Run("counts the open, in progress and closed tickets of each agent", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketHandledBy(t, store, "agent-1", "Open")
		ticketHandledBy(t, store, "agent-1", "In Progress")
		ticketHandledBy(t, store, "agent-2", "Closed")
		openTestTicket(t, store)

		output, err := query.NewSupportWorkloadUseCase(store, newTestCache(), directory).Execute(context.Background(), testAdmin)

		require.NoError(t, err)
		assert.Equal(t, []query.AgentWorkload{
			{ID: "agent-1", Name: "Ana Souza", Open: 1, InProgress: 1, Closed: 0},
			{ID: "agent-2", Name: "Bruno Lima", Open: 0, InProgress: 0, Closed: 1},
			{ID: "agent-3", Name: "Carla Melo", Open: 0, InProgress: 0, Closed: 0},
		}, output)
	})

	t.Run("is only for admins", func(t *testing.T) {
		_, err := query.NewSupportWorkloadUseCase(outtest.NewEventStore(), newTestCache(), directory).
			Execute(context.Background(), outtest.Actor("agent-1", "support"))

		assert.ErrorIs(t, err, domainErr.ErrForbidden)
	})
}
