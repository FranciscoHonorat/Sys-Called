package outbox_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/infra/outbox"
)

type fakeStore struct {
	pending   []outbox.Event
	published []uuid.UUID
	fetchErr  error
}

func (s *fakeStore) FetchPending(_ context.Context) ([]outbox.Event, error) {
	return s.pending, s.fetchErr
}

func (s *fakeStore) MarkPublished(_ context.Context, id uuid.UUID) error {
	s.published = append(s.published, id)
	return nil
}

type fakePublisher struct {
	published []outbox.Event
	err       error
}

func (p *fakePublisher) Publish(_ context.Context, eventType string, payload []byte) error {
	if p.err != nil {
		return p.err
	}
	p.published = append(p.published, outbox.Event{EventType: eventType, Payload: payload})
	return nil
}

func TestRelay(t *testing.T) {
	t.Run("should publish pending events and mark them as published", func(t *testing.T) {
		id1, id2 := uuid.New(), uuid.New()
		store := &fakeStore{pending: []outbox.Event{
			{ID: id1, EventType: "EmployeeRegistered", Payload: []byte(`{"id":"agent-1"}`)},
			{ID: id2, EventType: "EmployeeRegistered", Payload: []byte(`{"id":"agent-2"}`)},
		}}
		publisher := &fakePublisher{}
		relay := outbox.NewRelay(store, publisher)

		err := relay.PublishPending(context.Background())

		require.NoError(t, err)
		assert.Len(t, publisher.published, 2)
		assert.Equal(t, []uuid.UUID{id1, id2}, store.published)
	})

	t.Run("should do nothing when there are no pending events", func(t *testing.T) {
		store := &fakeStore{}
		publisher := &fakePublisher{}
		relay := outbox.NewRelay(store, publisher)

		err := relay.PublishPending(context.Background())

		require.NoError(t, err)
		assert.Empty(t, publisher.published)
	})

	t.Run("should stop and not mark as published when publishing fails", func(t *testing.T) {
		id1 := uuid.New()
		store := &fakeStore{pending: []outbox.Event{
			{ID: id1, EventType: "EmployeeRegistered", Payload: []byte(`{}`)},
		}}
		publisher := &fakePublisher{err: assert.AnError}
		relay := outbox.NewRelay(store, publisher)

		err := relay.PublishPending(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
		assert.Empty(t, store.published)
	})

	t.Run("should propagate errors from fetching pending events", func(t *testing.T) {
		store := &fakeStore{fetchErr: assert.AnError}
		publisher := &fakePublisher{}
		relay := outbox.NewRelay(store, publisher)

		err := relay.PublishPending(context.Background())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
