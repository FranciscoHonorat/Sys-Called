package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func TestAutoAssignTicketUseCase(t *testing.T) {
	t.Run("should assign to the responsible with the fewest open tickets", func(t *testing.T) {
		store := newFakeEventStore()
		cache := newTestCache()

		first := openTestTicket(t, store)
		second := openTestTicket(t, store)
		third := openTestTicket(t, store)

		require.NoError(t, application.NewAssignTicketUseCase(store, cache).Execute(context.Background(), application.AssignTicketInput{
			TicketID: first, AssigneeID: "agent-1",
		}))
		require.NoError(t, application.NewAssignTicketUseCase(store, cache).Execute(context.Background(), application.AssignTicketInput{
			TicketID: second, AssigneeID: "agent-2",
		}))

		responsibles := &fakeResponsibleDirectory{responsibles: []string{"agent-1", "agent-2", "agent-3"}}
		uc := application.NewAutoAssignTicketUseCase(store, cache, responsibles)

		output, err := uc.Execute(context.Background(), application.AutoAssignTicketInput{TicketID: third})

		require.NoError(t, err)
		assert.Equal(t, "agent-3", output.AssigneeID)
	})

	t.Run("should not count closed tickets as open load", func(t *testing.T) {
		store := newFakeEventStore()
		cache := newTestCache()

		busy := openTestTicket(t, store)
		target := openTestTicket(t, store)

		require.NoError(t, application.NewAssignTicketUseCase(store, cache).Execute(context.Background(), application.AssignTicketInput{
			TicketID: busy, AssigneeID: "agent-1",
		}))
		require.NoError(t, application.NewCloseTicketUseCase(store, cache).Execute(context.Background(), application.CloseTicketInput{
			TicketID: busy,
		}))

		responsibles := &fakeResponsibleDirectory{responsibles: []string{"agent-1", "agent-2"}}
		uc := application.NewAutoAssignTicketUseCase(store, cache, responsibles)

		output, err := uc.Execute(context.Background(), application.AutoAssignTicketInput{TicketID: target})

		require.NoError(t, err)
		assert.Equal(t, "agent-1", output.AssigneeID)
	})

	t.Run("should return an error when there are no responsibles available", func(t *testing.T) {
		store := newFakeEventStore()
		cache := newTestCache()
		ticketID := openTestTicket(t, store)

		responsibles := &fakeResponsibleDirectory{responsibles: nil}
		uc := application.NewAutoAssignTicketUseCase(store, cache, responsibles)

		_, err := uc.Execute(context.Background(), application.AutoAssignTicketInput{TicketID: ticketID})

		assert.ErrorIs(t, err, domainErr.ErrNoResponsiblesAvailable)
	})

	t.Run("should return an error for an invalid ticket ID", func(t *testing.T) {
		store := newFakeEventStore()
		responsibles := &fakeResponsibleDirectory{responsibles: []string{"agent-1"}}
		uc := application.NewAutoAssignTicketUseCase(store, newTestCache(), responsibles)

		_, err := uc.Execute(context.Background(), application.AutoAssignTicketInput{TicketID: "not-a-uuid"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidUUID)
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := newFakeEventStore()
		responsibles := &fakeResponsibleDirectory{responsibles: []string{"agent-1"}}
		uc := application.NewAutoAssignTicketUseCase(store, newTestCache(), responsibles)

		_, err := uc.Execute(context.Background(), application.AutoAssignTicketInput{TicketID: "00000000-0000-0000-0000-000000000001"})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})
}
