package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

type PublishPendingEventsUseCase struct {
	store     out.OutboxStore
	publisher out.EventPublisher
}

func NewPublishPendingEventsUseCase(store out.OutboxStore, publisher out.EventPublisher) *PublishPendingEventsUseCase {
	return &PublishPendingEventsUseCase{store: store, publisher: publisher}
}

func (uc *PublishPendingEventsUseCase) Execute(ctx context.Context) error {
	events, err := uc.store.FetchPending(ctx)
	if err != nil {
		return err
	}

	for _, e := range events {
		if err := uc.publisher.Publish(ctx, e.EventType, e.Payload); err != nil {
			return err
		}
		if err := uc.store.MarkPublished(ctx, e.ID); err != nil {
			return err
		}
	}

	return nil
}
