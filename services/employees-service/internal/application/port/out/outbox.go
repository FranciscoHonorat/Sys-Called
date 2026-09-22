package out

import (
	"context"

	"github.com/google/uuid"
)

type OutboxEvent struct {
	ID        uuid.UUID
	EventType string
	Payload   []byte
}

type OutboxStore interface {
	FetchPending(ctx context.Context) ([]OutboxEvent, error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
}

type EventPublisher interface {
	Publish(ctx context.Context, eventType string, payload []byte) error
}
