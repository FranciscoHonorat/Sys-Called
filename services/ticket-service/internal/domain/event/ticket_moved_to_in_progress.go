package event

import "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

type TicketMovedToInProgress struct {
	baseEvent
}

func NewTicketMovedToInProgress(id *valueobjects.ID) TicketMovedToInProgress {
	return TicketMovedToInProgress{
		baseEvent: newBaseEvent(id.GetID()),
	}
}

func (TicketMovedToInProgress) EventName() string {
	return "TicketMovedToInProgress"
}
