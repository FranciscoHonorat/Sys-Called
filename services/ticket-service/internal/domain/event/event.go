package event

import (
	"time"

	"github.com/google/uuid"
)

type Event interface {
	EventName() string
	AggregateID() uuid.UUID
	OccurredAt() time.Time
}

type baseEvent struct {
	aggregateID uuid.UUID
	occurredAt  time.Time
}

func newBaseEvent(aggregateID uuid.UUID) baseEvent {
	return baseEvent{
		aggregateID: aggregateID,
		occurredAt:  time.Now(),
	}
}

func (e baseEvent) AggregateID() uuid.UUID {
	return e.aggregateID
}

func (e baseEvent) OccurredAt() time.Time {
	return e.occurredAt
}
