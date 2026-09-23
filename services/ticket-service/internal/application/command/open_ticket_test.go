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

func TestOpenTicketUseCase(t *testing.T) {
	t.Run("should record the actor as the one who opened the ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache())

		output, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       outtest.Actor("user-7", "user"),
			Title:       "Valid Title",
			Description: "Valid Description",
		})
		assert.NoError(t, err)

		history, err := store.Load(context.Background(), uuid.MustParse(output.TicketID))
		assert.NoError(t, err)
		tk, err := ticket.LoadFromHistory(history)
		assert.NoError(t, err)
		assert.Equal(t, "user-7", tk.GetRequesterID())
	})

	t.Run("should open a new ticket and persist a TicketOpened event", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache())

		output, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testAdmin,
			Title:       "Valid Title",
			Description: "Valid Description",
			AssigneeID:  "agent-1",
			Priority:    "High",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, output.TicketID)

		ticketID, err := uuid.Parse(output.TicketID)
		assert.NoError(t, err)

		history, err := store.Load(context.Background(), ticketID)
		assert.NoError(t, err)
		assert.Len(t, history, 1)
		assert.Equal(t, "TicketOpened", history[0].EventName())

		tk, err := ticket.LoadFromHistory(history)
		assert.NoError(t, err)
		assert.Equal(t, "Valid Title", tk.GetTitle().GetTitle())
		assert.Equal(t, "agent-1", tk.GetAssigneeID().GetAssigneeID())
	})

	t.Run("should forbid a regular user from choosing the assignee when opening a ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Valid Title",
			Description: "Valid Description",
			AssigneeID:  "agent-1",
		})

		assert.ErrorIs(t, err, domainErr.ErrForbidden)
	})

	t.Run("should open a new ticket without an assignee or priority", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache())

		output, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Valid Title",
			Description: "Valid Description",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, output.TicketID)
	})

	t.Run("should return an error for an invalid title", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "",
			Description: "Valid Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error for an invalid priority", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewOpenTicketUseCase(store, newTestCache())

		_, err := uc.Execute(context.Background(), command.OpenTicketInput{
			Actor:       testUser,
			Title:       "Valid Title",
			Description: "Valid Description",
			Priority:    "bogus",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidPriority)
	})
}
