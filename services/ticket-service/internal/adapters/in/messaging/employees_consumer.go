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

type EmployeesConsumer struct {
	reader          *kafka.Reader
	syncResponsible responsibleSyncer
}

func NewEmployeesConsumer(reader *kafka.Reader, syncResponsible responsibleSyncer) *EmployeesConsumer {
	return &EmployeesConsumer{reader: reader, syncResponsible: syncResponsible}
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

func (c *EmployeesConsumer) Handle(ctx context.Context, msg kafka.Message) error {
	switch headerValue(msg.Headers, "event_type") {
	case "EmployeeRegistered":
		var payload struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			return err
		}
		return c.syncResponsible.Execute(ctx, command.SyncResponsibleInput{ID: payload.ID, Name: payload.Name})
	default:
		return nil
	}
}

func headerValue(headers []kafka.Header, key string) string {
	for _, h := range headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}
