package ticket_test

import (
	"testing"
	"time"

	"github.com/google/uuid"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
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

		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")

		assert.NoError(t, err)
		assert.Equal(t, id, tk.GetID())
		assert.Equal(t, title, tk.GetTitle())
		assert.Equal(t, description, tk.GetDescription())
		assert.Equal(t, status, tk.GetStatus())
		assert.Equal(t, assignee, tk.GetAssigneeID())
		assert.Equal(t, priority, tk.GetPriority())
		assert.WithinDuration(t, time.Now(), tk.GetCreatedAt(), time.Second)
		assert.Len(t, tk.GetUncommittedEvents(), 1)
		assert.Equal(t, "TicketOpened", tk.GetUncommittedEvents()[0].EventName())
	})

	t.Run("should create a new Ticket without an assignee or priority", func(t *testing.T) {
		id, title, description, status, _, _ := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, description, status, nil, nil, "user-1")

		assert.NoError(t, err)
		assert.Nil(t, tk.GetAssigneeID())
		assert.Nil(t, tk.GetPriority())
	})

	t.Run("should record who opened the ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")

		assert.NoError(t, err)
		assert.Equal(t, "user-1", tk.GetRequesterID())
		opened, ok := tk.GetUncommittedEvents()[0].(event.TicketOpened)
		assert.True(t, ok)
		assert.Equal(t, "user-1", opened.RequesterID)
	})

	t.Run("should require someone who opened the ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)

		_, err := ticket.NewTicket(id, title, description, status, assignee, priority, "")

		assert.ErrorIs(t, err, domainErr.ErrInvalidRequester)
	})

	t.Run("should restore who opened the ticket from its history", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		opened, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		tk, err := ticket.LoadFromHistory(opened.GetUncommittedEvents())

		assert.NoError(t, err)
		assert.Equal(t, "user-1", tk.GetRequesterID())
	})

	t.Run("should return an error when id is nil", func(t *testing.T) {
		_, title, description, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(nil, title, description, status, assignee, priority, "user-1")

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidID)
	})

	t.Run("should return an error when title is nil", func(t *testing.T) {
		id, _, description, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, nil, description, status, assignee, priority, "user-1")

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error when description is nil", func(t *testing.T) {
		id, title, _, status, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, nil, status, assignee, priority, "user-1")

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidDescription)
	})

	t.Run("should return an error when status is invalid", func(t *testing.T) {
		id, title, description, _, assignee, priority := validTicketParts(t)

		tk, err := ticket.NewTicket(id, title, description, valueobjects.Status("bogus"), assignee, priority, "user-1")

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidStatus)
	})

	t.Run("should edit the title and description of a Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		newTitle, err := valueobjects.NewTitle("New Title")
		assert.NoError(t, err)
		newDescription, err := valueobjects.NewDescription("New Description")
		assert.NoError(t, err)

		err = tk.Edit(newTitle, newDescription)

		assert.NoError(t, err)
		assert.Equal(t, newTitle, tk.GetTitle())
		assert.Equal(t, newDescription, tk.GetDescription())
		events := tk.GetUncommittedEvents()
		assert.Equal(t, "TicketEdited", events[len(events)-1].EventName())
	})

	t.Run("should return an error when editing with a nil title", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.Edit(nil, description)

		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error when editing with a nil description", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.Edit(title, nil)

		assert.ErrorIs(t, err, domainErr.ErrInvalidDescription)
	})

	t.Run("should return an error when editing a closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NoError(t, tk.Close("Resolvido"))

		err = tk.Edit(title, description)

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should assign a Ticket to a new assignee", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		newAssignee, err := valueobjects.NewAssigneeID("agent-2")
		assert.NoError(t, err)

		err = tk.AssignTo(&newAssignee)

		assert.NoError(t, err)
		assert.Equal(t, &newAssignee, tk.GetAssigneeID())
		events := tk.GetUncommittedEvents()
		assert.Equal(t, "TicketAssigned", events[len(events)-1].EventName())
	})

	t.Run("should return an error when assigning a nil AssigneeID", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.AssignTo(nil)

		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
	})

	t.Run("should return an error when assigning a closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NoError(t, tk.Close("Resolvido"))

		err = tk.AssignTo(assignee)

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should change the priority of a Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		newPriority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityLow))
		assert.NoError(t, err)

		err = tk.ChangePriority(&newPriority)

		assert.NoError(t, err)
		assert.Equal(t, &newPriority, tk.GetPriority())
		events := tk.GetUncommittedEvents()
		assert.Equal(t, "TicketPriorityChanged", events[len(events)-1].EventName())
	})

	t.Run("should return an error when changing to a nil Priority", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.ChangePriority(nil)

		assert.ErrorIs(t, err, domainErr.ErrInvalidPriority)
	})

	t.Run("should return an error when changing the priority of a closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NoError(t, tk.Close("Resolvido"))

		err = tk.ChangePriority(priority)

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should move an open Ticket to in progress", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.MoveToInProgress()

		assert.NoError(t, err)
		assert.Equal(t, valueobjects.TicketStatusInProgress, tk.GetStatus())
		events := tk.GetUncommittedEvents()
		assert.Equal(t, "TicketMovedToInProgress", events[len(events)-1].EventName())
	})

	t.Run("should return an error when moving a non-open Ticket to in progress", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NoError(t, tk.MoveToInProgress())

		err = tk.MoveToInProgress()

		assert.ErrorIs(t, err, domainErr.ErrInvalidStatusTransition)
	})

	t.Run("should require a report of what was done to close", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.Close("  ")

		assert.ErrorIs(t, err, domainErr.ErrInvalidResolution)
		assert.Equal(t, valueobjects.TicketStatusOpen, tk.GetStatus())
	})

	t.Run("should keep the closing report, also after being rebuilt from history", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		assert.NoError(t, tk.Close("Troquei o cabo de rede"))

		assert.Equal(t, "Troquei o cabo de rede", tk.GetResolution())
		rebuilt, err := ticket.LoadFromHistory(tk.GetUncommittedEvents())
		assert.NoError(t, err)
		assert.Equal(t, "Troquei o cabo de rede", rebuilt.GetResolution())
	})

	t.Run("should remember when it was closed, also after being rebuilt from history", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.Nil(t, tk.GetClosedAt())

		assert.NoError(t, tk.Close("Resolvido"))
		closed := tk.GetUncommittedEvents()[len(tk.GetUncommittedEvents())-1]

		assert.Equal(t, closed.OccurredAt(), *tk.GetClosedAt())
		rebuilt, err := ticket.LoadFromHistory(tk.GetUncommittedEvents())
		assert.NoError(t, err)
		assert.Equal(t, closed.OccurredAt(), *rebuilt.GetClosedAt())
	})

	t.Run("should close an open Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.Close("Resolvido")

		assert.NoError(t, err)
		assert.Equal(t, valueobjects.TicketStatusClosed, tk.GetStatus())
		events := tk.GetUncommittedEvents()
		assert.Equal(t, "TicketClosed", events[len(events)-1].EventName())
	})

	t.Run("should close a Ticket that is in progress", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NoError(t, tk.MoveToInProgress())

		err = tk.Close("Resolvido")

		assert.NoError(t, err)
		assert.Equal(t, valueobjects.TicketStatusClosed, tk.GetStatus())
	})

	t.Run("should return an error when closing an already closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NoError(t, tk.Close("Resolvido"))

		err = tk.Close("Resolvido")

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should add a response to a Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		r := validResponseFor(t, tk.GetID())

		err = tk.AddResponse(r)

		assert.NoError(t, err)
		assert.Equal(t, []*response.Response{r}, tk.GetResponses())
		events := tk.GetUncommittedEvents()
		assert.Equal(t, "TicketResponseAdded", events[len(events)-1].EventName())
	})

	t.Run("should clear the uncommitted events", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NotEmpty(t, tk.GetUncommittedEvents())

		tk.ClearUncommittedEvents()

		assert.Empty(t, tk.GetUncommittedEvents())
	})

	t.Run("should return an error when adding a nil response", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		err = tk.AddResponse(nil)

		assert.ErrorIs(t, err, domainErr.ErrInvalidResponse)
	})

	t.Run("should return an error when the response belongs to a different Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)

		otherTicketID := valueobjects.NewID(uuid.New())
		r := validResponseFor(t, otherTicketID)

		err = tk.AddResponse(r)

		assert.ErrorIs(t, err, domainErr.ErrResponseTicketMismatch)
	})

	t.Run("should return an error when adding a response to a closed Ticket", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		tk, err := ticket.NewTicket(id, title, description, status, assignee, priority, "user-1")
		assert.NoError(t, err)
		assert.NoError(t, tk.Close("Resolvido"))

		r := validResponseFor(t, tk.GetID())

		err = tk.AddResponse(r)

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})

	t.Run("should reconstruct a Ticket from its event history", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		opened := event.NewTicketOpened(id, title, description, status, assignee, priority, "user-1")

		newTitle, err := valueobjects.NewTitle("New Title")
		assert.NoError(t, err)
		newDescription, err := valueobjects.NewDescription("New Description")
		assert.NoError(t, err)
		edited := event.NewTicketEdited(id, newTitle, newDescription)

		newAssignee, err := valueobjects.NewAssigneeID("agent-2")
		assert.NoError(t, err)
		assigned := event.NewTicketAssigned(id, &newAssignee)

		newPriority, err := valueobjects.NewPriority(string(valueobjects.TicketPriorityLow))
		assert.NoError(t, err)
		priorityChanged := event.NewTicketPriorityChanged(id, &newPriority)

		movedToInProgress := event.NewTicketMovedToInProgress(id)

		r := validResponseFor(t, id)
		responseAdded := event.NewTicketResponseAdded(id, r)

		history := []event.Event{opened, edited, assigned, priorityChanged, movedToInProgress, responseAdded}

		tk, err := ticket.LoadFromHistory(history)

		assert.NoError(t, err)
		assert.Equal(t, id.GetID(), tk.GetID().GetID())
		assert.Equal(t, "New Title", tk.GetTitle().GetTitle())
		assert.Equal(t, "New Description", tk.GetDescription().GetDescription())
		assert.Equal(t, valueobjects.TicketStatusInProgress, tk.GetStatus())
		assert.Equal(t, "agent-2", tk.GetAssigneeID().GetAssigneeID())
		assert.Equal(t, string(valueobjects.TicketPriorityLow), tk.GetPriority().GetPriority())
		assert.Len(t, tk.GetResponses(), 1)
		assert.Equal(t, r.GetContent().GetContent(), tk.GetResponses()[0].GetContent().GetContent())
		assert.Empty(t, tk.GetUncommittedEvents())
	})

	t.Run("should keep the original creation time of responses when reconstructing from history", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		opened := event.NewTicketOpened(id, title, description, status, assignee, priority, "user-1")
		r := validResponseFor(t, id)
		responseAdded := event.NewTicketResponseAdded(id, r)

		time.Sleep(time.Millisecond)
		tk, err := ticket.LoadFromHistory([]event.Event{opened, responseAdded})

		assert.NoError(t, err)
		assert.Len(t, tk.GetResponses(), 1)
		assert.Equal(t, r.GetCreatedAt(), tk.GetResponses()[0].GetCreatedAt())
	})

	t.Run("should reconstruct a closed Ticket from its event history", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		history := []event.Event{
			event.NewTicketOpened(id, title, description, status, assignee, priority, "user-1"),
			event.NewTicketClosed(id, "Resolvido"),
		}

		tk, err := ticket.LoadFromHistory(history)

		assert.NoError(t, err)
		assert.Equal(t, valueobjects.TicketStatusClosed, tk.GetStatus())
	})

	t.Run("should return an error when the event history is empty", func(t *testing.T) {
		tk, err := ticket.LoadFromHistory(nil)

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrEmptyEventHistory)
	})

	t.Run("should return an error when the first event is not TicketOpened", func(t *testing.T) {
		id, _, _, _, _, _ := validTicketParts(t)
		history := []event.Event{event.NewTicketClosed(id, "Resolvido")}

		tk, err := ticket.LoadFromHistory(history)

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidEventHistory)
	})

	t.Run("should return an error when an event in the history has an invalid payload", func(t *testing.T) {
		id, title, description, status, assignee, priority := validTicketParts(t)
		history := []event.Event{
			event.NewTicketOpened(id, title, description, status, assignee, priority, "user-1"),
			event.NewTicketAssigned(id, &valueobjects.AssigneeID{}),
		}

		tk, err := ticket.LoadFromHistory(history)

		assert.Nil(t, tk)
		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
	})
}
