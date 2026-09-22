package event

import (
	"encoding/json"

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
	registerEvent(TicketClosed{}.EventName(), func(payload []byte, base baseEvent) (Event, error) {
		var e TicketClosed
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	})
}
