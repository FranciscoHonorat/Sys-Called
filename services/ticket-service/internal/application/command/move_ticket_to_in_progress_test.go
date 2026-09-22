package command_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestMoveTicketToInProgressUseCase(t *testing.T) {
	t.Run("should move an existing open ticket to in progress", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewMoveTicketToInProgressUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.MoveTicketToInProgressInput{
			TicketID: ticketID,
		})

		assert.NoError(t, err)

		tk := reloadTestTicket(t, store, ticketID)
		assert.Equal(t, valueobjects.TicketStatusInProgress, tk.GetStatus())
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewMoveTicketToInProgressUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.MoveTicketToInProgressInput{
			TicketID: "00000000-0000-0000-0000-000000000001",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error when the ticket is already in progress", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewMoveTicketToInProgressUseCase(store, newTestCache())
		assert.NoError(t, uc.Execute(context.Background(), command.MoveTicketToInProgressInput{TicketID: ticketID}))

		err := uc.Execute(context.Background(), command.MoveTicketToInProgressInput{TicketID: ticketID})

		assert.ErrorIs(t, err, domainErr.ErrInvalidStatusTransition)
	})
}
