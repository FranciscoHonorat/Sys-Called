package event

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type eventDecoder func(payload []byte, base baseEvent) (Event, error)

var eventDecoders = map[string]eventDecoder{}

func registerEvent[T Event, PT interface {
	*T
	setBase(baseEvent)
}]() {
	var zero T
	eventDecoders[zero.EventName()] = func(payload []byte, base baseEvent) (Event, error) {
		var e T
		if err := json.Unmarshal(payload, &e); err != nil {
			return nil, err
		}
		PT(&e).setBase(base)
		return e, nil
	}
}

func Hydrate(eventType string, aggregateID uuid.UUID, occurredAt time.Time, payload []byte) (Event, error) {
	decode, ok := eventDecoders[eventType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domainErr.ErrUnknownEventType, eventType)
	}

	return decode(payload, baseEvent{aggregateID: aggregateID, occurredAt: occurredAt})
}
