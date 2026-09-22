package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
)

func TestListTicketsUseCase(t *testing.T) {
	t.Run("should return an empty list when there are no tickets", func(t *testing.T) {
		store := newFakeEventStore()
		uc := application.NewListTicketsUseCase(store, newTestCache())

		output, err := uc.Execute(context.Background())

		assert.NoError(t, err)
		assert.Empty(t, output)
	})

	t.Run("should list tickets in the order they were opened", func(t *testing.T) {
		store := newFakeEventStore()
		openUC := application.NewOpenTicketUseCase(store, newTestCache())

		first, err := openUC.Execute(context.Background(), application.OpenTicketInput{
			Title:       "First Ticket",
			Description: "First Description",
		})
		require.NoError(t, err)

		second, err := openUC.Execute(context.Background(), application.OpenTicketInput{
			Title:       "Second Ticket",
			Description: "Second Description",
		})
		require.NoError(t, err)

		uc := application.NewListTicketsUseCase(store, newTestCache())
		output, err := uc.Execute(context.Background())

		require.NoError(t, err)
		require.Len(t, output, 2)
		assert.Equal(t, first.TicketID, output[0].TicketID)
		assert.Equal(t, "First Ticket", output[0].Title)
		assert.Equal(t, second.TicketID, output[1].TicketID)
		assert.Equal(t, "Second Ticket", output[1].Title)
	})

	t.Run("should reflect the current state of a mutated ticket", func(t *testing.T) {
		store := newFakeEventStore()
		opened, err := application.NewOpenTicketUseCase(store, newTestCache()).Execute(context.Background(), application.OpenTicketInput{
			Title:       "Valid Title",
			Description: "Valid Description",
		})
		require.NoError(t, err)

		err = application.NewCloseTicketUseCase(store, newTestCache()).Execute(context.Background(), application.CloseTicketInput{
			TicketID: opened.TicketID,
		})
		require.NoError(t, err)

		uc := application.NewListTicketsUseCase(store, newTestCache())
		output, err := uc.Execute(context.Background())

		require.NoError(t, err)
		require.Len(t, output, 1)
		assert.Equal(t, "Closed", output[0].Status)
	})
}
