package ticket_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/response"
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

func validResponseFor(t *testing.T, ticketID *valueobjects.ID) *response.Response {
	t.Helper()

	author, err := valueobjects.NewAuthorID("agent-1")
	assert.NoError(t, err)

	content, err := valueobjects.NewContent("Valid response content")
	assert.NoError(t, err)

	r, err := response.NewResponse(valueobjects.NewID(uuid.Nil), ticketID, &author, content)
	assert.NoError(t, err)

	return r
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

	t.Run("should include responses in the marshaled JSON", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, tk.AddResponse(validResponseFor(t, tk.GetID())))

		data, err := tk.MarshalJSON()

		assert.NoError(t, err)
		assert.Contains(t, string(data), `"content":"Valid response content"`)
	})

	t.Run("should omit assignee, priority and responses when absent", func(t *testing.T) {
		id, title, description, status, _, _ := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, nil, nil)
		assert.NoError(t, err)

		data, err := tk.MarshalJSON()

		assert.NoError(t, err)
		assert.NotContains(t, string(data), `"assignee_id"`)
		assert.NotContains(t, string(data), `"priority"`)
		assert.NotContains(t, string(data), `"responses"`)
	})

	t.Run("should round-trip marshal and unmarshal", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		original, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, original.AddResponse(validResponseFor(t, original.GetID())))

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
		assert.Len(t, decoded.GetResponses(), 1)
		assert.Equal(t, original.GetResponses()[0].GetContent().GetContent(), decoded.GetResponses()[0].GetContent().GetContent())
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

	t.Run("should return an error when unmarshaling a response for a different ticket", func(t *testing.T) {
		ticketID := uuid.New()
		otherTicketID := uuid.New()
		data := []byte(`{"id":"` + ticketID.String() + `","title":"T","description":"D","status":"Open","created_at":"` + time.Now().Format(time.RFC3339) +
			`","responses":[{"id":"` + uuid.New().String() + `","ticket_id":"` + otherTicketID.String() + `","author_id":"agent-1","content":"C","created_at":"` + time.Now().Format(time.RFC3339) + `"}]}`)

		var tk ticket.Ticket
		err := tk.UnmarshalJSON(data)

		assert.ErrorIs(t, err, domainErr.ErrResponseTicketMismatch)
	})

	t.Run("should assign a Ticket to a new assignee", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		newAssignee, err := valueobjects.NewAssigneeID("agent-2")
		assert.NoError(t, err)

		err = tk.AssignTo(&newAssignee)

		assert.NoError(t, err)
		assert.Equal(t, &newAssignee, tk.GetAssigneeID())
	})

	t.Run("should return an error when assigning a nil AssigneeID", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		err = tk.AssignTo(nil)

		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
	})

	t.Run("should return an error when assigning a closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, tk.Close())

		err = tk.AssignTo(assignee)

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should change the priority of a Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		newPriority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityLow))
		assert.NoError(t, err)

		err = tk.ChangePriority(&newPriority)

		assert.NoError(t, err)
		assert.Equal(t, &newPriority, tk.GetPriority())
	})

	t.Run("should return an error when changing to a nil Priority", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		err = tk.ChangePriority(nil)

		assert.ErrorIs(t, err, domainErr.ErrInvalidPriority)
	})

	t.Run("should return an error when changing the priority of a closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, tk.Close())

		err = tk.ChangePriority(priority)

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should move an open Ticket to in progress", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		err = tk.MoveToInProgress()

		assert.NoError(t, err)
		assert.Equal(t, valueobjects.TicketStatusInProgress, tk.GetStatus())
	})

	t.Run("should return an error when moving a non-open Ticket to in progress", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, tk.MoveToInProgress())

		err = tk.MoveToInProgress()

		assert.ErrorIs(t, err, domainErr.ErrInvalidStatusTransition)
	})

	t.Run("should close an open Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		err = tk.Close()

		assert.NoError(t, err)
		assert.Equal(t, valueobjects.TicketStatusClosed, tk.GetStatus())
	})

	t.Run("should close a Ticket that is in progress", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, tk.MoveToInProgress())

		err = tk.Close()

		assert.NoError(t, err)
		assert.Equal(t, valueobjects.TicketStatusClosed, tk.GetStatus())
	})

	t.Run("should return an error when closing an already closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, tk.Close())

		err = tk.Close()

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should add a response to a Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		r := validResponseFor(t, tk.GetID())

		err = tk.AddResponse(r)

		assert.NoError(t, err)
		assert.Equal(t, []*response.Response{r}, tk.GetResponses())
	})

	t.Run("should return an error when adding a nil response", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		err = tk.AddResponse(nil)

		assert.ErrorIs(t, err, domainErr.ErrInvalidResponse)
	})

	t.Run("should return an error when the response belongs to a different Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)

		otherTicketID := valueobjects.NewID(uuid.New())
		r := validResponseFor(t, otherTicketID)

		err = tk.AddResponse(r)

		assert.ErrorIs(t, err, domainErr.ErrResponseTicketMismatch)
	})

	t.Run("should return an error when adding a response to a closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority)
		assert.NoError(t, err)
		assert.NoError(t, tk.Close())

		r := validResponseFor(t, tk.GetID())

		err = tk.AddResponse(r)

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})
}
