package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type ChangeTicketPriorityInput struct {
	TicketID string
	Priority string
}

type ChangeTicketPriorityUseCase struct {
	eventSourcedUseCase
}

func NewChangeTicketPriorityUseCase(store repository.EventStore, cache repository.TicketCache) *ChangeTicketPriorityUseCase {
	return &ChangeTicketPriorityUseCase{eventSourcedUseCase{store: store, cache: cache}}
}

func (uc *ChangeTicketPriorityUseCase) Execute(ctx context.Context, input ChangeTicketPriorityInput) error {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
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

	return uc.commit(ctx, t, ticketID, version)
}
