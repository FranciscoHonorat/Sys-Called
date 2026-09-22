package event

import (
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
	registerEvent[TicketEdited]()
}
