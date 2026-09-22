package application_test

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/cache"
)

func newTestCache() *cache.InMemoryTicketCache {
	return cache.NewInMemoryTicketCache()
}

type fakeEventStore struct {
	streams map[uuid.UUID][]event.Event
	order   []uuid.UUID
}

func newFakeEventStore() *fakeEventStore {
	return &fakeEventStore{streams: make(map[uuid.UUID][]event.Event)}
}

func (s *fakeEventStore) Append(_ context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int) error {
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

func (s *fakeEventStore) Load(_ context.Context, aggregateID uuid.UUID) ([]event.Event, error) {
	events, ok := s.streams[aggregateID]
	if !ok {
		return nil, domainErr.ErrEventStreamNotFound
	}
	return events, nil
}

func (s *fakeEventStore) ListAggregateIDs(_ context.Context) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, len(s.order))
	copy(ids, s.order)
	return ids, nil
}

func (s *fakeEventStore) forgetAll() {
	s.streams = make(map[uuid.UUID][]event.Event)
	s.order = nil
}
