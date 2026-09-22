package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type AssignTicketInput struct {
	TicketID   string
	AssigneeID string
}

type AssignTicketUseCase struct {
	store repository.EventStore
}

func NewAssignTicketUseCase(store repository.EventStore) *AssignTicketUseCase {
	return &AssignTicketUseCase{store: store}
}

func (uc *AssignTicketUseCase) Execute(ctx context.Context, input AssignTicketInput) error {
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

	assigneeID, err := valueobjects.NewAssigneeID(input.AssigneeID)
	if err != nil {
		return err
	}

	if err := t.AssignTo(&assigneeID); err != nil {
		return err
	}

	if err := uc.store.Append(ctx, ticketID, t.GetUncommittedEvents(), len(history)); err != nil {
		return err
	}
	t.ClearUncommittedEvents()

	return nil
}
