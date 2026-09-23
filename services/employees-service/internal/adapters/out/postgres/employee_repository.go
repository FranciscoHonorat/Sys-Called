package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

var _ out.EmployeeRepository = (*EmployeeRepository)(nil)

const (
	uniqueViolation = "23505"
	employeeColumns = `id, name, COALESCE(username, ''), COALESCE(role, ''), COALESCE(password_hash, ''),
		status, must_change_password, password_reset_requested`
)

type EmployeeRepository struct {
	pool *pgxpool.Pool
}

func NewEmployeeRepository(pool *pgxpool.Pool) *EmployeeRepository {
	return &EmployeeRepository{pool: pool}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanEmployee(row rowScanner) (employee.Employee, error) {
	var s employee.Snapshot
	var role, status string
	if err := row.Scan(&s.ID, &s.Name, &s.Username, &role, &s.PasswordHash, &status, &s.MustChangePassword, &s.PasswordResetRequested); err != nil {
		return employee.Employee{}, err
	}
	s.Role, s.Status = employee.Role(role), employee.Status(status)
	return employee.Rehydrate(s), nil
}

func (r *EmployeeRepository) List(ctx context.Context) ([]employee.Employee, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+employeeColumns+` FROM employees ORDER BY registered_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []employee.Employee
	for rows.Next() {
		e, err := scanEmployee(rows)
		if err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}
	return employees, rows.Err()
}

func (r *EmployeeRepository) FindByUsername(ctx context.Context, username string) (employee.Employee, bool, error) {
	return r.findOne(ctx, `WHERE username = $1`, username)
}

func (r *EmployeeRepository) FindByID(ctx context.Context, id string) (employee.Employee, bool, error) {
	return r.findOne(ctx, `WHERE id = $1`, id)
}

func (r *EmployeeRepository) findOne(ctx context.Context, where string, arg string) (employee.Employee, bool, error) {
	e, err := scanEmployee(r.pool.QueryRow(ctx, `SELECT `+employeeColumns+` FROM employees `+where, arg))
	if errors.Is(err, pgx.ErrNoRows) {
		return employee.Employee{}, false, nil
	}
	if err != nil {
		return employee.Employee{}, false, err
	}
	return e, true, nil
}

func (r *EmployeeRepository) Register(ctx context.Context, e employee.Employee) error {
	return r.inTransaction(ctx, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx,
			`INSERT INTO employees (id, name, username, role, password_hash, status, must_change_password, password_reset_requested, registered_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			 ON CONFLICT (id) DO UPDATE SET username = EXCLUDED.username, role = EXCLUDED.role, password_hash = EXCLUDED.password_hash
			 WHERE employees.username IS NULL`,
			e.GetID(), e.GetName(), e.GetUsername(), e.GetRole().String(), e.GetPasswordHash(),
			string(e.GetStatus()), e.MustChangePassword(), e.HasRequestedPasswordReset(), time.Now(),
		)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation {
			return domainErr.ErrUsernameTaken
		}
		if err != nil || tag.RowsAffected() == 0 {
			return err
		}
		return appendToOutbox(ctx, tx, e.Events())
	})
}

func (r *EmployeeRepository) Save(ctx context.Context, e employee.Employee) error {
	return r.inTransaction(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx,
			`UPDATE employees SET password_hash = $2, status = $3, must_change_password = $4, password_reset_requested = $5
			 WHERE id = $1`,
			e.GetID(), e.GetPasswordHash(), string(e.GetStatus()), e.MustChangePassword(), e.HasRequestedPasswordReset(),
		); err != nil {
			return err
		}
		return appendToOutbox(ctx, tx, e.Events())
	})
}

func (r *EmployeeRepository) inTransaction(ctx context.Context, work func(tx pgx.Tx) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := work(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func appendToOutbox(ctx context.Context, tx pgx.Tx, events []employee.Event) error {
	for _, evt := range events {
		payload, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(ctx,
			`INSERT INTO outbox_events (id, event_type, payload, occurred_at) VALUES ($1, $2, $3, $4)`,
			uuid.New(), evt.EventType(), payload, time.Now(),
		); err != nil {
			return err
		}
	}
	return nil
}
