package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type MoveTicketToInProgressInput struct {
	TicketID string
}

type MoveTicketToInProgressUseCase struct {
	application.EventSourcedUseCase
}

func NewMoveTicketToInProgressUseCase(store out.EventStore, cache out.TicketCache) *MoveTicketToInProgressUseCase {
	return &MoveTicketToInProgressUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *MoveTicketToInProgressUseCase) Execute(ctx context.Context, input MoveTicketToInProgressInput) error {
	return uc.UpdateTicket(ctx, input.TicketID, func(t *ticket.Ticket) error {
		return t.MoveToInProgress()
	})
}
