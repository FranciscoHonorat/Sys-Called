package application_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/application/port/out/outtest"
	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

var agentCaller = out.Caller{ID: "agent-1", Role: employee.RoleSupport}

func repoWithPendingUser(t *testing.T) *outtest.EmployeeRepository {
	t.Helper()
	pending, err := employee.SignUp("user-9", "Maria Lima", "maria", "hashed:senha-forte")
	require.NoError(t, err)
	return &outtest.EmployeeRepository{Employees: []employee.Employee{pending}}
}

func TestApproveEmployeeUseCase(t *testing.T) {
	t.Run("lets an admin approve a pending account", func(t *testing.T) {
		repo := repoWithPendingUser(t)

		require.NoError(t, application.NewApproveEmployeeUseCase(repo).Execute(context.Background(), adminCaller, "user-9"))

		assert.Equal(t, employee.StatusActive, repo.Employees[0].GetStatus())
	})

	t.Run("is only for admins", func(t *testing.T) {
		err := application.NewApproveEmployeeUseCase(repoWithPendingUser(t)).Execute(context.Background(), agentCaller, "user-9")

		assert.ErrorIs(t, err, domainErr.ErrForbidden)
	})

	t.Run("tells when the account does not exist", func(t *testing.T) {
		err := application.NewApproveEmployeeUseCase(repoWithPendingUser(t)).Execute(context.Background(), adminCaller, "nobody")

		assert.ErrorIs(t, err, domainErr.ErrEmployeeNotFound)
	})
}

func TestIssueTemporaryPasswordUseCase(t *testing.T) {
	t.Run("gives the admin a password to hand over, which must be changed on the next login", func(t *testing.T) {
		repo := repoWithAdmin()
		uc := application.NewIssueTemporaryPasswordUseCase(repo, &outtest.PasswordHasher{}, outtest.PasswordGenerator{Password: "Temp-1234"}, newTestSessions(outtest.NewRefreshTokenStore()))

		password, err := uc.Execute(context.Background(), adminCaller, "admin-1")

		require.NoError(t, err)
		assert.Equal(t, "Temp-1234", password)
		assert.Equal(t, "hashed:Temp-1234", repo.Employees[0].GetPasswordHash())
		assert.True(t, repo.Employees[0].MustChangePassword())
	})

	t.Run("ends every open session of the employee", func(t *testing.T) {
		repo, sessions, store, _ := loggedInAdmin(t)
		uc := application.NewIssueTemporaryPasswordUseCase(repo, &outtest.PasswordHasher{}, outtest.PasswordGenerator{Password: "Temp-1234"}, sessions)

		_, err := uc.Execute(context.Background(), adminCaller, "admin-1")

		require.NoError(t, err)
		for hash, token := range store.Tokens {
			assert.True(t, token.IsRevoked(), hash)
		}
	})

	t.Run("is only for admins", func(t *testing.T) {
		uc := application.NewIssueTemporaryPasswordUseCase(repoWithAdmin(), &outtest.PasswordHasher{}, outtest.PasswordGenerator{Password: "x"}, newTestSessions(outtest.NewRefreshTokenStore()))

		_, err := uc.Execute(context.Background(), agentCaller, "admin-1")

		assert.ErrorIs(t, err, domainErr.ErrForbidden)
	})
}

func TestChangePasswordUseCase(t *testing.T) {
	t.Run("replaces the password after checking the current one", func(t *testing.T) {
		repo := repoWithAdmin()
		repo.Employees[0].SetTemporaryPassword("hashed:secret")

		err := application.NewChangePasswordUseCase(repo, &outtest.PasswordHasher{}).Execute(context.Background(), adminCaller,
			application.ChangePasswordInput{Current: "secret", New: "nova-senha-forte"})

		require.NoError(t, err)
		assert.Equal(t, "hashed:nova-senha-forte", repo.Employees[0].GetPasswordHash())
		assert.False(t, repo.Employees[0].MustChangePassword())
	})

	t.Run("refuses a wrong current password", func(t *testing.T) {
		err := application.NewChangePasswordUseCase(repoWithAdmin(), &outtest.PasswordHasher{}).Execute(context.Background(), adminCaller,
			application.ChangePasswordInput{Current: "errada", New: "nova-senha-forte"})

		assert.ErrorIs(t, err, domainErr.ErrInvalidCredentials)
	})

	t.Run("asks for a password of at least 8 characters", func(t *testing.T) {
		err := application.NewChangePasswordUseCase(repoWithAdmin(), &outtest.PasswordHasher{}).Execute(context.Background(), adminCaller,
			application.ChangePasswordInput{Current: "secret", New: "curta"})

		assert.ErrorIs(t, err, domainErr.ErrWeakPassword)
	})
}

func TestListEmployeesShowsAccountState(t *testing.T) {
	repo := repoWithPendingUser(t)

	output, err := application.NewListEmployeesUseCase(repo).Execute(context.Background(), adminCaller)

	require.NoError(t, err)
	assert.Equal(t, "pending", output[0].Status)
	assert.False(t, output[0].PasswordResetRequested)
}
