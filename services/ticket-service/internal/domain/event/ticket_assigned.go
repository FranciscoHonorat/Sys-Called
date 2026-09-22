package event

import (
	"encoding/json"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

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

func init() {
	registerEvent(TicketAssigned{}.EventName(), func(payload []byte, base baseEvent) (Event, error) {
		var e TicketAssigned
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	})
}
