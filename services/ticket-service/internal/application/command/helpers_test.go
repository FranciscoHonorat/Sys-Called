package command_test

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/cache"
)

func newTestCache() *cache.InMemoryTicketCache {
	return cache.NewInMemoryTicketCache()
}
