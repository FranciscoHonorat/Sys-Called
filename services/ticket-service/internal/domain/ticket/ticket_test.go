package ticket_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func validTicketParts(t *testing.T) (*valueobjects.ID, *valueobjects.Title, *valueobjects.Description, valueobjects.Status, *valueobjects.AssigneeID, *valueobjects.Priority) {
	t.Helper()

	id := valueobjects.NewID(uuid.New())

	title, err := valueobjects.NewTitle("Valid Title")
	assert.NoError(t, err)

	description, err := valueobjects.NewDescription("Valid Description")
	assert.NoError(t, err)

	status := valueobjects.TicketStatusOpen

	assignee, err := valueobjects.NewAssigneeID("agent-1")
	assert.NoError(t, err)

	priority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityHigh))
	assert.NoError(t, err)

	return id, title, description, status, &assignee, &priority
}

func TestTicket(t *testing.T) {
	t.Run("should create a new Ticket with valid fields", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)

		assert.NoError(t, err)
		assert.Equal(t, id, tk.GetID())
		assert.Equal(t, title, tk.GetTitle())
		assert.Equal(t, description, tk.GetDescription())
		assert.Equal(t, status, tk.GetStatus())
		assert.Equal(t, assignee, tk.GetAssigneeID())
		assert.Equal(t, priority, tk.GetPriority())
		assert.WithinDuration(t, time.Now(), tk.GetCreatedAt(), time.Second)
	})

	t.Run("should create a new Ticket without an assignee or priority", func(t *testing.T) {
		id, title, description, status, _, _ := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, description, status, nil, nil)

		assert.NoError(t, err)
		assert.Nil(t, tk.GetAssigneeID())
		assert.Nil(t, tk.GetPriority())
	})

	t.Run("should return an error when id is nil", func(t *testing.T) {
		_, title, description, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(nil, title, description, status, assignee, priority)

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidID)
	})

	t.Run("should return an error when title is nil", func(t *testing.T) {
		id, _, description, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, nil, description, status, assignee, priority)

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error when description is nil", func(t *testing.T) {
		id, title, _, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, nil, status, assignee, priority)

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidDescription)
	})

	t.Run("should return an error when status is invalid", func(t *testing.T) {
		id, title, description, _, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, description, valueobjects.Status("bogus"), assignee, priority)

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidStatus)
	})

	t.Run("should marshal a Ticket with all fields to JSON", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		data, err := tk.MarshalJSON()

		assert.NoError(t, err)
		assert.Contains(t, string(data), `"id":"`+id.String()+`"`)
		assert.Contains(t, string(data), `"title":"Valid Title"`)
		assert.Contains(t, string(data), `"assignee_id":"agent-1"`)
	})

	t.Run("should omit assignee and priority when absent", func(t *testing.T) {
		id, title, description, status, _, _ := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, nil, nil)
		assert.NoError(t, err)

		data, err := tk.MarshalJSON()

		assert.NoError(t, err)
		assert.NotContains(t, string(data), `"assignee_id"`)
		assert.NotContains(t, string(data), `"priority"`)
	})

	t.Run("should round-trip marshal and unmarshal", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		original, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		data, err := original.MarshalJSON()
		assert.NoError(t, err)

		var decoded ticket.Ticket
		err = decoded.UnmarshalJSON(data)
		assert.NoError(t, err)

		assert.Equal(t, original.GetID().GetID(), decoded.GetID().GetID())
		assert.Equal(t, original.GetTitle().GetTitle(), decoded.GetTitle().GetTitle())
		assert.Equal(t, original.GetDescription().GetDescription(), decoded.GetDescription().GetDescription())
		assert.Equal(t, original.GetStatus(), decoded.GetStatus())
		assert.Equal(t, original.GetAssigneeID().GetAssigneeID(), decoded.GetAssigneeID().GetAssigneeID())
		assert.Equal(t, original.GetPriority().GetPriority(), decoded.GetPriority().GetPriority())
	})

	t.Run("should leave assignee and priority nil when absent from JSON", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","title":"T","description":"D","status":"Open","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var tk ticket.Ticket
		err := tk.UnmarshalJSON(data)

		assert.NoError(t, err)
		assert.Nil(t, tk.GetAssigneeID())
		assert.Nil(t, tk.GetPriority())
	})

	t.Run("should return an error for an invalid UUID", func(t *testing.T) {
		data := []byte(`{"id":"not-a-uuid","title":"T","description":"D","status":"Open","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var tk ticket.Ticket
		err := tk.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrInvalidUUID)
	})

	t.Run("should return an error for an empty title", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","title":"","description":"D","status":"Open","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var tk ticket.Ticket
		err := tk.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error for an invalid status", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","title":"T","description":"D","status":"bogus","created_at":"` + time.Now().Format(time.RFC3339) + `"}`)

		var tk ticket.Ticket
		err := tk.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrInvalidStatus)
	})

	t.Run("should return an error for a malformed created_at", func(t *testing.T) {
		data := []byte(`{"id":"` + uuid.New().String() + `","title":"T","description":"D","status":"Open","created_at":"not-a-date"}`)

		var tk ticket.Ticket
		err := tk.UnmarshalJSON(data)

		assert.Error(t, err)
	})
}
