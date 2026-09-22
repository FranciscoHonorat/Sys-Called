package query_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
)

func newTestCache() *cache.InMemoryTicketCache {
	return cache.NewInMemoryTicketCache()
}

func openTestTicket(t *testing.T, store *outtest.EventStore) string {
	t.Helper()

	output, err := command.NewOpenTicketUseCase(store, newTestCache()).Execute(context.Background(), command.OpenTicketInput{
		Title:       "Valid Title",
		Description: "Valid Description",
	})
	assert.NoError(t, err)

	return output.TicketID
}
