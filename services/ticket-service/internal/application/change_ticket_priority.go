package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type ChangeTicketPriorityInput struct {
	TicketID string
	Priority string
}

type ChangeTicketPriorityUseCase struct {
	store repository.EventStore
	cache repository.TicketCache
}

func NewChangeTicketPriorityUseCase(store repository.EventStore, cache repository.TicketCache) *ChangeTicketPriorityUseCase {
	return &ChangeTicketPriorityUseCase{store: store, cache: cache}
}

func (uc *ChangeTicketPriorityUseCase) Execute(ctx context.Context, input ChangeTicketPriorityInput) error {
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

	priority, err := valueobjects.NewPriority(input.Priority)
	if err != nil {
		return err
	}

	if err := t.ChangePriority(&priority); err != nil {
		return err
	}

	if err := uc.store.Append(ctx, ticketID, t.GetUncommittedEvents(), len(history)); err != nil {
		return err
	}
	t.ClearUncommittedEvents()
	uc.cache.Set(ctx, t)

	return nil
}
