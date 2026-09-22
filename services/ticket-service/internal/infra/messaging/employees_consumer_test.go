package messaging_test

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/infra/messaging"
)

type fakeResponsibleUpserter struct {
	upserted map[string]string
	err      error
}

func newFakeResponsibleUpserter() *fakeResponsibleUpserter {
	return &fakeResponsibleUpserter{upserted: make(map[string]string)}
}

func (f *fakeResponsibleUpserter) Upsert(_ context.Context, id, name string) error {
	if f.err != nil {
		return f.err
	}
	f.upserted[id] = name
	return nil
}

func TestEmployeesConsumer(t *testing.T) {
	t.Run("should upsert a responsible on EmployeeRegistered", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		consumer := messaging.NewEmployeesConsumer(nil, upserter)

		msg := kafka.Message{
			Value:   []byte(`{"id":"agent-1","name":"Ana Souza"}`),
			Headers: []kafka.Header{{Key: "event_type", Value: []byte("EmployeeRegistered")}},
		}

		err := consumer.Handle(context.Background(), msg)

		require.NoError(t, err)
		assert.Equal(t, "Ana Souza", upserter.upserted["agent-1"])
	})

	t.Run("should ignore unknown event types", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		consumer := messaging.NewEmployeesConsumer(nil, upserter)

		msg := kafka.Message{
			Value:   []byte(`{}`),
			Headers: []kafka.Header{{Key: "event_type", Value: []byte("SomethingElse")}},
		}

		err := consumer.Handle(context.Background(), msg)

		require.NoError(t, err)
		assert.Empty(t, upserter.upserted)
	})

	t.Run("should return an error for a malformed payload", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		consumer := messaging.NewEmployeesConsumer(nil, upserter)

		msg := kafka.Message{
			Value:   []byte(`not-json`),
			Headers: []kafka.Header{{Key: "event_type", Value: []byte("EmployeeRegistered")}},
		}

		err := consumer.Handle(context.Background(), msg)

		assert.Error(t, err)
	})

	t.Run("should propagate upserter errors", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		upserter.err = assert.AnError
		consumer := messaging.NewEmployeesConsumer(nil, upserter)

		msg := kafka.Message{
			Value:   []byte(`{"id":"agent-1","name":"Ana Souza"}`),
			Headers: []kafka.Header{{Key: "event_type", Value: []byte("EmployeeRegistered")}},
		}

		err := consumer.Handle(context.Background(), msg)

		assert.ErrorIs(t, err, assert.AnError)
	})
}
