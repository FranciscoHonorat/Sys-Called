package command_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func TestAuditTrail(t *testing.T) {
	t.Run("should record who performed each action on the ticket", func(t *testing.T) {
		store := outtest.NewEventStore()
		cache := newTestCache()
		admin, agent := outtest.Actor("admin-1", "admin"), outtest.Actor("agent-1", "support")

		opened, err := command.NewOpenTicketUseCase(store, cache).Execute(context.Background(), command.OpenTicketInput{
			Actor: testUser, Title: "Printer", Description: "Broken",
		})
		require.NoError(t, err)
		require.NoError(t, command.NewAssignTicketUseCase(store, cache).Execute(context.Background(), command.AssignTicketInput{
			Actor: admin, TicketID: opened.TicketID, AssigneeID: "agent-1",
		}))
		require.NoError(t, command.NewCloseTicketUseCase(store, cache).Execute(context.Background(), command.CloseTicketInput{
			Resolution: "Resolvido",
			Actor:      agent, TicketID: opened.TicketID,
		}))

		assert.Equal(t, []string{"user-1", "admin-1", "agent-1"}, store.ActorsOf(uuid.MustParse(opened.TicketID)))
	})

	t.Run("should refuse a change without knowing who is making it", func(t *testing.T) {
		store := outtest.NewEventStore()
		ticketID := openTestTicket(t, store)

		err := command.NewCloseTicketUseCase(store, newTestCache()).Execute(context.Background(), command.CloseTicketInput{
			Resolution: "Resolvido",
			Actor:      actor.Actor{}, TicketID: ticketID,
		})

		assert.ErrorIs(t, err, domainErr.ErrInvalidActor)
	})
}
