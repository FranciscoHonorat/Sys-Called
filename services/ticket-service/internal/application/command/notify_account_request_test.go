package command_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/command"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
)

func TestNotifyAccountRequestUseCase(t *testing.T) {
	at := time.Date(2026, 9, 23, 10, 0, 0, 0, time.UTC)
	now := func() time.Time { return at }

	t.Run("tells the admins about a new account", func(t *testing.T) {
		store := &outtest.NotificationStore{}

		err := command.NewNotifyAccountRequestUseCase(store, now).Execute(context.Background(), command.AccountRequestInput{
			Kind: command.AccountSignUp, EmployeeID: "user-9", Name: "Maria Lima",
		})

		require.NoError(t, err)
		require.Len(t, store.Added, 1)
		assert.Equal(t, "Nova conta aguardando aprovação: Maria Lima", store.Added[0].Message)
		assert.Equal(t, at, store.Added[0].CreatedAt)
	})

	t.Run("tells the admins about a password request", func(t *testing.T) {
		store := &outtest.NotificationStore{}

		err := command.NewNotifyAccountRequestUseCase(store, now).Execute(context.Background(), command.AccountRequestInput{
			Kind: command.AccountPasswordReset, EmployeeID: "user-1", Name: "Usuário Padrão",
		})

		require.NoError(t, err)
		require.Len(t, store.Added, 1)
		assert.Equal(t, "Pedido de nova senha: Usuário Padrão", store.Added[0].Message)
	})
}
