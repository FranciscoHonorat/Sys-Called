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
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func repoWithAdmin() *outtest.EmployeeRepository {
	return &outtest.EmployeeRepository{Employees: []employee.Employee{
		employee.NewEmployee("admin-1", "Admin", "admin", employee.RoleAdmin, "hashed:secret"),
	}}
}

func TestLoginUseCase(t *testing.T) {
	t.Run("should issue a bearer token when the credentials match", func(t *testing.T) {
		repo := repoWithAdmin()
		uc := application.NewLoginUseCase(repo, &outtest.PasswordHasher{}, newTestSessions(outtest.NewRefreshTokenStore()))

		output, err := uc.Execute(context.Background(), application.LoginInput{Username: "admin", Password: "secret"})

		require.NoError(t, err)
		assert.Equal(t, "token-for-admin-1", output.AccessToken)
		assert.Equal(t, "Bearer", output.TokenType)
		assert.Equal(t, 900, output.ExpiresIn)
	})

	t.Run("should open a session with a refresh token stored only as a hash", func(t *testing.T) {
		store := outtest.NewRefreshTokenStore()
		uc := application.NewLoginUseCase(repoWithAdmin(), &outtest.PasswordHasher{}, newTestSessions(store))

		output, err := uc.Execute(context.Background(), application.LoginInput{Username: "admin", Password: "secret"})

		require.NoError(t, err)
		assert.Equal(t, "refresh-1", output.RefreshToken)
		assert.Equal(t, 7*24*3600, output.RefreshExpiresIn)
		stored, found := store.Tokens["sha:refresh-1"]
		require.True(t, found)
		assert.Equal(t, "admin-1", stored.EmployeeID())
		assert.False(t, stored.IsExpired(testNow.Add(7*24*time.Hour-time.Second)))
		assert.True(t, stored.IsExpired(testNow.Add(7*24*time.Hour)))
		assert.NotContains(t, store.Tokens, "refresh-1")
	})

	rejected := []struct {
		name  string
		input application.LoginInput
	}{
		{"an unknown username", application.LoginInput{Username: "nobody", Password: "secret"}},
		{"a wrong password", application.LoginInput{Username: "admin", Password: "wrong"}},
	}
	for _, tc := range rejected {
		t.Run("should reject "+tc.name+" with a generic error", func(t *testing.T) {
			repo := repoWithAdmin()
			uc := application.NewLoginUseCase(repo, &outtest.PasswordHasher{}, newTestSessions(outtest.NewRefreshTokenStore()))

			output, err := uc.Execute(context.Background(), tc.input)

			assert.ErrorIs(t, err, domainErr.ErrInvalidCredentials)
			assert.Empty(t, output.AccessToken)
		})
	}

	t.Run("should still compare a password for an unknown username so both failures take the same time", func(t *testing.T) {
		hasher := &outtest.PasswordHasher{}
		uc := application.NewLoginUseCase(repoWithAdmin(), hasher, newTestSessions(outtest.NewRefreshTokenStore()))

		_, err := uc.Execute(context.Background(), application.LoginInput{Username: "nobody", Password: "secret"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidCredentials)
		assert.Equal(t, 1, hasher.Comparisons)
	})
}
