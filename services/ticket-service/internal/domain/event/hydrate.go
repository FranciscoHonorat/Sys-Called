package event

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

type eventDecoder func(payload []byte, base baseEvent) (Event, error)

var eventDecoders = map[string]eventDecoder{}

func registerEvent(name string, decode eventDecoder) {
	eventDecoders[name] = decode
}

func Hydrate(eventType string, aggregateID uuid.UUID, occurredAt time.Time, payload []byte) (Event, error) {
	decode, ok := eventDecoders[eventType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", domainErr.ErrUnknownEventType, eventType)
	}

	return decode(payload, baseEvent{aggregateID: aggregateID, occurredAt: occurredAt})
}
