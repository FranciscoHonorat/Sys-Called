package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type OpenTicketInput struct {
	Title       string
	Description string
	AssigneeID  string
	Priority    string
}

type OpenTicketOutput struct {
	TicketID string
}

type OpenTicketUseCase struct {
	store repository.EventStore
}

func NewOpenTicketUseCase(store repository.EventStore) *OpenTicketUseCase {
	return &OpenTicketUseCase{store: store}
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

	t, err := ticket.NewTicket(id, title, description, valueobjects.TicketStatusOpen, assigneeID, priority)
	if err != nil {
		return OpenTicketOutput{}, err
	}

	if err := uc.store.Append(ctx, t.GetID().GetID(), t.GetUncommittedEvents(), 0); err != nil {
		return OpenTicketOutput{}, err
	}
	t.ClearUncommittedEvents()

	return OpenTicketOutput{TicketID: t.GetID().String()}, nil
}
