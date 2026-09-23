package notification

import (
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/event"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/ticket"
)

type AudienceKind string

const (
	AudienceUser AudienceKind = "user"
	AudienceRole AudienceKind = "role"
)

type Audience struct {
	Kind  AudienceKind
	Value string
}

func ToUser(id string) Audience {
	return Audience{Kind: AudienceUser, Value: id}
}

func ToRole(role actor.Role) Audience {
	return Audience{Kind: AudienceRole, Value: string(role)}
}

type Notification struct {
	ID        uuid.UUID
	Audience  Audience
	Message   string
	TicketID  uuid.UUID
	ActorID   string
	CreatedAt time.Time
}

func ForTicketEvent(t *ticket.Ticket, e event.Event, actorID string) []Notification {
	title := t.GetTitle().GetTitle()
	assignee := ""
	if t.GetAssigneeID() != nil {
		assignee = t.GetAssigneeID().GetAssigneeID()
	}

	var audiences []Audience
	var message string
	switch e.(type) {
	case event.TicketOpened:
		audiences, message = []Audience{ToRole(actor.RoleSupport), ToRole(actor.RoleAdmin)}, "Novo chamado: "+title
	case event.TicketAssigned:
		audiences, message = users(actorID, assignee), "Chamado atribuído a você: "+title
	case event.TicketResponseAdded:
		audiences, message = users(actorID, t.GetRequesterID(), assignee), "Nova mensagem em: "+title
	case event.TicketMovedToInProgress:
		audiences, message = users(actorID, t.GetRequesterID()), "Seu chamado está em atendimento: "+title
	case event.TicketClosed:
		audiences, message = users(actorID, t.GetRequesterID()), "Seu chamado foi fechado: "+title
	}

	notifications := make([]Notification, 0, len(audiences))
	for _, audience := range audiences {
		notifications = append(notifications, Notification{
			ID:        uuid.New(),
			Audience:  audience,
			Message:   message,
			TicketID:  e.AggregateID(),
			ActorID:   actorID,
			CreatedAt: e.OccurredAt(),
		})
	}
	return notifications
}

func users(actorID string, ids ...string) []Audience {
	var audiences []Audience
	seen := map[string]bool{actorID: true, "": true}
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			audiences = append(audiences, ToUser(id))
		}
	}
	return audiences
}
