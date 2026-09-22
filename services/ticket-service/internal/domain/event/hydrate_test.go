package event_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHydrate(t *testing.T) {
	aggregateID := uuid.New()
	occurredAt := time.Now().Truncate(time.Second)

	t.Run("should hydrate a TicketOpened event", func(t *testing.T) {
		id := valueobjects.NewID(aggregateID)
		title, err := valueobjects.NewTitle("Valid Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("Valid Description")
		require.NoError(t, err)
		assignee, err := valueobjects.NewAssigneeID("agent-1")
		require.NoError(t, err)
		priority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityHigh))
		require.NoError(t, err)

		original := event.NewTicketOpened(id, title, description, valueobjects.TicketStatusOpen, &assignee, &priority)
		payload, err := json.Marshal(original)
		require.NoError(t, err)

		hydrated, err := event.Hydrate("TicketOpened", aggregateID, occurredAt, payload)

		assert.NoError(t, err)
		assert.Equal(t, aggregateID, hydrated.AggregateID())
		assert.Equal(t, occurredAt, hydrated.OccurredAt())
		opened, ok := hydrated.(event.TicketOpened)
		assert.True(t, ok)
		assert.Equal(t, original.Title, opened.Title)
		assert.Equal(t, original.Description, opened.Description)
		assert.Equal(t, original.Status, opened.Status)
		assert.Equal(t, original.AssigneeID, opened.AssigneeID)
		assert.Equal(t, original.Priority, opened.Priority)
	})

	t.Run("should hydrate a TicketEdited event", func(t *testing.T) {
		id := valueobjects.NewID(aggregateID)
		title, err := valueobjects.NewTitle("New Title")
		require.NoError(t, err)
		description, err := valueobjects.NewDescription("New Description")
		require.NoError(t, err)

		original := event.NewTicketEdited(id, title, description)
		payload, err := json.Marshal(original)
		require.NoError(t, err)

		hydrated, err := event.Hydrate("TicketEdited", aggregateID, occurredAt, payload)

		assert.NoError(t, err)
		edited, ok := hydrated.(event.TicketEdited)
		assert.True(t, ok)
		assert.Equal(t, original.Title, edited.Title)
		assert.Equal(t, original.Description, edited.Description)
	})

	t.Run("should hydrate a TicketAssigned event", func(t *testing.T) {
		id := valueobjects.NewID(aggregateID)
		assignee, err := valueobjects.NewAssigneeID("agent-1")
		require.NoError(t, err)

		original := event.NewTicketAssigned(id, &assignee)
		payload, err := json.Marshal(original)
		require.NoError(t, err)

		hydrated, err := event.Hydrate("TicketAssigned", aggregateID, occurredAt, payload)

		assert.NoError(t, err)
		assigned, ok := hydrated.(event.TicketAssigned)
		assert.True(t, ok)
		assert.Equal(t, original.AssigneeID, assigned.AssigneeID)
	})

	t.Run("should hydrate a TicketPriorityChanged event", func(t *testing.T) {
		id := valueobjects.NewID(aggregateID)
		priority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityLow))
		require.NoError(t, err)

		original := event.NewTicketPriorityChanged(id, &priority)
		payload, err := json.Marshal(original)
		require.NoError(t, err)

		hydrated, err := event.Hydrate("TicketPriorityChanged", aggregateID, occurredAt, payload)

		assert.NoError(t, err)
		changed, ok := hydrated.(event.TicketPriorityChanged)
		assert.True(t, ok)
		assert.Equal(t, original.Priority, changed.Priority)
	})

	t.Run("should hydrate a TicketMovedToInProgress event", func(t *testing.T) {
		id := valueobjects.NewID(aggregateID)
		original := event.NewTicketMovedToInProgress(id)
		payload, err := json.Marshal(original)
		require.NoError(t, err)

		hydrated, err := event.Hydrate("TicketMovedToInProgress", aggregateID, occurredAt, payload)

		assert.NoError(t, err)
		_, ok := hydrated.(event.TicketMovedToInProgress)
		assert.True(t, ok)
	})

	t.Run("should hydrate a TicketClosed event", func(t *testing.T) {
		id := valueobjects.NewID(aggregateID)
		original := event.NewTicketClosed(id)
		payload, err := json.Marshal(original)
		require.NoError(t, err)

		hydrated, err := event.Hydrate("TicketClosed", aggregateID, occurredAt, payload)

		assert.NoError(t, err)
		_, ok := hydrated.(event.TicketClosed)
		assert.True(t, ok)
	})

	t.Run("should hydrate a TicketResponseAdded event", func(t *testing.T) {
		ticketID := valueobjects.NewID(aggregateID)
		author, err := valueobjects.NewAuthorID("agent-1")
		require.NoError(t, err)
		content, err := valueobjects.NewContent("Valid response content")
		require.NoError(t, err)
		r, err := response.NewResponse(valueobjects.NewID(uuid.New()), ticketID, &author, content)
		require.NoError(t, err)

		original := event.NewTicketResponseAdded(ticketID, r)
		payload, err := json.Marshal(original)
		require.NoError(t, err)

		hydrated, err := event.Hydrate("TicketResponseAdded", aggregateID, occurredAt, payload)

		assert.NoError(t, err)
		added, ok := hydrated.(event.TicketResponseAdded)
		assert.True(t, ok)
		assert.Equal(t, original.ResponseID, added.ResponseID)
		assert.Equal(t, original.AuthorID, added.AuthorID)
		assert.Equal(t, original.Content, added.Content)
	})

	t.Run("should return an error for an unknown event type", func(t *testing.T) {
		_, err := event.Hydrate("SomethingElse", aggregateID, occurredAt, []byte(`{}`))

		assert.ErrorIs(t, err, domainErr.ErrUnknownEventType)
	})
}
