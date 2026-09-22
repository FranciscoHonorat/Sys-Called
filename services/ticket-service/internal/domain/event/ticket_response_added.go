package event

import (
	"encoding/json"

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

func init() {
	registerEvent(TicketResponseAdded{}.EventName(), func(payload []byte, base baseEvent) (Event, error) {
		var e TicketResponseAdded
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	})
}
