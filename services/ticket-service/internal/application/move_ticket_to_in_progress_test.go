package application_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/valueobjects"

	"github.com/stretchr/testify/assert"
)

func TestMoveTicketToInProgressUseCase(t *testing.T) {
	t.Run("should move an existing open ticket to in progress", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		uc := application.NewMoveTicketToInProgressUseCase(store)

		err := uc.Execute(context.Background(), application.MoveTicketToInProgressInput{
			TicketID: ticketID,
		})

		assert.NoError(t, err)

		tk := reloadTestTicket(t, store, ticketID)
		assert.Equal(t, valueobjects.TicketStatusInProgress, tk.GetStatus())
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := newFakeEventStore()
		uc := application.NewMoveTicketToInProgressUseCase(store)

		err := uc.Execute(context.Background(), application.MoveTicketToInProgressInput{
			TicketID: "00000000-0000-0000-0000-000000000001",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error when the ticket is already in progress", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		uc := application.NewMoveTicketToInProgressUseCase(store)
		assert.NoError(t, uc.Execute(context.Background(), application.MoveTicketToInProgressInput{TicketID: ticketID}))

		err := uc.Execute(context.Background(), application.MoveTicketToInProgressInput{TicketID: ticketID})

		assert.ErrorIs(t, err, domainErr.ErrInvalidStatusTransition)
	})
}
