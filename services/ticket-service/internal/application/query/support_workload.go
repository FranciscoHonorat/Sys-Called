package query

import (
	"context"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type AgentWorkload struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Open       int    `json:"open"`
	InProgress int    `json:"in_progress"`
	Closed     int    `json:"closed"`
}

type SupportWorkloadUseCase struct {
	application.EventSourcedUseCase
	directory out.ResponsibleDirectory
}

func NewSupportWorkloadUseCase(store out.EventStore, cache out.TicketCache, directory out.ResponsibleDirectory) *SupportWorkloadUseCase {
	return &SupportWorkloadUseCase{EventSourcedUseCase: application.NewEventSourcedUseCase(store, cache), directory: directory}
}

func (uc *SupportWorkloadUseCase) Execute(ctx context.Context, viewer actor.Actor) ([]AgentWorkload, error) {
	if !viewer.IsAdmin() {
		return nil, domainErr.ErrForbidden
	}

	agents, err := uc.directory.ListWithNames(ctx)
	if err != nil {
		return nil, err
	}
	tickets, err := uc.LoadAllTickets(ctx)
	if err != nil {
		return nil, err
	}

	byAgent := make(map[string]*AgentWorkload, len(agents))
	workloads := make([]AgentWorkload, len(agents))
	for i, agent := range agents {
		workloads[i] = AgentWorkload{ID: agent.ID, Name: agent.Name}
		byAgent[agent.ID] = &workloads[i]
	}
	for _, t := range tickets {
		if t.GetAssigneeID() == nil {
			continue
		}
		workload, ok := byAgent[t.GetAssigneeID().GetAssigneeID()]
		if !ok {
			continue
		}
		workload.count(t.GetStatus())
	}
	return workloads, nil
}

func (w *AgentWorkload) count(status valueobjects.Status) {
	switch status {
	case valueobjects.TicketStatusOpen:
		w.Open++
	case valueobjects.TicketStatusInProgress:
		w.InProgress++
	case valueobjects.TicketStatusClosed:
		w.Closed++
	}
}
