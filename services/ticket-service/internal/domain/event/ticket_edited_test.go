package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTicketEdited(t *testing.T) {
	t.Run("should carry the new title, description and the aggregate ID", func(t *testing.T) {
		id := valueobjects.NewID(uuid.New())
		title, err := valueobjects.NewTitle("New Title")
		assert.NoError(t, err)
		description, err := valueobjects.NewDescription("New Description")
		assert.NoError(t, err)

		e := event.NewTicketEdited(id, title, description)

		assert.Equal(t, "TicketEdited", e.EventName())
		assert.Equal(t, id.GetID(), e.AggregateID())
		assert.WithinDuration(t, time.Now(), e.OccurredAt(), time.Second)
		assert.Equal(t, "New Title", e.Title)
		assert.Equal(t, "New Description", e.Description)
	})
}
