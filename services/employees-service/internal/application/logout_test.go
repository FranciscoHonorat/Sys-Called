package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
)

func TestLogoutUseCase(t *testing.T) {
	t.Run("should revoke the session's refresh token", func(t *testing.T) {
		_, sessions, store, refreshToken := loggedInAdmin(t)
		uc := application.NewLogoutUseCase(sessions)

		err := uc.Execute(context.Background(), refreshToken)

		require.NoError(t, err)
		assert.True(t, store.Tokens["sha:"+refreshToken].IsRevoked())
	})

	t.Run("should succeed for an unknown refresh token", func(t *testing.T) {
		_, sessions, _, _ := loggedInAdmin(t)
		uc := application.NewLogoutUseCase(sessions)

		assert.NoError(t, uc.Execute(context.Background(), "unknown"))
	})
}
