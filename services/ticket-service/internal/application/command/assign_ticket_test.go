package command_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"

	"github.com/stretchr/testify/assert"
)

func openTestTicket(t *testing.T, store *outtest.EventStore) string {
	t.Helper()

	output, err := command.NewOpenTicketUseCase(store, newTestCache()).Execute(context.Background(), command.OpenTicketInput{
		Title:       "Valid Title",
		Description: "Valid Description",
	})
	assert.NoError(t, err)

	return output.TicketID
}

func TestAssignTicketUseCase(t *testing.T) {
	t.Run("should assign an existing ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewAssignTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.AssignTicketInput{
			TicketID:   ticketID,
			AssigneeID: "agent-1",
		})

		assert.NoError(t, err)

		tk := reloadTestTicket(t, store, ticketID)
		assert.Equal(t, "agent-1", tk.GetAssigneeID().GetAssigneeID())
	})

	t.Run("should return an error for an invalid ticket ID", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewAssignTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.AssignTicketInput{
			TicketID:   "not-a-uuid",
			AssigneeID: "agent-1",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidUUID)
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewAssignTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.AssignTicketInput{
			TicketID:   "00000000-0000-0000-0000-000000000001",
			AssigneeID: "agent-1",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error for an invalid assignee", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewAssignTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.AssignTicketInput{
			TicketID:   ticketID,
			AssigneeID: "",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidAssignee)
	})
}

func reloadTestTicket(t *testing.T, store *outtest.EventStore, ticketID string) *ticket.Ticket {
	t.Helper()

	id, err := uuid.Parse(ticketID)
	assert.NoError(t, err)

	history, err := store.Load(context.Background(), id)
	assert.NoError(t, err)

	tk, err := ticket.LoadFromHistory(history)
	assert.NoError(t, err)

	return tk
}
