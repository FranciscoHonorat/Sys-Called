package application

import (
	"context"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type AutoAssignTicketInput struct {
	TicketID string
}

type AutoAssignTicketOutput struct {
	AssigneeID string `json:"assignee_id"`
}

type AutoAssignTicketUseCase struct {
	eventSourcedUseCase
	responsibles repository.ResponsibleDirectory
}

func NewAutoAssignTicketUseCase(store repository.EventStore, cache repository.TicketCache, responsibles repository.ResponsibleDirectory) *AutoAssignTicketUseCase {
	return &AutoAssignTicketUseCase{
		eventSourcedUseCase: eventSourcedUseCase{store: store, cache: cache},
		responsibles:        responsibles,
	}
}

func (uc *AutoAssignTicketUseCase) Execute(ctx context.Context, input AutoAssignTicketInput) (AutoAssignTicketOutput, error) {
	ticketID, err := uuid.Parse(input.TicketID)
	if err != nil {
		return AutoAssignTicketOutput{}, domainErr.ErrInvalidUUID
	}

	t, version, err := uc.loadTicket(ctx, ticketID)
	if err != nil {
		return AutoAssignTicketOutput{}, err
	}

	candidates, err := uc.responsibles.List(ctx)
	if err != nil {
		return AutoAssignTicketOutput{}, err
	}
	if len(candidates) == 0 {
		return AutoAssignTicketOutput{}, domainErr.ErrNoResponsiblesAvailable
	}

	allTickets, err := uc.loadAllTickets(ctx)
	if err != nil {
		return AutoAssignTicketOutput{}, err
	}

	chosenID := leastBusyResponsible(candidates, countOpenTicketsByAssignee(allTickets))

	assigneeID, err := valueobjects.NewAssigneeID(chosenID)
	if err != nil {
		return AutoAssignTicketOutput{}, err
	}

	if err := t.AssignTo(&assigneeID); err != nil {
		return AutoAssignTicketOutput{}, err
	}

	if err := uc.commit(ctx, t, ticketID, version); err != nil {
		return AutoAssignTicketOutput{}, err
	}

	return AutoAssignTicketOutput{AssigneeID: chosenID}, nil
}

func countOpenTicketsByAssignee(tickets []*ticket.Ticket) map[string]int {
	counts := make(map[string]int)
	for _, t := range tickets {
		if t.GetStatus() == valueobjects.TicketStatusClosed {
			continue
		}
		if t.GetAssigneeID() == nil {
			continue
		}
		counts[t.GetAssigneeID().GetAssigneeID()]++
	}
	return counts
}

func leastBusyResponsible(candidates []string, openCounts map[string]int) string {
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if openCounts[candidate] < openCounts[best] {
			best = candidate
		}
	}
	return best
}
