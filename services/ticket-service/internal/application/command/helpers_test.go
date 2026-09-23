package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
)

func newTestCache() *cache.InMemoryTicketCache {
	return cache.NewInMemoryTicketCache()
}

var testUser = outtest.Actor("user-1", "user")

func openAssignedTicket(t *testing.T, store *outtest.EventStore) string {
	t.Helper()
	ticketID := openTestTicket(t, store)
	require.NoError(t, command.NewAssignTicketUseCase(store, newTestCache()).Execute(context.Background(),
		command.AssignTicketInput{Actor: testAdmin, TicketID: ticketID, AssigneeID: assignedAgent.ID()}))
	return ticketID
}
