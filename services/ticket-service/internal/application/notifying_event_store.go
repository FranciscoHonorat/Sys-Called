package application

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type NotifyingEventStore struct {
	out.EventStore
	notifications out.NotificationStore
}

func NewNotifyingEventStore(store out.EventStore, notifications out.NotificationStore) *NotifyingEventStore {
	return &NotifyingEventStore{EventStore: store, notifications: notifications}
}

func (s *NotifyingEventStore) Append(ctx context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int, actorID string) error {
	if err := s.EventStore.Append(ctx, aggregateID, events, expectedVersion, actorID); err != nil {
		return err
	}

	if err := s.notify(ctx, aggregateID, events, actorID); err != nil {
		log.Printf("notifications for ticket %s were not recorded: %v", aggregateID, err)
	}
	return nil
}

func (s *NotifyingEventStore) notify(ctx context.Context, aggregateID uuid.UUID, events []event.Event, actorID string) error {
	history, err := s.EventStore.Load(ctx, aggregateID)
	if err != nil {
		return err
	}
	t, err := ticket.LoadFromHistory(history)
	if err != nil {
		return err
	}

	var pending []notification.Notification
	for _, e := range events {
		pending = append(pending, notification.ForTicketEvent(t, e, actorID)...)
	}
	if len(pending) == 0 {
		return nil
	}
	return s.notifications.Add(ctx, pending)
}
