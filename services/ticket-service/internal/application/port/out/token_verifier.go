package out

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

type TokenVerifier interface {
	Verify(ctx context.Context, token string) (actor.Actor, error)
}
