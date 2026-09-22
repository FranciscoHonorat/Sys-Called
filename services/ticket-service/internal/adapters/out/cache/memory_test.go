package cache_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/cache"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestTicket(t *testing.T) *ticket.Ticket {
	t.Helper()

	id := valueobjects.NewID(uuid.New())
	title, err := valueobjects.NewTitle("Valid Title")
	require.NoError(t, err)
	description, err := valueobjects.NewDescription("Valid Description")
	require.NoError(t, err)

	tk, err := ticket.NewTicket(id, title, description, valueobjects.TicketStatusOpen, nil, nil)
	require.NoError(t, err)

	return tk
}

func TestInMemoryTicketCache(t *testing.T) {
	t.Run("should return a miss for an unknown ticket", func(t *testing.T) {
		c := cache.NewInMemoryTicketCache()

		_, ok := c.Get(context.Background(), uuid.New())

		assert.False(t, ok)
	})

	t.Run("should store and retrieve a ticket", func(t *testing.T) {
		c := cache.NewInMemoryTicketCache()
		tk := newTestTicket(t)

		c.Set(context.Background(), tk)
		got, ok := c.Get(context.Background(), tk.GetID().GetID())

		assert.True(t, ok)
		assert.Same(t, tk, got)
	})

}
