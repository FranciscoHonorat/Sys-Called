package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type MoveTicketToInProgressInput struct {
	TicketID string
}

type MoveTicketToInProgressUseCase struct {
	store repository.EventStore
	cache repository.TicketCache
}

func NewMoveTicketToInProgressUseCase(store repository.EventStore, cache repository.TicketCache) *MoveTicketToInProgressUseCase {
	return &MoveTicketToInProgressUseCase{store: store, cache: cache}
}

func (uc *MoveTicketToInProgressUseCase) Execute(ctx context.Context, input MoveTicketToInProgressInput) error {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}

	history, err := uc.store.Load(ctx, ticketID)
	if err != nil {
		return err
	}

	t, err := ticket.LoadFromHistory(history)
	if err != nil {
		return err
	}

	if err := t.MoveToInProgress(); err != nil {
		return err
	}

	if err := uc.store.Append(ctx, ticketID, t.GetUncommittedEvents(), len(history)); err != nil {
		return err
	}
	t.ClearUncommittedEvents()
	uc.cache.Set(ctx, t)

	return nil
}
