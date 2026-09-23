package command

import (
	"context"
	"time"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/notification"
)

type AccountRequest string

const (
	AccountSignUp        AccountRequest = "signup"
	AccountPasswordReset AccountRequest = "password_reset"
)

var accountNotifications = map[AccountRequest]func(employeeID, name string, at time.Time) notification.Notification{
	AccountSignUp:        notification.ForNewAccount,
	AccountPasswordReset: notification.ForPasswordResetRequest,
}

type AccountRequestInput struct {
	Kind       AccountRequest
	EmployeeID string
	Name       string
}

type NotifyAccountRequestUseCase struct {
	notifications out.NotificationStore
	now           func() time.Time
}

func NewNotifyAccountRequestUseCase(notifications out.NotificationStore, now func() time.Time) *NotifyAccountRequestUseCase {
	return &NotifyAccountRequestUseCase{notifications: notifications, now: now}
}

func (uc *NotifyAccountRequestUseCase) Execute(ctx context.Context, input AccountRequestInput) error {
	notify, known := accountNotifications[input.Kind]
	if !known {
		return nil
	}
	return uc.notifications.Add(ctx, []notification.Notification{notify(input.EmployeeID, input.Name, uc.now())})
}
