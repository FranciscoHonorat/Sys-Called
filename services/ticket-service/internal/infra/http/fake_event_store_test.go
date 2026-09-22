package httpapi_test

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
)

type fakeEventStore struct {
	streams map[uuid.UUID][]event.Event
}

func newFakeEventStore() *fakeEventStore {
	return &fakeEventStore{streams: make(map[uuid.UUID][]event.Event)}
}

func (s *fakeEventStore) Append(_ context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int) error {
	current := s.streams[aggregateID]
	if len(current) != expectedVersion {
		return domainErr.ErrConcurrencyConflict
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
