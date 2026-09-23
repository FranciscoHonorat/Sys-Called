package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
)

type responsibleSyncer interface {
	Execute(ctx context.Context, input command.SyncResponsibleInput) error
}

type accountNotifier interface {
	Execute(ctx context.Context, input command.AccountRequestInput) error
}

type EmployeesConsumer struct {
	reader          *kafka.Reader
	syncResponsible responsibleSyncer
	notifyAccount   accountNotifier
}

func NewEmployeesConsumer(reader *kafka.Reader, syncResponsible responsibleSyncer, notifyAccount accountNotifier) *EmployeesConsumer {
	return &EmployeesConsumer{reader: reader, syncResponsible: syncResponsible, notifyAccount: notifyAccount}
}

type employeePayload struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func (c *EmployeesConsumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("employees consumer: read error: %v", err)
			continue
		}

		if err := c.Handle(ctx, msg); err != nil {
			log.Printf("employees consumer: handle error: %v", err)
		}
	}
}

var accountRequests = map[string]command.AccountRequest{
	"EmployeeSignedUp":       command.AccountSignUp,
	"PasswordResetRequested": command.AccountPasswordReset,
}

func (c *EmployeesConsumer) Handle(ctx context.Context, msg kafka.Message) error {
	eventType := headerValue(msg.Headers, "event_type")
	request, isAccountRequest := accountRequests[eventType]
	if eventType != "EmployeeRegistered" && !isAccountRequest {
		return nil
	}

	var payload employeePayload
	if err := json.Unmarshal(msg.Value, &payload); err != nil {
		return err
	}
	if isAccountRequest {
		return c.notifyAccount.Execute(ctx, command.AccountRequestInput{Kind: request, EmployeeID: payload.ID, Name: payload.Name})
	}
	return c.syncResponsible.Execute(ctx, command.SyncResponsibleInput{ID: payload.ID, Name: payload.Name, Role: payload.Role})
}

func headerValue(headers []kafka.Header, key string) string {
	for _, h := range headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}
