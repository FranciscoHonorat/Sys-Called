package application_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

type fakeOutboxStore struct {
	pending   []out.OutboxEvent
	published []uuid.UUID
	fetchErr  error
}

func (s *fakeOutboxStore) FetchPending(_ context.Context) ([]out.OutboxEvent, error) {
	return s.pending, s.fetchErr
}

func (s *fakeOutboxStore) MarkPublished(_ context.Context, id uuid.UUID) error {
	s.published = append(s.published, id)
	return nil
}

type fakeEventPublisher struct {
	published []out.OutboxEvent
	err       error
}

func (p *fakeEventPublisher) Publish(_ context.Context, eventType string, payload []byte) error {
	if p.err != nil {
		return p.err
	}
	p.published = append(p.published, out.OutboxEvent{EventType: eventType, Payload: payload})
	return nil
}

func TestPublishPendingEventsUseCase(t *testing.T) {
	t.Run("should publish pending events and mark them as published", func(t *testing.T) {
		id1, id2 := uuid.New(), uuid.New()
		store := &fakeOutboxStore{pending: []out.OutboxEvent{
			{ID: id1, EventType: "EmployeeRegistered", Payload: []byte(`{"id":"agent-1"}`)},
			{ID: id2, EventType: "EmployeeRegistered", Payload: []byte(`{"id":"agent-2"}`)},
		}}
		publisher := &fakeEventPublisher{}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Len(t, publisher.published, 2)
		assert.Equal(t, []uuid.UUID{id1, id2}, store.published)
	})

	t.Run("should do nothing when there are no pending events", func(t *testing.T) {
		store := &fakeOutboxStore{}
		publisher := &fakeEventPublisher{}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		require.NoError(t, err)
		assert.Empty(t, publisher.published)
	})

	t.Run("should stop and not mark as published when publishing fails", func(t *testing.T) {
		id1 := uuid.New()
		store := &fakeOutboxStore{pending: []out.OutboxEvent{
			{ID: id1, EventType: "EmployeeRegistered", Payload: []byte(`{}`)},
		}}
		publisher := &fakeEventPublisher{err: assert.AnError}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, store.published)
	})

	t.Run("should propagate errors from fetching pending events", func(t *testing.T) {
		store := &fakeOutboxStore{fetchErr: assert.AnError}
		publisher := &fakeEventPublisher{}
		uc := application.NewPublishPendingEventsUseCase(store, publisher)

		err := uc.Execute(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
