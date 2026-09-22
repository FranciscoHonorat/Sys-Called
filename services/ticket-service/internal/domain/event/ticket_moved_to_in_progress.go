package event

import (
	"encoding/json"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

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

func init() {
	registerEvent(TicketMovedToInProgress{}.EventName(), func(payload []byte, base baseEvent) (Event, error) {
		var e TicketMovedToInProgress
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	})
}
