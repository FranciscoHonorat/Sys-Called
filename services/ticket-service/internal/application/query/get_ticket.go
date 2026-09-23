package query

import (
	"context"
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type GetTicketInput struct {
	Viewer   actor.Actor
	TicketID string
}

type GetTicketResponseOutput struct {
	ResponseID string    `json:"response_id"`
	AuthorID   string    `json:"author_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetTicketOutput struct {
	TicketID    string                    `json:"ticket_id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Status      string                    `json:"status"`
	AssigneeID  string                    `json:"assignee_id,omitempty"`
	Priority    string                    `json:"priority,omitempty"`
	RequesterID string                    `json:"requester_id,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
	ClosedAt    *time.Time                `json:"closed_at,omitempty"`
	Resolution  string                    `json:"resolution,omitempty"`
	Responses   []GetTicketResponseOutput `json:"responses,omitempty"`
}

type GetTicketUseCase struct {
	application.EventSourcedUseCase
}

func NewGetTicketUseCase(store out.EventStore, cache out.TicketCache) *GetTicketUseCase {
	return &GetTicketUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *GetTicketUseCase) Execute(ctx context.Context, input GetTicketInput) (GetTicketOutput, error) {
	t, err := uc.LoadCachedTicket(ctx, input.TicketID)
	if err != nil {
		return GetTicketOutput{}, err
	}
	if !t.IsVisibleTo(input.Viewer) {
		return GetTicketOutput{}, domainErr.ErrEventStreamNotFound
	}

	return toGetTicketOutput(t), nil
}

func toGetTicketOutput(t *ticket.Ticket) GetTicketOutput {
	output := GetTicketOutput{
		TicketID:    t.GetID().String(),
		Title:       t.GetTitle().GetTitle(),
		Description: t.GetDescription().GetDescription(),
		Status:      t.GetStatus().String(),
		CreatedAt:   t.GetCreatedAt(),
		RequesterID: t.GetRequesterID(),
		ClosedAt:    t.GetClosedAt(),
		Resolution:  t.GetResolution(),
	}

	if t.GetAssigneeID() != nil {
		output.AssigneeID = t.GetAssigneeID().GetAssigneeID()
	}
	if t.GetPriority() != nil {
		output.Priority = t.GetPriority().GetPriority()
	}

	for _, r := range t.GetResponses() {
		output.Responses = append(output.Responses, GetTicketResponseOutput{
			ResponseID: r.GetID().String(),
			AuthorID:   r.GetAuthorID().GetAuthorID(),
			Content:    r.GetContent().GetContent(),
			CreatedAt:  r.GetCreatedAt(),
		})
	}

	return output
}
