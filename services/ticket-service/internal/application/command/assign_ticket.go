package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type AssignTicketInput struct {
	Actor      actor.Actor
	TicketID   string
	AssigneeID string
}

type AssignTicketUseCase struct {
	application.EventSourcedUseCase
}

func NewAssignTicketUseCase(store out.EventStore, cache out.TicketCache) *AssignTicketUseCase {
	return &AssignTicketUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *AssignTicketUseCase) Execute(ctx context.Context, input AssignTicketInput) error {
	return uc.UpdateTicket(ctx, input.Actor, input.TicketID, ticket.CanManage, func(t *ticket.Ticket) error {
		assigneeID, err := valueobjects.NewAssigneeID(input.AssigneeID)
		if err != nil {
			return err
		}

		return t.AssignTo(&assigneeID)
	})
}
