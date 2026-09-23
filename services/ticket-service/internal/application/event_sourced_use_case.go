package application

import (
	"context"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type EventSourcedUseCase struct {
	store out.EventStore
	cache out.TicketCache
}

func NewEventSourcedUseCase(store out.EventStore, cache out.TicketCache) EventSourcedUseCase {
	return EventSourcedUseCase{store: store, cache: cache}
}

func (uc EventSourcedUseCase) CreateTicket(ctx context.Context, by actor.Actor, t *ticket.Ticket) error {
	if err := requireActor(by); err != nil {
		return err
	}
	return uc.commit(ctx, by, t, t.GetID().GetID(), 0)
}

func (uc EventSourcedUseCase) UpdateTicket(ctx context.Context, by actor.Actor, rawTicketID string, allowed ticket.Permission, change func(t *ticket.Ticket) error) error {
	if err := requireActor(by); err != nil {
		return err
	}

	ticketID, err := parseTicketID(rawTicketID)
	if err != nil {
		return err
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return err
	}
	if err := authorize(t, by, allowed); err != nil {
		return err
	}

	if err := change(t); err != nil {
		return err
	}

	return uc.commit(ctx, by, t, ticketID, version)
}

func (uc EventSourcedUseCase) loadTicket(ctx context.Context, ticketID uuid.UUID) (*ticket.Ticket, int, error) {
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

func (uc EventSourcedUseCase) LoadCachedTicket(ctx context.Context, rawTicketID string) (*ticket.Ticket, error) {
	ticketID, err := parseTicketID(rawTicketID)
	if err != nil {
		return nil, err
	}

	if t, ok := uc.cache.Get(ctx, ticketID); ok {
		return t, nil
	}

	t, _, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	uc.cache.Set(ctx, t)
	return t, nil
}

func (uc EventSourcedUseCase) LoadAllTickets(ctx context.Context) ([]*ticket.Ticket, error) {
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

func (uc EventSourcedUseCase) commit(ctx context.Context, by actor.Actor, t *ticket.Ticket, ticketID uuid.UUID, expectedVersion int) error {
	if err := uc.store.Append(ctx, ticketID, t.GetUncommittedEvents(), expectedVersion, by.ID()); err != nil {
		return err
	}
	t.ClearUncommittedEvents()
	uc.cache.Set(ctx, t)
	return nil
}

func parseTicketID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, domainErr.ErrInvalidUUID
	}
	return id, nil
}

func requireActor(by actor.Actor) error {
	if by.ID() == "" {
		return domainErr.ErrInvalidActor
	}
	return nil
}

func authorize(t *ticket.Ticket, by actor.Actor, allowed ticket.Permission) error {
	if !t.IsVisibleTo(by) {
		return domainErr.ErrEventStreamNotFound
	}
	if !allowed(t, by) {
		return domainErr.ErrForbidden
	}
	return nil
}
