package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

type EmployeeRepository struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{pool: pool}
}

func (r *EmployeeRepository) List(ctx context.Context) ([]employee.Employee, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, name FROM employees ORDER BY registered_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []employee.Employee
	for rows.Next() {
		var id, name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		employees = append(employees, employee.NewEmployee(id, name))
	}

	return employees, rows.Err()
}

func (r *EmployeeRepository) Register(ctx context.Context, e employee.Employee) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx,
		`INSERT INTO employees (id, name, registered_at) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING`,
		e.GetID(), e.GetName(), time.Now(),
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return tx.Commit(ctx)
	}

	payload, err := json.Marshal(struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}{ID: e.GetID(), Name: e.GetName()})
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`INSERT INTO outbox_events (id, event_type, payload, occurred_at) VALUES ($1, $2, $3, $4)`,
		uuid.New(), "EmployeeRegistered", payload, time.Now(),
	); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
