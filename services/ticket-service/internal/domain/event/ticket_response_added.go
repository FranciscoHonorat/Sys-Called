package event

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"
)

type TicketResponseAdded struct {
	baseEvent
	ResponseID string
	AuthorID   string
	Content    string
}

func NewTicketResponseAdded(id *valueobjects.ID, r *response.Response) TicketResponseAdded {
	return TicketResponseAdded{
		baseEvent:  newBaseEvent(id.GetID()),
		ResponseID: r.GetID().String(),
		AuthorID:   r.GetAuthorID().GetAuthorID(),
		Content:    r.GetContent().GetContent(),
	}
}

func (TicketResponseAdded) EventName() string {
	return "TicketResponseAdded"
}
