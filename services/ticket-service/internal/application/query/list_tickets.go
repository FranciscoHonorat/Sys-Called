package query

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
)

type ListTicketsUseCase struct {
	application.EventSourcedUseCase
}

func NewListTicketsUseCase(store out.EventStore, cache out.TicketCache) *ListTicketsUseCase {
	return &ListTicketsUseCase{application.NewEventSourcedUseCase(store, cache)}
}

func (uc *ListTicketsUseCase) Execute(ctx context.Context) ([]GetTicketOutput, error) {
	tickets, err := uc.LoadAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]GetTicketOutput, 0, len(tickets))
	for _, t := range tickets {
		outputs = append(outputs, toGetTicketOutput(t))
	}

	return outputs, nil
}
