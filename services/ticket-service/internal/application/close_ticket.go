package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type CloseTicketInput struct {
	TicketID string
}

type CloseTicketUseCase struct {
	store repository.EventStore
}

func NewCloseTicketUseCase(store repository.EventStore) *CloseTicketUseCase {
	return &CloseTicketUseCase{store: store}
}

func (uc *CloseTicketUseCase) Execute(ctx context.Context, input CloseTicketInput) error {
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

	if err := t.Close(); err != nil {
		return err
	}

	if err := uc.store.Append(ctx, ticketID, t.GetUncommittedEvents(), len(history)); err != nil {
		return err
	}
	t.ClearUncommittedEvents()

	return nil
}
