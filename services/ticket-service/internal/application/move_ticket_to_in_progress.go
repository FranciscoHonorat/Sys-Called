package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
)

type MoveTicketToInProgressInput struct {
	TicketID string
}

type MoveTicketToInProgressUseCase struct {
	eventSourcedUseCase
}

func NewMoveTicketToInProgressUseCase(store repository.EventStore, cache repository.TicketCache) *MoveTicketToInProgressUseCase {
	return &MoveTicketToInProgressUseCase{eventSourcedUseCase{store: store, cache: cache}}
}

func (uc *MoveTicketToInProgressUseCase) Execute(ctx context.Context, input MoveTicketToInProgressInput) error {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return err
	}

	if err := t.MoveToInProgress(); err != nil {
		return err
	}

	return uc.commit(ctx, t, ticketID, version)
}
