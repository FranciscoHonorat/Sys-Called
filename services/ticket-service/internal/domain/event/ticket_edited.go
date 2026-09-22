package event

import (
	"encoding/json"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type TicketEdited struct {
	baseEvent
	Title       string
	Description string
}

func NewTicketEdited(id *valueobjects.ID, title *valueobjects.Title, description *valueobjects.Description) TicketEdited {
	return TicketEdited{
		baseEvent:   newBaseEvent(id.GetID()),
		Title:       title.GetTitle(),
		Description: description.GetDescription(),
	}
}

func (TicketEdited) EventName() string {
	return "TicketEdited"
}

func init() {
	registerEvent(TicketEdited{}.EventName(), func(payload []byte, base baseEvent) (Event, error) {
		var e TicketEdited
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	})
}
