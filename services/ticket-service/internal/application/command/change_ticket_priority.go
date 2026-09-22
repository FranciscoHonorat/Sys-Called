package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type ChangeTicketPriorityInput struct {
	TicketID string
	Priority string
}

type ChangeTicketPriorityUseCase struct {
	application.EventSourcedUseCase
}

func NewChangeTicketPriorityUseCase(store out.EventStore, cache out.TicketCache) *ChangeTicketPriorityUseCase {
	return &ChangeTicketPriorityUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *ChangeTicketPriorityUseCase) Execute(ctx context.Context, input ChangeTicketPriorityInput) error {
	return uc.UpdateTicket(ctx, input.TicketID, func(t *ticket.Ticket) error {
		priority, err := valueobjects.NewPriority(input.Priority)
		if err != nil {
			return err
		}

		return t.ChangePriority(&priority)
	})
}
