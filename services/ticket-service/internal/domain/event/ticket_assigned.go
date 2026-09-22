package event

import "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

type TicketAssigned struct {
	baseEvent
	AssigneeID string
}

func NewTicketAssigned(id *valueobjects.ID, assigneeID *valueobjects.AssigneeID) TicketAssigned {
	return TicketAssigned{
		baseEvent:  newBaseEvent(id.GetID()),
		AssigneeID: assigneeID.GetAssigneeID(),
	}
}

func (TicketAssigned) EventName() string {
	return "TicketAssigned"
}
