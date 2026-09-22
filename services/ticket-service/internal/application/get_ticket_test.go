package application_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"

	"github.com/stretchr/testify/assert"
)

func TestGetTicketUseCase(t *testing.T) {
	t.Run("should get an existing ticket", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		uc := application.NewGetTicketUseCase(store, newTestCache())

		output, err := uc.Execute(context.Background(), application.GetTicketInput{TicketID: ticketID})

		assert.NoError(t, err)
		assert.Equal(t, ticketID, output.TicketID)
		assert.Equal(t, "Valid Title", output.Title)
		assert.Equal(t, "Open", output.Status)
	})

	t.Run("should return an error for an invalid ticket ID", func(t *testing.T) {
		store := newFakeEventStore()
		uc := application.NewGetTicketUseCase(store, newTestCache())

		_, err := uc.Execute(context.Background(), application.GetTicketInput{TicketID: "not-a-uuid"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidUUID)
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := newFakeEventStore()
		uc := application.NewGetTicketUseCase(store, newTestCache())

		_, err := uc.Execute(context.Background(), application.GetTicketInput{TicketID: "00000000-0000-0000-0000-000000000001"})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should serve from cache without reading the store again", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		sharedCache := newTestCache()
		getUC := application.NewGetTicketUseCase(store, sharedCache)

		_, err := getUC.Execute(context.Background(), application.GetTicketInput{TicketID: ticketID})
		assert.NoError(t, err)

		store.forgetAll()

		output, err := getUC.Execute(context.Background(), application.GetTicketInput{TicketID: ticketID})
		assert.NoError(t, err)
		assert.Equal(t, ticketID, output.TicketID)
	})

	t.Run("should populate the cache after a write use case", func(t *testing.T) {
		store := newFakeEventStore()
		sharedCache := newTestCache()
		openUC := application.NewOpenTicketUseCase(store, sharedCache)
		getUC := application.NewGetTicketUseCase(store, sharedCache)

		output, err := openUC.Execute(context.Background(), application.OpenTicketInput{
			Title:       "Valid Title",
			Description: "Valid Description",
		})
		assert.NoError(t, err)

		store.forgetAll()

		got, err := getUC.Execute(context.Background(), application.GetTicketInput{TicketID: output.TicketID})
		assert.NoError(t, err)
		assert.Equal(t, output.TicketID, got.TicketID)
	})
}
