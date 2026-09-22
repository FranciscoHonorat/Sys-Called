package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type EditTicketInput struct {
	TicketID    string
	Title       string
	Description string
}

type EditTicketUseCase struct {
	application.EventSourcedUseCase
}

func NewEditTicketUseCase(store out.EventStore, cache out.TicketCache) *EditTicketUseCase {
	return &EditTicketUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *EditTicketUseCase) Execute(ctx context.Context, input EditTicketInput) error {
	return uc.UpdateTicket(ctx, input.TicketID, func(t *ticket.Ticket) error {
		title, err := valueobjects.NewTitle(input.Title)
		if err != nil {
			return err
		}

		description, err := valueobjects.NewDescription(input.Description)
		if err != nil {
			return err
		}

		return t.Edit(title, description)
	})
}
