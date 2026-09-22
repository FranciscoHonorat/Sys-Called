package application_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"

	"github.com/stretchr/testify/assert"
)

func TestEditTicketUseCase(t *testing.T) {
	t.Run("should edit the title and description of an existing ticket", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		uc := application.NewEditTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), application.EditTicketInput{
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
		store := newFakeEventStore()
		uc := application.NewEditTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), application.EditTicketInput{
			TicketID:    "00000000-0000-0000-0000-000000000001",
			Title:       "New Title",
			Description: "New Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error for an invalid title", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		uc := application.NewEditTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), application.EditTicketInput{
			TicketID:    ticketID,
			Title:       "",
			Description: "New Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidTitle)
	})

	t.Run("should return an error when editing a closed ticket", func(t *testing.T) {
		store := newFakeEventStore()
		cache := newTestCache()
		ticketID := openTestTicket(t, store)
		assert.NoError(t, application.NewCloseTicketUseCase(store, cache).Execute(context.Background(), application.CloseTicketInput{
			TicketID: ticketID,
		}))

		uc := application.NewEditTicketUseCase(store, cache)
		err := uc.Execute(context.Background(), application.EditTicketInput{
			TicketID:    ticketID,
			Title:       "New Title",
			Description: "New Description",
		})

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})
}
