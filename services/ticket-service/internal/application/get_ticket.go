package application

import (
	"context"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type GetTicketInput struct {
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
	CreatedAt   time.Time                 `json:"created_at"`
	Responses   []GetTicketResponseOutput `json:"responses,omitempty"`
}

type GetTicketUseCase struct {
	store repository.EventStore
	cache repository.TicketCache
}

func NewGetTicketUseCase(store repository.EventStore, cache repository.TicketCache) *GetTicketUseCase {
	return &GetTicketUseCase{store: store, cache: cache}
}

func (uc *GetTicketUseCase) Execute(ctx context.Context, input GetTicketInput) (GetTicketOutput, error) {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return GetTicketOutput{}, domainErr.ErrInvalidUUID
	}

	if t, ok := uc.cache.Get(ctx, ticketID); ok {
		return toGetTicketOutput(t), nil
	}

	history, err := uc.store.Load(ctx, ticketID)
	if err != nil {
		return GetTicketOutput{}, err
	}

	t, err := ticket.LoadFromHistory(history)
	if err != nil {
		return GetTicketOutput{}, err
	}

	uc.cache.Set(ctx, t)

	return toGetTicketOutput(t), nil
}

func toGetTicketOutput(t *ticket.Ticket) GetTicketOutput {
	output := GetTicketOutput{
		TicketID:    t.GetID().String(),
		Title:       t.GetTitle().GetTitle(),
		Description: t.GetDescription().GetDescription(),
		Status:      t.GetStatus().String(),
		CreatedAt:   t.GetCreatedAt(),
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
