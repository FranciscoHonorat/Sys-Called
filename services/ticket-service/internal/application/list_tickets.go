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
	ids, err := uc.store.ListAggregateIDs(ctx)
	if err != nil {
		return nil, err
	}

	outputs := make([]GetTicketOutput, 0, len(ids))
	for _, id := range ids {
		t, _, err := uc.loadTicket(ctx, id)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, toGetTicketOutput(t))
	}

	return outputs, nil
}
