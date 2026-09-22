package command

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
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
	application.EventSourcedUseCase
	responsibles out.ResponsibleDirectory
}

func NewAutoAssignTicketUseCase(store out.EventStore, cache out.TicketCache, responsibles out.ResponsibleDirectory) *AutoAssignTicketUseCase {
	return &AutoAssignTicketUseCase{
		EventSourcedUseCase: application.NewEventSourcedUseCase(store, cache),
		responsibles:        responsibles,
	}
}

func (uc *AutoAssignTicketUseCase) Execute(ctx context.Context, input AutoAssignTicketInput) (AutoAssignTicketOutput, error) {
	var chosenID string
	err := uc.UpdateTicket(ctx, input.TicketID, func(t *ticket.Ticket) error {
		candidates, err := uc.responsibles.List(ctx)
		if err != nil {
			return err
		}
		if len(candidates) == 0 {
			return domainErr.ErrNoResponsiblesAvailable
		}

		allTickets, err := uc.LoadAllTickets(ctx)
		if err != nil {
			return err
		}

		chosenID = leastBusyResponsible(candidates, countOpenTicketsByAssignee(allTickets))

		assigneeID, err := valueobjects.NewAssigneeID(chosenID)
		if err != nil {
			return err
		}

		return t.AssignTo(&assigneeID)
	})
	if err != nil {
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
