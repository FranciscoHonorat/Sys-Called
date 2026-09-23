package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTicketOpened(t *testing.T) {
	t.Run("should carry the ticket data and the aggregate ID", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		assert.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		assert.NoError(t, err)
		priority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityHigh))
		assert.NoError(t, err)
		assignee, err := valueobjects.NewAssigneeID("agent-1")
		assert.NoError(t, err)

		e := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, &assignee, &priority, "user-1")

		assert.Equal(t, "TicketOpened", e.EventName())
		assert.Equal(t, id.GetID(), e.AggregateID())
		assert.WithinDuration(t, time.Now(), e.OccurredAt(), time.Second)
		assert.Equal(t, "Valid Title", e.Title)
		assert.Equal(t, "Valid Description", e.Description)
		assert.Equal(t, string(valueobjects.TicketStatusOpen), e.Status)
		assert.Equal(t, string(valueobjects.TicketPriorityHigh), e.Priority)
		assert.Equal(t, "agent-1", e.AssigneeID)
	})

	t.Run("should leave priority and assignee empty when absent", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("Valid Title")
		assert.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		assert.NoError(t, err)

		e := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, nil, nil, "user-1")

		assert.Empty(t, e.Priority)
		assert.Empty(t, e.AssigneeID)
	})
}
