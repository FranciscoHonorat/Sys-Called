package query

import (
	"context"
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
)

const notificationsShown = 20

type NotificationOutput struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	TicketID  string    `json:"ticket_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Unread    bool      `json:"unread"`
}

type NotificationsOutput struct {
	Unread int                  `json:"unread"`
	Items  []NotificationOutput `json:"items"`
}

type ListNotificationsUseCase struct {
	store out.NotificationStore
}

func NewListNotificationsUseCase(store out.NotificationStore) *ListNotificationsUseCase {
	return &ListNotificationsUseCase{store: store}
}

func (uc *ListNotificationsUseCase) Execute(ctx context.Context, viewer actor.Actor) (NotificationsOutput, error) {
	notifications, err := uc.store.ListFor(ctx, viewer, notificationsShown)
	if err != nil {
		return NotificationsOutput{}, err
	}
	seenUntil, err := uc.store.SeenUntil(ctx, viewer.ID())
	if err != nil {
		return NotificationsOutput{}, err
	}

	output := NotificationsOutput{Items: make([]NotificationOutput, 0, len(notifications))}
	for _, n := range notifications {
		unread := n.CreatedAt.After(seenUntil)
		if unread {
			output.Unread++
		}
		item := NotificationOutput{ID: n.ID.String(), Message: n.Message, CreatedAt: n.CreatedAt, Unread: unread}
		if n.IsAboutTicket() {
			item.TicketID = n.TicketID.String()
		}
		output.Items = append(output.Items, item)
	}
	return output, nil
}
