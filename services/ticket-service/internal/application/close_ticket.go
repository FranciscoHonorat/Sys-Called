package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
)

type CloseTicketInput struct {
	TicketID string
}

type CloseTicketUseCase struct {
	eventSourcedUseCase
}

func NewCloseTicketUseCase(store repository.EventStore, cache repository.TicketCache) *CloseTicketUseCase {
	return &CloseTicketUseCase{eventSourcedUseCase{store: store, cache: cache}}
}

func (uc *CloseTicketUseCase) Execute(ctx context.Context, input CloseTicketInput) error {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return err
	}

	if err := t.Close(); err != nil {
		return err
	}

	return uc.commit(ctx, t, ticketID, version)
}
