package outbox

import (
	"context"
	"log"
	"time"
)

type pendingEventsPublisher interface {
	Execute(ctx context.Context) error
}

type Relay struct {
	publishPending pendingEventsPublisher
}

func NewRelay(publishPending pendingEventsPublisher) *Relay {
	return &Relay{publishPending: publishPending}
}

func (r *Relay) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.publishPending.Execute(ctx); err != nil {
				log.Printf("outbox relay: %v", err)
			}
		}
	}
}
