package cache

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type InMemoryTicketCache struct {
	mu      sync.RWMutex
	tickets map[uuid.UUID]*ticket.Ticket
}

func NewInMemoryTicketCache() *InMemoryTicketCache {
	return &InMemoryTicketCache{tickets: make(map[uuid.UUID]*ticket.Ticket)}
}

func (c *InMemoryTicketCache) Get(_ context.Context, aggregateID uuid.UUID) (*ticket.Ticket, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	t, ok := c.tickets[aggregateID]
	return t, ok
}

func (c *InMemoryTicketCache) Set(_ context.Context, t *ticket.Ticket) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tickets[t.GetID().GetID()] = t
}
