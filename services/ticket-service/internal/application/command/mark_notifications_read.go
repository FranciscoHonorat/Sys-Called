package command

import (
	"context"
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

type MarkNotificationsReadUseCase struct {
	store out.NotificationStore
	now   func() time.Time
}

func NewMarkNotificationsReadUseCase(store out.NotificationStore, now func() time.Time) *MarkNotificationsReadUseCase {
	return &MarkNotificationsReadUseCase{store: store, now: now}
}

func (uc *MarkNotificationsReadUseCase) Execute(ctx context.Context, viewer actor.Actor) error {
	return uc.store.MarkSeen(ctx, viewer.ID(), uc.now())
}
