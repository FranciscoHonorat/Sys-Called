package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type responsibleUpserter interface {
	Upsert(ctx context.Context, id, name string) error
}

type EmployeesConsumer struct {
	reader       *kafka.Reader
	responsibles responsibleUpserter
}

func NewEmployeesConsumer(reader *kafka.Reader, responsibles responsibleUpserter) *EmployeesConsumer {
	return &EmployeesConsumer{reader: reader, responsibles: responsibles}
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
		return c.responsibles.Upsert(ctx, payload.ID, payload.Name)
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
