package httpapi

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"
)

type UseCases struct {
	OpenTicket             *command.OpenTicketUseCase
	GetTicket              *query.GetTicketUseCase
	ListTickets            *query.ListTicketsUseCase
	EditTicket             *command.EditTicketUseCase
	AssignTicket           *command.AssignTicketUseCase
	AutoAssignTicket       *command.AutoAssignTicketUseCase
	ChangeTicketPriority   *command.ChangeTicketPriorityUseCase
	MoveTicketToInProgress *command.MoveTicketToInProgressUseCase
	CloseTicket            *command.CloseTicketUseCase
	AddTicketResponse      *command.AddTicketResponseUseCase
}

type Handler struct {
	useCases UseCases
}

func NewHandler(useCases UseCases) *Handler {
	return &Handler{useCases: useCases}
}
