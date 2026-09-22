package outbox

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID
	EventType string
	Payload   []byte
}

type Store interface {
	FetchPending(ctx context.Context) ([]Event, error)
	MarkPublished(ctx context.Context, id uuid.UUID) error
}

type Publisher interface {
	Publish(ctx context.Context, eventType string, payload []byte) error
}

type Relay struct {
	store     Store
	publisher Publisher
}

func NewRelay(store Store, publisher Publisher) *Relay {
	return &Relay{store: store, publisher: publisher}
}

func (r *Relay) PublishPending(ctx context.Context) error {
	events, err := r.store.FetchPending(ctx)
	if err != nil {
		return err
	}

	for _, e := range events {
		if err := r.publisher.Publish(ctx, e.EventType, e.Payload); err != nil {
			return err
		}
		if err := r.store.MarkPublished(ctx, e.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *Relay) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.PublishPending(ctx); err != nil {
				log.Printf("outbox relay: %v", err)
			}
		}
	}
}
