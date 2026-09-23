package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func TestAutoAssignTicketUseCase(t *testing.T) {
	t.Run("should assign to the responsible with the fewest open tickets", func(t *testing.T) {
		store := outtest.NewEventStore()
		cache := newTestCache()

		first := openTestTicket(t, store)
		second := openTestTicket(t, store)
		third := openTestTicket(t, store)

		require.NoError(t, command.NewAssignTicketUseCase(store, cache).Execute(context.Background(), command.AssignTicketInput{
			Actor:    testAdmin,
			TicketID: first, AssigneeID: "agent-1",
		}))
		require.NoError(t, command.NewAssignTicketUseCase(store, cache).Execute(context.Background(), command.AssignTicketInput{
			Actor:    testAdmin,
			TicketID: second, AssigneeID: "agent-2",
		}))

		responsibles := &outtest.ResponsibleDirectory{Responsibles: []string{"agent-1", "agent-2", "agent-3"}}
		uc := command.NewAutoAssignTicketUseCase(store, cache, responsibles)

		output, err := uc.Execute(context.Background(), command.AutoAssignTicketInput{Actor: testAdmin, TicketID: third})

		require.NoError(t, err)
		assert.Equal(t, "agent-3", output.AssigneeID)
	})

	t.Run("should not count closed tickets as open load", func(t *testing.T) {
		store := outtest.NewEventStore()
		cache := newTestCache()

		busy := openTestTicket(t, store)
		target := openTestTicket(t, store)

		require.NoError(t, command.NewAssignTicketUseCase(store, cache).Execute(context.Background(), command.AssignTicketInput{
			Actor:    testAdmin,
			TicketID: busy, AssigneeID: "agent-1",
		}))
		require.NoError(t, command.NewCloseTicketUseCase(store, cache).Execute(context.Background(), command.CloseTicketInput{
			Resolution: "Resolvido",
			Actor:      assignedAgent,
			TicketID:   busy,
		}))

		responsibles := &outtest.ResponsibleDirectory{Responsibles: []string{"agent-1", "agent-2"}}
		uc := command.NewAutoAssignTicketUseCase(store, cache, responsibles)

		output, err := uc.Execute(context.Background(), command.AutoAssignTicketInput{Actor: testAdmin, TicketID: target})

		require.NoError(t, err)
		assert.Equal(t, "agent-1", output.AssigneeID)
	})

	t.Run("should return an error when there are no responsibles available", func(t *testing.T) {
		store := outtest.NewEventStore()
		cache := newTestCache()
		ticketID := openTestTicket(t, store)

		responsibles := &outtest.ResponsibleDirectory{Responsibles: nil}
		uc := command.NewAutoAssignTicketUseCase(store, cache, responsibles)

		_, err := uc.Execute(context.Background(), command.AutoAssignTicketInput{Actor: testAdmin, TicketID: ticketID})

		assert.ErrorIs(t, err, domainErr.ErrNoResponsiblesAvailable)
	})

	t.Run("should return an error for an invalid ticket ID", func(t *testing.T) {
		store := outtest.NewEventStore()
		responsibles := &outtest.ResponsibleDirectory{Responsibles: []string{"agent-1"}}
		uc := command.NewAutoAssignTicketUseCase(store, newTestCache(), responsibles)

		_, err := uc.Execute(context.Background(), command.AutoAssignTicketInput{Actor: testAdmin, TicketID: "not-a-uuid"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidUUID)
	})

	t.Run("should return an error when the ticket does not exist", func(t *testing.T) {
		store := outtest.NewEventStore()
		responsibles := &outtest.ResponsibleDirectory{Responsibles: []string{"agent-1"}}
		uc := command.NewAutoAssignTicketUseCase(store, newTestCache(), responsibles)

		_, err := uc.Execute(context.Background(), command.AutoAssignTicketInput{Actor: testAdmin, TicketID: "00000000-0000-0000-0000-000000000001"})

		assert.ErrorIs(t, err, domainErr.ErrEventStreamNotFound)
	})
}
