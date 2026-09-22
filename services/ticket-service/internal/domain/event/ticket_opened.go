package event

import "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

type TicketOpened struct {
	baseEvent
	Title       string
	Description string
	Priority    string
	AssigneeID  string
}

func NewTicketOpened(id *valueobjects.ID, title *valueobjects.Title, description *valueobjects.Description, priority *valueobjects.Priority, assigneeID *valueobjects.AssigneeID) TicketOpened {
	e := TicketOpened{
		baseEvent:   newBaseEvent(id.GetID()),
		Title:       title.GetTitle(),
		Description: description.GetDescription(),
	}
	if priority != nil {
		e.Priority = priority.GetPriority()
	}
	if assigneeID != nil {
		e.AssigneeID = assigneeID.GetAssigneeID()
	}
	return e
}

func (TicketOpened) EventName() string {
	return "TicketOpened"
}
