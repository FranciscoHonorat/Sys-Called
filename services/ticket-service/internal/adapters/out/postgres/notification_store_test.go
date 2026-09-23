package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/out/postgres"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
)

func testActor(t *testing.T, pool *pgxpool.Pool, role string) actor.Actor {
	t.Helper()
	a, err := actor.New("test-"+uuid.NewString(), "Test", role)
	require.NoError(t, err)
	t.Cleanup(func() {
		ctx := context.Background()
		_, _ = pool.Exec(ctx, `DELETE FROM notifications WHERE audience = $1 OR actor_id = $1`, a.ID())
		_, _ = pool.Exec(ctx, `DELETE FROM notification_reads WHERE user_id = $1`, a.ID())
	})
	return a
}

func TestNotificationStore(t *testing.T) {
	pool := testPool(t)
	store := postgres.NewNotificationStore(pool)
	ctx := context.Background()
	now := time.Now().UTC().Truncate(time.Microsecond)

	t.Run("lists what was addressed to the viewer, newest first, without what they caused", func(t *testing.T) {
		viewer, other := testActor(t, pool, "user"), testActor(t, pool, "user")
		older := notification.Notification{ID: uuid.New(), Audience: notification.ToUser(viewer.ID()), Message: "antiga", TicketID: uuid.New(), ActorID: other.ID(), CreatedAt: now}
		newer := notification.Notification{ID: uuid.New(), Audience: notification.ToUser(viewer.ID()), Message: "nova", TicketID: uuid.New(), ActorID: other.ID(), CreatedAt: now.Add(time.Minute)}
		ownAction := notification.Notification{ID: uuid.New(), Audience: notification.ToUser(viewer.ID()), Message: "própria", TicketID: uuid.New(), ActorID: viewer.ID(), CreatedAt: now}
		toSomeoneElse := notification.Notification{ID: uuid.New(), Audience: notification.ToUser(other.ID()), Message: "outra pessoa", TicketID: uuid.New(), ActorID: viewer.ID(), CreatedAt: now}
		require.NoError(t, store.Add(ctx, []notification.Notification{older, newer, ownAction, toSomeoneElse}))

		got, err := store.ListFor(ctx, viewer, 10)

		require.NoError(t, err)
		require.Len(t, got, 2)
		assert.Equal(t, "nova", got[0].Message)
		assert.Equal(t, newer, got[0])
		assert.Equal(t, "antiga", got[1].Message)
	})

	t.Run("reaches everyone with the addressed role", func(t *testing.T) {
		author := testActor(t, pool, "user")
		agent := testActor(t, pool, "support")
		broadcast := notification.Notification{ID: uuid.New(), Audience: notification.ToRole("support"), Message: "para o suporte", TicketID: uuid.New(), ActorID: author.ID(), CreatedAt: now}
		require.NoError(t, store.Add(ctx, []notification.Notification{broadcast}))
		t.Cleanup(func() { _, _ = pool.Exec(ctx, `DELETE FROM notifications WHERE id = $1`, broadcast.ID) })

		got, err := store.ListFor(ctx, agent, 50)

		require.NoError(t, err)
		assert.Contains(t, got, broadcast)
	})

	t.Run("keeps notifications that are not about a ticket", func(t *testing.T) {
		newcomer := testActor(t, pool, "user")
		admin := testActor(t, pool, "admin")
		signUp := notification.ForNewAccount(newcomer.ID(), "Maria Lima", now)
		require.NoError(t, store.Add(ctx, []notification.Notification{signUp}))

		got, err := store.ListFor(ctx, admin, 50)

		require.NoError(t, err)
		assert.Contains(t, got, signUp)
	})

	t.Run("remembers until when each user has seen their notifications", func(t *testing.T) {
		viewer := testActor(t, pool, "support")

		never, err := store.SeenUntil(ctx, viewer.ID())
		require.NoError(t, err)
		require.NoError(t, store.MarkSeen(ctx, viewer.ID(), now))
		require.NoError(t, store.MarkSeen(ctx, viewer.ID(), now.Add(time.Hour)))
		seen, err := store.SeenUntil(ctx, viewer.ID())

		require.NoError(t, err)
		assert.True(t, never.IsZero())
		assert.True(t, seen.Equal(now.Add(time.Hour)))
	})
}
