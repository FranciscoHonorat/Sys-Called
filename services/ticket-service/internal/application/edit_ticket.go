package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type EditTicketInput struct {
	TicketID    string
	Title       string
	Description string
}

type EditTicketUseCase struct {
	eventSourcedUseCase
}

func NewEditTicketUseCase(store repository.EventStore, cache repository.TicketCache) *EditTicketUseCase {
	return &EditTicketUseCase{eventSourcedUseCase{store: store, cache: cache}}
}

func (uc *EditTicketUseCase) Execute(ctx context.Context, input EditTicketInput) error {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return err
	}

	title, err := valueobjects.NewTitle(input.Title)
	if err != nil {
		return err
	}

	description, err := valueobjects.NewDescription(input.Description)
	if err != nil {
		return err
	}

	if err := t.Edit(title, description); err != nil {
		return err
	}

	return uc.commit(ctx, t, ticketID, version)
}
