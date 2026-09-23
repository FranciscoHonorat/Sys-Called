package messaging_test

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/adapters/in/messaging"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
)

type fakeResponsibleSyncer struct {
	upserted map[string]string
	inputs   []command.SyncResponsibleInput
	err      error
}

func newFakeResponsibleUpserter() *fakeResponsibleSyncer {
	return &fakeResponsibleSyncer{upserted: make(map[string]string)}
}

func (f *fakeResponsibleSyncer) Execute(_ context.Context, input command.SyncResponsibleInput) error {
	if f.err != nil {
		return f.err
	}
	f.upserted[input.ID] = input.Name
	f.inputs = append(f.inputs, input)
	return nil
}

type fakeAccountNotifier struct {
	inputs []command.AccountRequestInput
}

func (f *fakeAccountNotifier) Execute(_ context.Context, input command.AccountRequestInput) error {
	f.inputs = append(f.inputs, input)
	return nil
}

func employeeEvent(eventType, payload string) kafka.Message {
	return kafka.Message{
		Value:   []byte(payload),
		Headers: []kafka.Header{{Key: "event_type", Value: []byte(eventType)}},
	}
}

func employeeRegistered(payload string) kafka.Message {
	return kafka.Message{
		Value:   []byte(payload),
		Headers: []kafka.Header{{Key: "event_type", Value: []byte("EmployeeRegistered")}},
	}
}

func TestEmployeesConsumer(t *testing.T) {
	t.Run("should upsert a responsible on EmployeeRegistered", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		consumer := messaging.NewEmployeesConsumer(nil, upserter, &fakeAccountNotifier{})

		msg := employeeRegistered(`{"id":"agent-1","name":"Ana Souza"}`)

		err := consumer.Handle(context.Background(), msg)

		require.NoError(t, err)
		assert.Equal(t, "Ana Souza", upserter.upserted["agent-1"])
	})

	t.Run("should forward the employee role", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		consumer := messaging.NewEmployeesConsumer(nil, upserter, &fakeAccountNotifier{})

		msg := employeeRegistered(`{"id":"admin-1","name":"Admin","role":"admin"}`)

		require.NoError(t, consumer.Handle(context.Background(), msg))
		assert.Equal(t, []command.SyncResponsibleInput{{ID: "admin-1", Name: "Admin", Role: "admin"}}, upserter.inputs)
	})

	t.Run("should ignore unknown event types", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		consumer := messaging.NewEmployeesConsumer(nil, upserter, &fakeAccountNotifier{})

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
		consumer := messaging.NewEmployeesConsumer(nil, upserter, &fakeAccountNotifier{})

		msg := employeeRegistered(`not-json`)

		err := consumer.Handle(context.Background(), msg)

		assert.Error(t, err)
	})

	t.Run("should propagate sync errors", func(t *testing.T) {
		upserter := newFakeResponsibleUpserter()
		upserter.err = assert.AnError
		consumer := messaging.NewEmployeesConsumer(nil, upserter, &fakeAccountNotifier{})

		msg := employeeRegistered(`{"id":"agent-1","name":"Ana Souza"}`)

		err := consumer.Handle(context.Background(), msg)

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("should tell the admins about a new account", func(t *testing.T) {
		notifier := &fakeAccountNotifier{}
		consumer := messaging.NewEmployeesConsumer(nil, newFakeResponsibleUpserter(), notifier)

		require.NoError(t, consumer.Handle(context.Background(), employeeEvent("EmployeeSignedUp", `{"id":"user-9","name":"Maria Lima"}`)))

		assert.Equal(t, []command.AccountRequestInput{{Kind: command.AccountSignUp, EmployeeID: "user-9", Name: "Maria Lima"}}, notifier.inputs)
	})

	t.Run("should tell the admins about a password request", func(t *testing.T) {
		notifier := &fakeAccountNotifier{}
		consumer := messaging.NewEmployeesConsumer(nil, newFakeResponsibleUpserter(), notifier)

		require.NoError(t, consumer.Handle(context.Background(), employeeEvent("PasswordResetRequested", `{"id":"user-1","name":"Usuário Padrão"}`)))

		assert.Equal(t, []command.AccountRequestInput{{Kind: command.AccountPasswordReset, EmployeeID: "user-1", Name: "Usuário Padrão"}}, notifier.inputs)
	})
}
