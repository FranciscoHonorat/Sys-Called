package httpapi

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/repository"
)

type Handler struct {
	openTicket             *application.OpenTicketUseCase
	getTicket              *application.GetTicketUseCase
	listTickets            *application.ListTicketsUseCase
	assignTicket           *application.AssignTicketUseCase
	changeTicketPriority   *application.ChangeTicketPriorityUseCase
	moveTicketToInProgress *application.MoveTicketToInProgressUseCase
	closeTicket            *application.CloseTicketUseCase
	addTicketResponse      *application.AddTicketResponseUseCase
}

func NewHandler(store repository.EventStore, cache repository.TicketCache) *Handler {
	return &Handler{
		openTicket:             application.NewOpenTicketUseCase(store, cache),
		getTicket:              application.NewGetTicketUseCase(store, cache),
		listTickets:            application.NewListTicketsUseCase(store, cache),
		assignTicket:           application.NewAssignTicketUseCase(store, cache),
		changeTicketPriority:   application.NewChangeTicketPriorityUseCase(store, cache),
		moveTicketToInProgress: application.NewMoveTicketToInProgressUseCase(store, cache),
		closeTicket:            application.NewCloseTicketUseCase(store, cache),
		addTicketResponse:      application.NewAddTicketResponseUseCase(store, cache),
	}
}
