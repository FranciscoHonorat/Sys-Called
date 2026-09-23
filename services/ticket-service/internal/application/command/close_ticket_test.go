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

func TestCloseTicketUseCase(t *testing.T) {
	t.Run("should close an existing ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openAssignedTicket(t, store)
		uc := command.NewCloseTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.CloseTicketInput{
			Resolution: "Resolvido",
			Actor:      assignedAgent,
			TicketID:   ticketID,
		})

		assert.NoError(t, err)

		tk := reloadTestTicket(t, store, ticketID)
		assert.Equal(t, valueobjects.TicketStatusClosed, tk.GetStatus())
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := outtest.NewEventStore()
		uc := command.NewCloseTicketUseCase(store, newTestCache())

		err := uc.Execute(context.Background(), command.CloseTicketInput{
			Resolution: "Resolvido",
			Actor:      assignedAgent,
			TicketID:   "00000000-0000-0000-0000-000000000001",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error when the ticket is already closed", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openAssignedTicket(t, store)
		uc := command.NewCloseTicketUseCase(store, newTestCache())
		assert.NoError(t, uc.Execute(context.Background(), command.CloseTicketInput{Resolution: "Resolvido", Actor: assignedAgent, TicketID: ticketID}))

		err := uc.Execute(context.Background(), command.CloseTicketInput{Resolution: "Resolvido", Actor: assignedAgent, TicketID: ticketID})

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})
}
