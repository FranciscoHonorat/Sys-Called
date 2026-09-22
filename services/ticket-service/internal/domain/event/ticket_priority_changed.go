package event

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type TicketPriorityChanged struct {
	baseEvent
	Priority string
}

func NewTicketPriorityChanged(id *valueobjects.ID, priority *valueobjects.Priority) TicketPriorityChanged {
	return TicketPriorityChanged{
		baseEvent: newBaseEvent(id.GetID()),
		Priority:  priority.GetPriority(),
	}
}

func (TicketPriorityChanged) EventName() string {
	return "TicketPriorityChanged"
}

func init() {
	registerEvent[TicketPriorityChanged]()
}
