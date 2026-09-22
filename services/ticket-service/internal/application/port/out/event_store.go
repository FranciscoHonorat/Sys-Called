package out

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
)

type EventStore interface {
	Append(ctx context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int) error
	Load(ctx context.Context, aggregateID uuid.UUID) ([]event.Event, error)
	ListAggregateIDs(ctx context.Context) ([]uuid.UUID, error)
}
