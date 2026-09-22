package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type eventSourcedUseCase struct {
	store repository.EventStore
	cache repository.TicketCache
}

func (uc eventSourcedUseCase) loadTicket(ctx context.Context, ticketID uuid.UUID) (*ticket.Ticket, int, error) {
	history, err := uc.store.Load(ctx, ticketID)
	if err != nil {
		return nil, 0, err
	}

	t, err := ticket.LoadFromHistory(history)
	if err != nil {
		return nil, 0, err
	}

	return t, len(history), nil
}

func (uc eventSourcedUseCase) loadAllTickets(ctx context.Context) ([]*ticket.Ticket, error) {
	ids, err := uc.store.ListAggregateIDs(ctx)
	if err != nil {
		return nil, err
	}

	tickets := make([]*ticket.Ticket, 0, len(ids))
	for _, id := range ids {
		t, _, err := uc.loadTicket(ctx, id)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, t)
	}

	return tickets, nil
}

func (uc eventSourcedUseCase) commit(ctx context.Context, t *ticket.Ticket, ticketID uuid.UUID, expectedVersion int) error {
	if err := uc.store.Append(ctx, ticketID, t.GetUncommittedEvents(), expectedVersion); err != nil {
		return err
	}
	t.ClearUncommittedEvents()
	uc.cache.Set(ctx, t)
	return nil
}
