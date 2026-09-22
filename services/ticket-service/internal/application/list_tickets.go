package application

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
)

type ListTicketsUseCase struct {
	eventSourcedUseCase
}

func NewListTicketsUseCase(store repository.EventStore, cache repository.TicketCache) *ListTicketsUseCase {
	return &ListTicketsUseCase{eventSourcedUseCase{store: store, cache: cache}}
}

func (uc *ListTicketsUseCase) Execute(ctx context.Context) ([]GetTicketOutput, error) {
	tickets, err := uc.loadAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]GetTicketOutput, 0, len(tickets))
	for _, t := range tickets {
		outputs = append(outputs, toGetTicketOutput(t))
	}

	return outputs, nil
}
