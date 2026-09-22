package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
)

type OutboxStore struct {
	pool *pgxpool.Pool
}

func NewOutboxStore(pool *pgxpool.Pool) *OutboxStore {
	return &OutboxStore{pool: pool}
}

func (s *OutboxStore) FetchPending(ctx context.Context) ([]out.OutboxEvent, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, event_type, payload FROM outbox_events WHERE published_at IS NULL ORDER BY occurred_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []out.OutboxEvent
	for rows.Next() {
		var e out.OutboxEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, rows.Err()
}

func (s *OutboxStore) MarkPublished(ctx context.Context, id uuid.UUID) error {
	_, err := s.pool.Exec(ctx, `UPDATE outbox_events SET published_at = now() WHERE id = $1`, id)
	return err
}
