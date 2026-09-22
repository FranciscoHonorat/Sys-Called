package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTicketClosed(t *testing.T) {
	t.Run("should carry the aggregate ID", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())

		e := event.NewTicketClosed(id)

		assert.Equal(t, "TicketClosed", e.EventName())
		assert.Equal(t, id.GetID(), e.AggregateID())
		assert.WithinDuration(t, time.Now(), e.OccurredAt(), time.Second)
	})
}
