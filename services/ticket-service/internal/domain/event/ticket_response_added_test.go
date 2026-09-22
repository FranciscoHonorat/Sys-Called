package event_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestTicketResponseAdded(t *testing.T) {
	t.Run("should carry the response data and the aggregate ID", func(t *testing.T) {
		ticketID := valueobjects.NewID(uuid.New())

		responseID := valueobjects.NewID(uuid.New())
		author, err := valueobjects.NewAuthorID("agent-1")
		assert.NoError(t, err)
		content, err := valueobjects.NewContent("Valid response content")
		assert.NoError(t, err)

		r, err := response.NewResponse(responseID, ticketID, &author, content)
		assert.NoError(t, err)

		e := event.NewTicketResponseAdded(ticketID, r)

		assert.Equal(t, "TicketResponseAdded", e.EventName())
		assert.Equal(t, ticketID.GetID(), e.AggregateID())
		assert.WithinDuration(t, time.Now(), e.OccurredAt(), time.Second)
		assert.Equal(t, responseID.String(), e.ResponseID)
		assert.Equal(t, "agent-1", e.AuthorID)
		assert.Equal(t, "Valid response content", e.Content)
	})
}
