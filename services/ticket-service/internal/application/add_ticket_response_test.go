package application_test

import (
	"context"
	"testing"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"

	"github.com/stretchr/testify/assert"
)

func TestAddTicketResponseUseCase(t *testing.T) {
	t.Run("should add a response to an existing ticket", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		uc := application.NewAddTicketResponseUseCase(store)

		err := uc.Execute(context.Background(), application.AddTicketResponseInput{
			TicketID: ticketID,
			AuthorID: "agent-1",
			Content:  "Valid response content",
		})

		assert.NoError(t, err)

		tk := reloadTestTicket(t, store, ticketID)
		assert.Len(t, tk.GetResponses(), 1)
		assert.Equal(t, "Valid response content", tk.GetResponses()[0].GetContent().GetContent())
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := newFakeEventStore()
		uc := application.NewAddTicketResponseUseCase(store)

		err := uc.Execute(context.Background(), application.AddTicketResponseInput{
			TicketID: "00000000-0000-0000-0000-000000000001",
			AuthorID: "agent-1",
			Content:  "Valid response content",
		})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})

	t.Run("should return an error for an empty content", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		uc := application.NewAddTicketResponseUseCase(store)

		err := uc.Execute(context.Background(), application.AddTicketResponseInput{
			TicketID: ticketID,
			AuthorID: "agent-1",
			Content:  "",
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidContent)
	})

	t.Run("should return an error when adding a response to a closed ticket", func(t *testing.T) {
		store := newFakeEventStore()
		ticketID := openTestTicket(t, store)
		assert.NoError(t, application.NewCloseTicketUseCase(store).Execute(context.Background(), application.CloseTicketInput{TicketID: ticketID}))
		uc := application.NewAddTicketResponseUseCase(store)

		err := uc.Execute(context.Background(), application.AddTicketResponseInput{
			TicketID: ticketID,
			AuthorID: "agent-1",
			Content:  "Valid response content",
		})

		assert.ErrorIs(t, err, domainErr.ErrTicketAlreadyClosed)
	})
}
