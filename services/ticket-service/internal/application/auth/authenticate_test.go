package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/auth"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/application/port/out/outtest"
	"github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/actor"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/ticket-service/internal/domain/domain-errors"
)

func TestAuthenticateUseCase(t *testing.T) {
	admin, err := actor.New("admin-1", "Administradora", "admin")
	require.NoError(t, err)

	t.Run("should return the actor of a valid token", func(t *testing.T) {
		verifier := &outtest.TokenVerifier{Tokens: map[string]actor.Actor{"valid": admin}}
		uc := auth.NewAuthenticateUseCase(verifier)

		got, err := uc.Execute(context.Background(), "valid")

		require.NoError(t, err)
		assert.Equal(t, admin, got)
	})

	t.Run("should not even try to verify an empty token", func(t *testing.T) {
		verifier := &outtest.TokenVerifier{}
		uc := auth.NewAuthenticateUseCase(verifier)

		_, err := uc.Execute(context.Background(), "")

		assert.ErrorIs(t, err, domainErr.ErrUnauthenticated)
		assert.Zero(t, verifier.Calls)
	})

	t.Run("should hide why a token was rejected", func(t *testing.T) {
		verifier := &outtest.TokenVerifier{Err: errors.New("signature is invalid")}
		uc := auth.NewAuthenticateUseCase(verifier)

		_, err := uc.Execute(context.Background(), "tampered")

		assert.ErrorIs(t, err, domainErr.ErrUnauthenticated)
		assert.NotContains(t, err.Error(), "signature")
	})
}
