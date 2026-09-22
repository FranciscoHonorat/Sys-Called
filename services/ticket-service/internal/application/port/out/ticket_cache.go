package out

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type TicketCache interface {
	Get(ctx context.Context, aggregateID uuid.UUID) (*ticket.Ticket, bool)
	Set(ctx context.Context, t *ticket.Ticket)
}
