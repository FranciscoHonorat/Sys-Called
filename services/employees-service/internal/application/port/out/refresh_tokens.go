package out

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/session"
)

type RefreshTokenGenerator interface {
	Generate() (string, error)
	Hash(token string) string
}

type RefreshTokenStore interface {
	Save(ctx context.Context, token session.RefreshToken) error
	FindByHash(ctx context.Context, hash string) (session.RefreshToken, bool, error)
	Revoke(ctx context.Context, hash string) (bool, error)
	RevokeAllForEmployee(ctx context.Context, employeeID string) error
}
