package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTicketAssigned(t *testing.T) {
	t.Run("should carry the assignee and the aggregate ID", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		assignee, err := valueobjects.NewAssigneeID("agent-1")
		assert.NoError(t, err)

		e := event.NewTicketAssigned(id, &assignee)

		assert.Equal(t, "TicketAssigned", e.EventName())
		assert.Equal(t, id.GetID(), e.AggregateID())
		assert.WithinDuration(t, time.Now(), e.OccurredAt(), time.Second)
		assert.Equal(t, "agent-1", e.AssigneeID)
	})
}
