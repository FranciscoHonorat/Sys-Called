package command_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"

	"github.com/stretchr/testify/assert"
)

func TestEditTicketUseCase(t *testing.T) {
	t.Run("should edit the title and description of an existing ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewEditTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.EditTicketInput{
			TicketID:    ticketID,
			Title:       "New Title",
			Description: "New Description",
		})

		assert.NoError(t, err)

		tk := reloadTestTicket(t, store, ticketID)
		assert.Equal(t, "New Title", tk.GetTitle().GetTitle())
		assert.Equal(t, "New Description", tk.GetDescription().GetDescription())
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewEditTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.EditTicketInput{
			TicketID:    "00000000-0000-0000-0000-000000000001",
			Title:       "New Title",
			Description: "New Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error for an invalid title", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewEditTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.EditTicketInput{
			TicketID:    ticketID,
			Title:       "",
			Description: "New Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error when editing a closed ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		cache := newTestCache()
		ticketID := openTestTicket(t, store)
		assert.NoError(t, command.NewCloseTicketUseCase(store, cache).Execute(context.Background(), command.CloseTicketInput{
			TicketID: ticketID,
		}))

		uc := command.NewEditTicketUseCase(store, cache)
		err := uc.Execute(context.Background(), command.EditTicketInput{
			TicketID:    ticketID,
			Title:       "New Title",
			Description: "New Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})
}
