package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type AddTicketResponseInput struct {
	TicketID string
	AuthorID string
	Content  string
}

type AddTicketResponseUseCase struct {
	application.EventSourcedUseCase
}

func NewAddTicketResponseUseCase(store out.EventStore, cache out.TicketCache) *AddTicketResponseUseCase {
	return &AddTicketResponseUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *AddTicketResponseUseCase) Execute(ctx context.Context, input AddTicketResponseInput) error {
	return uc.UpdateTicket(ctx, input.TicketID, func(t *ticket.Ticket) error {
		authorID, err := valueobjects.NewAuthorID(input.AuthorID)
		if err != nil {
			return err
		}

		content, err := valueobjects.NewContent(input.Content)
		if err != nil {
			return err
		}

		r, err := response.NewResponse(valueobjects.NewID(uuid.Nil), t.GetID(), &authorID, content)
		if err != nil {
			return err
		}

		return t.AddResponse(r)
	})
}
