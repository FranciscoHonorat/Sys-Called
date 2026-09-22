package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ResponsibleDirectory struct {
	pool *pgxpool.Pool
}

func NewResponsibleDirectory(pool *pgxpool.Pool) *ResponsibleDirectory {
	return &ResponsibleDirectory{pool: pool}
}

func (d *ResponsibleDirectory) List(ctx context.Context) ([]string, error) {
	rows, err := d.pool.Query(ctx, `SELECT id FROM responsibles ORDER BY registered_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (d *ResponsibleDirectory) Upsert(ctx context.Context, id, name string) error {
	_, err := d.pool.Exec(ctx,
		`INSERT INTO responsibles (id, name) VALUES ($1, $2)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`,
		id, name,
	)
	return err
}
