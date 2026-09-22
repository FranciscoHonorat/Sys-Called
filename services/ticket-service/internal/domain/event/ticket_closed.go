package event

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type TicketClosed struct {
	baseEvent
}

func NewTicketClosed(id *valueobjects.ID) TicketClosed {
	return TicketClosed{
		baseEvent: newBaseEvent(id.GetID()),
	}
}

func (TicketClosed) EventName() string {
	return "TicketClosed"
}

func init() {
	registerEvent[TicketClosed]()
}
