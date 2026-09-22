package messaging

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(writer *kafka.Writer) *KafkaPublisher {
	return &KafkaPublisher{writer: writer}
}

func (p *KafkaPublisher) Publish(ctx context.Context, eventType string, payload []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(eventType),
		Value: payload,
		Headers: []kafka.Header{
			{Key: "event_type", Value: []byte(eventType)},
		},
	})
}
