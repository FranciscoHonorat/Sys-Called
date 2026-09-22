package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
)

const uniqueViolationCode = "23505"

type EventStore struct {
	pool *pgxpool.Pool
}

func NewEventStore(pool *pgxpool.Pool) *EventStore {
	return &EventStore{pool: pool}
}

func (s *EventStore) Append(ctx context.Context, aggregateID uuid.UUID, events []event.Event, expectedVersion int) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var currentVersion int
	err = tx.QueryRow(ctx, `SELECT COUNT(*) FROM ticket_events WHERE aggregate_id = $1`, aggregateID).Scan(&currentVersion)
	if err != nil {
		return err
	}
	if currentVersion != expectedVersion {
		return domainErr.ErrConcurrencyConflict
	}

	for i, e := range events {
		payload, err := json.Marshal(e)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO ticket_events (id, aggregate_id, version, event_type, payload, occurred_at) VALUES ($1, $2, $3, $4, $5, $6)`,
			uuid.New(), aggregateID, expectedVersion+i+1, e.EventName(), payload, e.OccurredAt(),
		)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
				return domainErr.ErrConcurrencyConflict
			}
			return err
		}
	}

	return tx.Commit(ctx)
}

func (s *EventStore) Load(ctx context.Context, aggregateID uuid.UUID) ([]event.Event, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT event_type, payload, occurred_at FROM ticket_events WHERE aggregate_id = $1 ORDER BY version ASC`,
		aggregateID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []event.Event
	for rows.Next() {
		var eventType string
		var payload []byte
		var occurredAt time.Time
		if err := rows.Scan(&eventType, &payload, &occurredAt); err != nil {
			return nil, err
		}

		e, err := event.Hydrate(eventType, aggregateID, occurredAt, payload)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, domainErr.ErrEventStreamNotFound
	}

	return events, nil
}

func (s *EventStore) ListAggregateIDs(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT aggregate_id FROM ticket_events WHERE version = 1 ORDER BY occurred_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
