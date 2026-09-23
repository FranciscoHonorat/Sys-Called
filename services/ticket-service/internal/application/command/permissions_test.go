package command_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

var (
	requester     = testUser
	stranger      = outtest.Actor("user-2", "user")
	assignedAgent = outtest.Actor("agent-1", "support")
	anotherAgent  = outtest.Actor("agent-2", "support")
	testAdmin     = outtest.Actor("admin-1", "admin")
)

type action func(store *outtest.EventStore, by actor.Actor, ticketID string) error

func ticketAssignedToAgent(t *testing.T) (*outtest.EventStore, string) {
	t.Helper()
	store := outtest.NewEventStore()
	ticketID := openTestTicket(t, store)
	require.NoError(t, command.NewAssignTicketUseCase(store, newTestCache()).Execute(context.Background(),
		command.AssignTicketInput{Actor: testAdmin, TicketID: ticketID, AssigneeID: "agent-1"}))
	return store, ticketID
}

func TestCommandPermissions(t *testing.T) {
	edit := func(store *outtest.EventStore, by actor.Actor, id string) error {
		return command.NewEditTicketUseCase(store, newTestCache()).Execute(context.Background(),
			command.EditTicketInput{Actor: by, TicketID: id, Title: "New", Description: "New"})
	}
	respond := func(store *outtest.EventStore, by actor.Actor, id string) error {
		return command.NewAddTicketResponseUseCase(store, newTestCache()).Execute(context.Background(),
			command.AddTicketResponseInput{Actor: by, TicketID: id, Content: "On it"})
	}
	assign := func(store *outtest.EventStore, by actor.Actor, id string) error {
		return command.NewAssignTicketUseCase(store, newTestCache()).Execute(context.Background(),
			command.AssignTicketInput{Actor: by, TicketID: id, AssigneeID: "agent-2"})
	}
	autoAssign := func(store *outtest.EventStore, by actor.Actor, id string) error {
		directory := &outtest.ResponsibleDirectory{Responsibles: []string{"agent-1", "agent-2"}}
		_, err := command.NewAutoAssignTicketUseCase(store, newTestCache(), directory).Execute(context.Background(),
			command.AutoAssignTicketInput{Actor: by, TicketID: id})
		return err
	}
	changePriority := func(store *outtest.EventStore, by actor.Actor, id string) error {
		return command.NewChangeTicketPriorityUseCase(store, newTestCache()).Execute(context.Background(),
			command.ChangeTicketPriorityInput{Actor: by, TicketID: id, Priority: "High"})
	}
	start := func(store *outtest.EventStore, by actor.Actor, id string) error {
		return command.NewMoveTicketToInProgressUseCase(store, newTestCache()).Execute(context.Background(),
			command.MoveTicketToInProgressInput{Actor: by, TicketID: id})
	}
	closeTicket := func(store *outtest.EventStore, by actor.Actor, id string) error {
		return command.NewCloseTicketUseCase(store, newTestCache()).Execute(context.Background(),
			command.CloseTicketInput{Resolution: "Resolvido", Actor: by, TicketID: id})
	}

	cases := []struct {
		name      string
		run       action
		allowed   actor.Actor
		forbidden []actor.Actor
	}{
		{"edit", edit, requester, []actor.Actor{anotherAgent}},
		{"respond", respond, assignedAgent, nil},
		{"assign", assign, anotherAgent, []actor.Actor{requester}},
		{"auto assign", autoAssign, anotherAgent, []actor.Actor{requester}},
		{"change priority", changePriority, anotherAgent, []actor.Actor{requester}},
		{"start", start, assignedAgent, []actor.Actor{anotherAgent, testAdmin}},
		{"close", closeTicket, assignedAgent, []actor.Actor{anotherAgent, testAdmin}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" is allowed for "+tc.allowed.ID(), func(t *testing.T) {
			store, ticketID := ticketAssignedToAgent(t)

			assert.NoError(t, tc.run(store, tc.allowed, ticketID))
		})
		for _, forbidden := range tc.forbidden {
			t.Run(tc.name+" is forbidden for "+forbidden.ID(), func(t *testing.T) {
				store, ticketID := ticketAssignedToAgent(t)

				assert.ErrorIs(t, tc.run(store, forbidden, ticketID), domainErr.ErrForbidden)
			})
		}
		t.Run(tc.name+" pretends the ticket does not exist for someone who cannot see it", func(t *testing.T) {
			store, ticketID := ticketAssignedToAgent(t)

			assert.ErrorIs(t, tc.run(store, stranger, ticketID), domainErr.ErrEventStreamNotFound)
		})
	}
}
