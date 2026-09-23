package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type CloseTicketInput struct {
	Actor      actor.Actor
	TicketID   string
	Resolution string
}

type CloseTicketUseCase struct {
	application.EventSourcedUseCase
}

func NewCloseTicketUseCase(store out.EventStore, cache out.TicketCache) *CloseTicketUseCase {
	return &CloseTicketUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *CloseTicketUseCase) Execute(ctx context.Context, input CloseTicketInput) error {
	return uc.UpdateTicket(ctx, input.Actor, input.TicketID, ticket.CanWork, func(t *ticket.Ticket) error {
		return t.Close(input.Resolution)
	})
}
