package query_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/query"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"

	"github.com/stretchr/testify/assert"
)

func TestGetTicketUseCase(t *testing.T) {
	t.Run("should get an existing ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := query.NewGetTicketUseCase(store, newTestCache())

		output, err := uc.Execute(context.Background(), query.GetTicketInput{TicketID: ticketID})

		assert.NoError(t, err)
		assert.Equal(t, ticketID, output.TicketID)
		assert.Equal(t, "Valid Title", output.Title)
		assert.Equal(t, "Open", output.Status)
	})

	t.Run("should return an error for an invalid ticket ID", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := query.NewGetTicketUseCase(store, newTestCache())

		_, err := uc.Execute(context.Background(), query.GetTicketInput{TicketID: "not-a-uuid"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidUUID)
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := query.NewGetTicketUseCase(store, newTestCache())

		_, err := uc.Execute(context.Background(), query.GetTicketInput{TicketID: "00000000-0000-0000-0000-000000000001"})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should serve from cache without reading the store again", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		sharedCache := newTestCache()
		getUC := query.NewGetTicketUseCase(store, sharedCache)

		_, err := getUC.Execute(context.Background(), query.GetTicketInput{TicketID: ticketID})
		assert.NoError(t, err)

		store.ForgetAll()

		output, err := getUC.Execute(context.Background(), query.GetTicketInput{TicketID: ticketID})
		assert.NoError(t, err)
		assert.Equal(t, ticketID, output.TicketID)
	})

	t.Run("should populate the cache after a write use case", func(t *testing.T) {
		store := outtest.NewEventStore()
		sharedCache := newTestCache()
		openUC := command.NewOpenTicketUseCase(store, sharedCache)
		getUC := query.NewGetTicketUseCase(store, sharedCache)

		output, err := openUC.Execute(context.Background(), command.OpenTicketInput{
			Title:       "Valid Title",
			Description: "Valid Description",
		})
		assert.NoError(t, err)

		store.ForgetAll()

		got, err := getUC.Execute(context.Background(), query.GetTicketInput{TicketID: output.TicketID})
		assert.NoError(t, err)
		assert.Equal(t, output.TicketID, got.TicketID)
	})
}
