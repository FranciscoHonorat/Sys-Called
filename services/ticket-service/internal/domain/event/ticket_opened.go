package event

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type TicketOpened struct {
	baseEvent
	Title       string
	Description string
	Status      string
	AssigneeID  string
	Priority    string
}

func NewTicketOpened(id *valueobjects.ID, title *valueobjects.Title, description *valueobjects.Description, status valueobjects.Status, assigneeID *valueobjects.AssigneeID, priority *valueobjects.Priority) TicketOpened {
	e := TicketOpened{
		baseEvent:   newBaseEvent(id.GetID()),
		Title:       title.GetTitle(),
		Description: description.GetDescription(),
		Status:      status.String(),
	}
	if assigneeID != nil {
		e.AssigneeID = assigneeID.GetAssigneeID()
	}
	if priority != nil {
		e.Priority = priority.GetPriority()
	}
	return e
}

func (TicketOpened) EventName() string {
	return "TicketOpened"
}

func init() {
	registerEvent[TicketOpened]()
}
