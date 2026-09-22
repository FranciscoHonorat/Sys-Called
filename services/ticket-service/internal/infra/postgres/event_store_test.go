package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/postgres"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TICKET_SERVICE_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TICKET_SERVICE_TEST_DATABASE_URL not set; skipping Postgres integration tests")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	return pool
}

func TestEventStore(t *testing.T) {
	pool := testPool(t)
	store := postgres.NewEventStore(pool)

	t.Run("should append and load an event stream", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		opened := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil)

		err = store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0)
		assert.NoError(t, err)

		loaded, err := store.Load(context.Background(), id.GetID())
		assert.NoError(t, err)
		assert.Len(t, loaded, 1)
		assert.Equal(t, "TicketOpened", loaded[0].EventName())
		assert.Equal(t, id.GetID(), loaded[0].AggregateID())

		reopened, ok := loaded[0].(event.TicketOpened)
		assert.True(t, ok)
		assert.Equal(t, "Valid Title", reopened.Title)
	})

	t.Run("should append multiple events across separate calls in order", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		opened := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil)
		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0))

		assignee, err := valueobjects.NewAssigneeID("agent-1")
		require.NoError(t, err)
		assigned := event.NewTicketAssigned(id, &assignee)
		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{assigned}, 1))

		loaded, err := store.Load(context.Background(), id.GetID())
		assert.NoError(t, err)
		assert.Len(t, loaded, 2)
		assert.Equal(t, "TicketOpened", loaded[0].EventName())
		assert.Equal(t, "TicketAssigned", loaded[1].EventName())
	})

	t.Run("should return an error when expectedVersion is stale", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		opened := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil)
		require.NoError(t, store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0))

		err = store.Append(context.Background(), id.GetID(), []event.Event{opened}, 0)
		assert.ErrorIs(t, err, domainErr.ErrConcurrencyConflict)
	})

	t.Run("should return an error when loading a stream that does not exist", func(t *testing.T) {
		_, err := store.Load(context.Background(), uuid.New())
		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should list the aggregate IDs of opened tickets", func(t *testing.T) {
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)

		id1 := valueobjects.NewID(uuid.New())
		opened1 := event.NewTicketOpened(id1, title, description, valueobjects.TicketStatusOpen, nil, nil)
		require.NoError(t, store.Append(context.Background(), id1.GetID(), []event.Event{opened1}, 0))

		id2 := valueobjects.NewID(uuid.New())
		opened2 := event.NewTicketOpened(id2, title, description, valueobjects.TicketStatusOpen, nil, nil)
		require.NoError(t, store.Append(context.Background(), id2.GetID(), []event.Event{opened2}, 0))

		ids, err := store.ListAggregateIDs(context.Background())
		assert.NoError(t, err)
		assert.Contains(t, ids, id1.GetID())
		assert.Contains(t, ids, id2.GetID())
	})
}
