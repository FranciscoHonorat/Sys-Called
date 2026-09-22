package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTicketPriorityChanged(t *testing.T) {
	t.Run("should carry the new priority and the aggregate ID", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		priority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityLow))
		assert.NoError(t, err)

		e := event.NewTicketPriorityChanged(id, &priority)

		assert.Equal(t, "TicketPriorityChanged", e.EventName())
		assert.Equal(t, id.GetID(), e.AggregateID())
		assert.WithinDuration(t, time.Now(), e.OccurredAt(), time.Second)
		assert.Equal(t, string(valueobjects.TicketPriorityLow), e.Priority)
	})
}
