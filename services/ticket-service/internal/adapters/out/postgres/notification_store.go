package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
)

var _ out.NotificationStore = (*NotificationStore)(nil)

type NotificationStore struct {
	pool *pgxpool.Pool
}

func NewNotificationStore(pool *pgxpool.Pool) *NotificationStore {
	return &NotificationStore{pool: pool}
}

func (s *NotificationStore) Add(ctx context.Context, notifications []notification.Notification) error {
	batch := &pgx.Batch{}
	for _, n := range notifications {
		batch.Queue(
			`INSERT INTO notifications (id, audience_kind, audience, message, ticket_id, actor_id, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			n.ID, string(n.Audience.Kind), n.Audience.Value, n.Message, n.TicketID, n.ActorID, n.CreatedAt,
		)
	}
	return s.pool.SendBatch(ctx, batch).Close()
}

func (s *NotificationStore) ListFor(ctx context.Context, viewer actor.Actor, limit int) ([]notification.Notification, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, audience_kind, audience, message, ticket_id, actor_id, created_at
		 FROM notifications
		 WHERE ((audience_kind = 'user' AND audience = $1) OR (audience_kind = 'role' AND audience = $2))
		   AND actor_id <> $1
		 ORDER BY created_at DESC
		 LIMIT $3`,
		viewer.ID(), string(viewer.Role()), limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []notification.Notification
	for rows.Next() {
		var n notification.Notification
		var kind string
		if err := rows.Scan(&n.ID, &kind, &n.Audience.Value, &n.Message, &n.TicketID, &n.ActorID, &n.CreatedAt); err != nil {
			return nil, err
		}
		n.Audience.Kind = notification.AudienceKind(kind)
		n.CreatedAt = n.CreatedAt.UTC()
		notifications = append(notifications, n)
	}
	return notifications, rows.Err()
}

func (s *NotificationStore) SeenUntil(ctx context.Context, userID string) (time.Time, error) {
	var seen time.Time
	err := s.pool.QueryRow(ctx, `SELECT seen_until FROM notification_reads WHERE user_id = $1`, userID).Scan(&seen)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, nil
	}
	return seen, err
}

func (s *NotificationStore) MarkSeen(ctx context.Context, userID string, at time.Time) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO notification_reads (user_id, seen_until) VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET seen_until = EXCLUDED.seen_until`,
		userID, at,
	)
	return err
}
