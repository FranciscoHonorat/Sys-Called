package event

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type TicketClosed struct {
	baseEvent
	Resolution string
}

func NewTicketClosed(id *valueobjects.ID, resolution string) TicketClosed {
	return TicketClosed{
		baseEvent:  newBaseEvent(id.GetID()),
		Resolution: resolution,
	}
}

func (TicketClosed) EventName() string {
	return "TicketClosed"
}

func init() {
	registerEvent[TicketClosed]()
}
