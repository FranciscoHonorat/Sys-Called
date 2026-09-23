package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type OpenTicketInput struct {
	Actor       actor.Actor
	Title       string
	Description string
	AssigneeID  string
	Priority    string
}

type OpenTicketOutput struct {
	TicketID string `json:"ticket_id"`
}

type OpenTicketUseCase struct {
	application.EventSourcedUseCase
}

func NewOpenTicketUseCase(store out.EventStore, cache out.TicketCache) *OpenTicketUseCase {
	return &OpenTicketUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *OpenTicketUseCase) Execute(ctx context.Context, input OpenTicketInput) (OpenTicketOutput, error) {
	title, err := valueobjects.NewTitle(input.Title)
	if err != nil {
		return OpenTicketOutput{}, err
	}

	description, err := valueobjects.NewDescription(input.Description)
	if err != nil {
		return OpenTicketOutput{}, err
	}

	var assigneeID *valueobjects.AssigneeID
	if input.AssigneeID != "" {
		if !ticket.CanManage(nil, input.Actor) {
			return OpenTicketOutput{}, domainErr.ErrForbidden
		}
		a, err := valueobjects.NewAssigneeID(input.AssigneeID)
		if err != nil {
			return OpenTicketOutput{}, err
		}
		assigneeID = &a
	}

	var priority *valueobjects.Priority
	if input.Priority != "" {
		p, err := valueobjects.NewPriority(input.Priority)
		if err != nil {
			return OpenTicketOutput{}, err
		}
		priority = &p
	}

	id := valueobjects.NewID(uuid.Nil)

	t, err := ticket.NewTicket(id, title, description, valueobjects.TicketStatusOpen, assigneeID, priority, input.Actor.ID())
	if err != nil {
		return OpenTicketOutput{}, err
	}

	if err := uc.CreateTicket(ctx, input.Actor, t); err != nil {
		return OpenTicketOutput{}, err
	}

	return OpenTicketOutput{TicketID: t.GetID().String()}, nil
}
