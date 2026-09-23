package ticket

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

func (t *Ticket) IsVisibleTo(a actor.Actor) bool {
	switch {
	case a.IsAdmin():
		return true
	case a.IsSupport():
		return t.isOpen() || t.isAssignedTo(a.ID())
	default:
		return t.isRequestedBy(a.ID())
	}
}

func (t *Ticket) isOpen() bool {
	return t.status == valueobjects.TicketStatusOpen
}

func (t *Ticket) isAssignedTo(id string) bool {
	return t.assigneeID != nil && t.assigneeID.GetAssigneeID() == id
}

func (t *Ticket) isRequestedBy(id string) bool {
	return t.requesterID == id
}
