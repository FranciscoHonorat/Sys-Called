package notification

import (
	"time"

	"github.com/google/uuid"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

func ForNewAccount(employeeID, name string, at time.Time) Notification {
	return toAdmins("Nova conta aguardando aprovação: "+name, employeeID, at)
}

func ForPasswordResetRequest(employeeID, name string, at time.Time) Notification {
	return toAdmins("Pedido de nova senha: "+name, employeeID, at)
}

func toAdmins(message, actorID string, at time.Time) Notification {
	return Notification{
		ID:        uuid.New(),
		Audience:  ToRole(actor.RoleAdmin),
		Message:   message,
		TicketID:  uuid.Nil,
		ActorID:   actorID,
		CreatedAt: at,
	}
}

func (n Notification) IsAboutTicket() bool {
	return n.TicketID != uuid.Nil
}
