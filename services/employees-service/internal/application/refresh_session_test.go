package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/session"
)

type staleReadStore struct {
	*outtest.RefreshTokenStore
}

func (s staleReadStore) FindByHash(ctx context.Context, hash string) (session.RefreshToken, bool, error) {
	token, found, err := s.RefreshTokenStore.FindByHash(ctx, hash)
	return session.NewRefreshToken(token.Hash(), token.EmployeeID(), token.ExpiresAt()), found, err
}

func loggedInAdmin(t *testing.T) (*outtest.EmployeeRepository, *application.Sessions, *outtest.RefreshTokenStore, string) {
	t.Helper()
	repo := repoWithAdmin()
	store := outtest.NewRefreshTokenStore()
	sessions := newTestSessions(store)

	login, err := application.NewLoginUseCase(repo, &outtest.PasswordHasher{}, sessions).
		Execute(context.Background(), application.LoginInput{Username: "admin", Password: "secret"})
	require.NoError(t, err)

	return repo, sessions, store, login.RefreshToken
}

func TestRefreshSessionUseCase(t *testing.T) {
	t.Run("should rotate a valid refresh token into a brand new session", func(t *testing.T) {
		repo, sessions, store, refreshToken := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, sessions)

		output, err := uc.Execute(context.Background(), refreshToken)

		require.NoError(t, err)
		assert.Equal(t, "token-for-admin-1", output.AccessToken)
		assert.Equal(t, "refresh-2", output.RefreshToken)
		assert.True(t, store.Tokens["sha:"+refreshToken].IsRevoked())
		assert.False(t, store.Tokens["sha:refresh-2"].IsRevoked())
	})

	t.Run("should reject an unknown refresh token", func(t *testing.T) {
		repo, sessions, store, _ := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, sessions)

		_, err := uc.Execute(context.Background(), "forged-token")

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Len(t, store.Tokens, 1)
	})

	t.Run("should revoke every session of the employee when a used refresh token is replayed", func(t *testing.T) {
		repo, sessions, store, stolenToken := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, sessions)
		_, err := uc.Execute(context.Background(), stolenToken)
		require.NoError(t, err)

		_, err = uc.Execute(context.Background(), stolenToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		for hash, token := range store.Tokens {
			assert.True(t, token.IsRevoked(), hash)
		}
	})

	t.Run("should let only one of two concurrent refreshes with the same token succeed", func(t *testing.T) {
		repo, _, store, refreshToken := loggedInAdmin(t)
		_, err := store.Revoke(context.Background(), "sha:"+refreshToken)
		require.NoError(t, err)
		sessions := application.NewSessions(outtest.TokenIssuer{}, &outtest.RefreshTokenGenerator{}, staleReadStore{store},
			testRefreshTTL, func() time.Time { return testNow })
		uc := application.NewRefreshSessionUseCase(repo, sessions)

		output, err := uc.Execute(context.Background(), refreshToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Empty(t, output.AccessToken)
		assert.Len(t, store.Tokens, 1)
	})

	t.Run("should reject the refresh token of an employee that no longer exists", func(t *testing.T) {
		repo, sessions, _, refreshToken := loggedInAdmin(t)
		repo.Employees = nil
		uc := application.NewRefreshSessionUseCase(repo, sessions)

		output, err := uc.Execute(context.Background(), refreshToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Empty(t, output.AccessToken)
	})

	t.Run("should reject an expired refresh token", func(t *testing.T) {
		repo, _, store, refreshToken := loggedInAdmin(t)
		uc := application.NewRefreshSessionUseCase(repo, newTestSessionsAt(store, testNow.Add(testRefreshTTL)))

		_, err := uc.Execute(context.Background(), refreshToken)

		assert.ErrorIs(t, err, domainErr.ErrInvalidRefreshToken)
		assert.Len(t, store.Tokens, 1)
	})
}
