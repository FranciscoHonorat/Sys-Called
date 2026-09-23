package application_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestSignUpUseCase(t *testing.T) {
	t.Run("creates an account waiting for approval with a hashed password", func(t *testing.T) {
		repo := &outtest.EmployeeRepository{}

		err := application.NewSignUpUseCase(repo, &outtest.PasswordHasher{}).Execute(context.Background(),
			application.SignUpInput{Name: "Maria Lima", Username: "maria", Password: "senha-forte"})

		require.NoError(t, err)
		require.Len(t, repo.Employees, 1)
		created := repo.Employees[0]
		assert.NotEmpty(t, created.GetID())
		assert.Equal(t, employee.StatusPending, created.GetStatus())
		assert.Equal(t, employee.RoleUser, created.GetRole())
		assert.Equal(t, "hashed:senha-forte", created.GetPasswordHash())
	})

	t.Run("asks for a password of at least 8 characters", func(t *testing.T) {
		err := application.NewSignUpUseCase(&outtest.EmployeeRepository{}, &outtest.PasswordHasher{}).Execute(context.Background(),
			application.SignUpInput{Name: "Maria Lima", Username: "maria", Password: "curta"})

		assert.ErrorIs(t, err, domainErr.ErrWeakPassword)
	})

	t.Run("refuses a password longer than 72 bytes", func(t *testing.T) {
		err := application.NewSignUpUseCase(&outtest.EmployeeRepository{}, &outtest.PasswordHasher{}).Execute(context.Background(),
			application.SignUpInput{Name: "Maria Lima", Username: "maria", Password: strings.Repeat("a", 73)})

		assert.ErrorIs(t, err, domainErr.ErrPasswordTooLong)
	})

	t.Run("refuses a username that is already taken", func(t *testing.T) {
		repo := repoWithAdmin()

		err := application.NewSignUpUseCase(repo, &outtest.PasswordHasher{}).Execute(context.Background(),
			application.SignUpInput{Name: "Outra", Username: "admin", Password: "senha-forte"})

		assert.ErrorIs(t, err, domainErr.ErrUsernameTaken)
	})
}

func TestLoginOfPendingAccounts(t *testing.T) {
	pending, err := employee.SignUp("user-9", "Maria Lima", "maria", "hashed:senha-forte")
	require.NoError(t, err)
	repo := &outtest.EmployeeRepository{Employees: []employee.Employee{pending}}
	uc := application.NewLoginUseCase(repo, &outtest.PasswordHasher{}, newTestSessions(outtest.NewRefreshTokenStore()))

	t.Run("tells the right password owner that the account waits for approval", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), application.LoginInput{Username: "maria", Password: "senha-forte"})

		assert.ErrorIs(t, err, domainErr.ErrPendingApproval)
	})

	t.Run("reveals nothing to someone without the password", func(t *testing.T) {
		_, err := uc.Execute(context.Background(), application.LoginInput{Username: "maria", Password: "chute"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidCredentials)
	})
}

func TestRequestPasswordResetUseCase(t *testing.T) {
	t.Run("records the request of an existing user", func(t *testing.T) {
		repo := repoWithAdmin()

		require.NoError(t, application.NewRequestPasswordResetUseCase(repo).Execute(context.Background(), "admin"))

		assert.True(t, repo.Employees[0].HasRequestedPasswordReset())
	})

	t.Run("answers the same for an unknown username", func(t *testing.T) {
		repo := repoWithAdmin()

		assert.NoError(t, application.NewRequestPasswordResetUseCase(repo).Execute(context.Background(), "nobody"))
		assert.False(t, repo.Employees[0].HasRequestedPasswordReset())
	})
}
