package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func Hydrate(eventType string, aggregateID uuid.UUID, occurredAt time.Time, payload []byte) (Event, error) {
	base := baseEvent{aggregateID: aggregateID, occurredAt: occurredAt}

	switch eventType {
	case "TicketOpened":
		var e TicketOpened
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	case "TicketAssigned":
		var e TicketAssigned
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	case "TicketPriorityChanged":
		var e TicketPriorityChanged
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	case "TicketMovedToInProgress":
		var e TicketMovedToInProgress
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	case "TicketClosed":
		var e TicketClosed
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	case "TicketResponseAdded":
		var e TicketResponseAdded
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		e.baseEvent = base
		return e, nil
	default:
		return nil, fmt.Errorf("%w: %s", domainErr.ErrUnknownEventType, eventType)
	}
}
