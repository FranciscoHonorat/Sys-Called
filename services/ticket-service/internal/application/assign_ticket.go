package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type AssignTicketInput struct {
	TicketID   string
	AssigneeID string
}

type AssignTicketUseCase struct {
	eventSourcedUseCase
}

func NewAssignTicketUseCase(store repository.EventStore, cache repository.TicketCache) *AssignTicketUseCase {
	return &AssignTicketUseCase{eventSourcedUseCase{store: store, cache: cache}}
}

func (uc *AssignTicketUseCase) Execute(ctx context.Context, input AssignTicketInput) error {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return err
	}

	assigneeID, err := valueobjects.NewAssigneeID(input.AssigneeID)
	if err != nil {
		return err
	}

	if err := t.AssignTo(&assigneeID); err != nil {
		return err
	}

	return uc.commit(ctx, t, ticketID, version)
}
