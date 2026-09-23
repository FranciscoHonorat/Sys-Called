package ticket

import (
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

type Permission func(t *Ticket, a actor.Actor) bool

func CanEdit(t *Ticket, a actor.Actor) bool {
	return a.IsAdmin() || (t.isRequestedBy(a.ID()) && t.isOpen())
}

func CanRespond(t *Ticket, a actor.Actor) bool {
	return t.IsVisibleTo(a)
}

func CanManage(_ *Ticket, a actor.Actor) bool {
	return a.IsAdmin() || a.IsSupport()
}

func CanWork(t *Ticket, a actor.Actor) bool {
	return a.IsSupport() && t.isAssignedTo(a.ID())
}
