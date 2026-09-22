package outtest

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
)

type EventStore struct {
	streams map[uuid.UUID][]event.Event
	order   []uuid.UUID
}

func NewEventStore() *EventStore {
	return &EventStore{streams: make(map[uuid.UUID][]event.Event)}
}

func (s *EventStore) Append(_ context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int) error {
	current := s.streams[aggregateID]
	if len(current) != expectedVersion {
		return domainErr.ErrConcurrencyConflict
	}
	if len(current) == 0 {
		s.order = append(s.order, aggregateID)
	}
	s.streams[aggregateID] = append(current, events...)
	return nil
}

func (s *EventStore) Load(_ context.Context, aggregateID uuid.UUID) ([]event.Event, error) {
	events, ok := s.streams[aggregateID]
	if !ok {
		return nil, domainErr.ErrEventStreamNotFound
	}
	return events, nil
}

func (s *EventStore) ListAggregateIDs(_ context.Context) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(s.order))
	copy(ids, s.order)
	return ids, nil
}

func (s *EventStore) ForgetAll() {
	s.streams = make(map[uuid.UUID][]event.Event)
	s.order = nil
}

type ResponsibleDirectory struct {
	Responsibles []string
	Names        map[string]string
	Err          error
}

func (d *ResponsibleDirectory) List(_ context.Context) ([]string, error) {
	return d.Responsibles, d.Err
}

func (d *ResponsibleDirectory) Upsert(_ context.Context, id, name string) error {
	if d.Err != nil {
		return d.Err
	}
	if d.Names == nil {
		d.Names = make(map[string]string)
	}
	if _, exists := d.Names[id]; !exists {
		d.Responsibles = append(d.Responsibles, id)
	}
	d.Names[id] = name
	return nil
}
