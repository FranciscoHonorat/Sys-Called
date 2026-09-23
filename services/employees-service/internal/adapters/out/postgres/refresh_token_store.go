package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/session"
)

var _ out.RefreshTokenStore = (*RefreshTokenStore)(nil)

type RefreshTokenStore struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenStore(pool *pgxpool.Pool) *RefreshTokenStore {
	return &RefreshTokenStore{pool: pool}
}

func (s *RefreshTokenStore) Save(ctx context.Context, token session.RefreshToken) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (token_hash, employee_id, expires_at) VALUES ($1, $2, $3)`,
		token.Hash(), token.EmployeeID(), token.ExpiresAt(),
	)
	return err
}

func (s *RefreshTokenStore) FindByHash(ctx context.Context, hash string) (session.RefreshToken, bool, error) {
	var employeeID string
	var expiresAt time.Time
	var revoked bool
	err := s.pool.QueryRow(ctx,
		`SELECT employee_id, expires_at, revoked_at IS NOT NULL FROM refresh_tokens WHERE token_hash = $1`, hash,
	).Scan(&employeeID, &expiresAt, &revoked)
	if errors.Is(err, pgx.ErrNoRows) {
		return session.RefreshToken{}, false, nil
	}
	if err != nil {
		return session.RefreshToken{}, false, err
	}
	return session.RestoreRefreshToken(hash, employeeID, expiresAt, revoked), true, nil
}

func (s *RefreshTokenStore) Revoke(ctx context.Context, hash string) (bool, error) {
	tag, err := s.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = now() WHERE token_hash = $1 AND revoked_at IS NULL`, hash,
	)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

func (s *RefreshTokenStore) RevokeAllForEmployee(ctx context.Context, employeeID string) error {
	_, err := s.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = now() WHERE employee_id = $1 AND revoked_at IS NULL`, employeeID,
	)
	return err
}
