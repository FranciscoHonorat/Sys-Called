package command_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"

	"github.com/stretchr/testify/assert"
)

func TestChangeTicketPriorityUseCase(t *testing.T) {
	t.Run("should change the priority of an existing ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewChangeTicketPriorityUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.ChangeTicketPriorityInput{
			TicketID: ticketID,
			Priority: "High",
		})

		assert.NoError(t, err)

		tk := reloadTestTicket(t, store, ticketID)
		assert.Equal(t, "High", tk.GetPriority().GetPriority())
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewChangeTicketPriorityUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.ChangeTicketPriorityInput{
			TicketID: "00000000-0000-0000-0000-000000000001",
			Priority: "High",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error for an invalid priority", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)
		uc := command.NewChangeTicketPriorityUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.ChangeTicketPriorityInput{
			TicketID: ticketID,
			Priority: "bogus",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidPriority)
	})
}
