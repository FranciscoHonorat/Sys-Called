package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type AddTicketResponseInput struct {
	TicketID string
	AuthorID string
	Content  string
}

type AddTicketResponseUseCase struct {
	eventSourcedUseCase
}

func NewAddTicketResponseUseCase(store repository.EventStore, cache repository.TicketCache) *AddTicketResponseUseCase {
	return &AddTicketResponseUseCase{eventSourcedUseCase{store: store, cache: cache}}
}

func (uc *AddTicketResponseUseCase) Execute(ctx context.Context, input AddTicketResponseInput) error {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return domainErr.ErrInvalidUUID
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return err
	}

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

	if err := t.AddResponse(r); err != nil {
		return err
	}

	return uc.commit(ctx, t, ticketID, version)
}
