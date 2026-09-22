package httpapi

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
)

type Handler struct {
	openTicket             *application.OpenTicketUseCase
	getTicket              *application.GetTicketUseCase
	listTickets            *application.ListTicketsUseCase
	editTicket             *application.EditTicketUseCase
	assignTicket           *application.AssignTicketUseCase
	autoAssignTicket       *application.AutoAssignTicketUseCase
	changeTicketPriority   *application.ChangeTicketPriorityUseCase
	moveTicketToInProgress *application.MoveTicketToInProgressUseCase
	closeTicket            *application.CloseTicketUseCase
	addTicketResponse      *application.AddTicketResponseUseCase
}

func NewHandler(store repository.EventStore, cache repository.TicketCache, responsibles repository.ResponsibleDirectory) *Handler {
	return &Handler{
		openTicket:             application.NewOpenTicketUseCase(store, cache),
		getTicket:              application.NewGetTicketUseCase(store, cache),
		listTickets:            application.NewListTicketsUseCase(store, cache),
		editTicket:             application.NewEditTicketUseCase(store, cache),
		assignTicket:           application.NewAssignTicketUseCase(store, cache),
		autoAssignTicket:       application.NewAutoAssignTicketUseCase(store, cache, responsibles),
		changeTicketPriority:   application.NewChangeTicketPriorityUseCase(store, cache),
		moveTicketToInProgress: application.NewMoveTicketToInProgressUseCase(store, cache),
		closeTicket:            application.NewCloseTicketUseCase(store, cache),
		addTicketResponse:      application.NewAddTicketResponseUseCase(store, cache),
	}
}
