package employee_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainErr "github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/domain-errors"
	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/domain/employee"
)

func TestAccountLifecycle(t *testing.T) {
	t.Run("a sign up creates a regular user waiting for approval", func(t *testing.T) {
		e, err := employee.SignUp("user-9", "Maria Lima", "maria", "hash")

		require.NoError(t, err)
		assert.Equal(t, employee.RoleUser, e.GetRole())
		assert.Equal(t, employee.StatusPending, e.GetStatus())
		assert.False(t, e.CanLogIn())
		assert.Equal(t, []employee.Event{employee.SignedUp{ID: "user-9", Name: "Maria Lima"}}, e.Events())
	})

	t.Run("a sign up needs name and username", func(t *testing.T) {
		_, err := employee.SignUp("user-9", "", "maria", "hash")
		assert.ErrorIs(t, err, domainErr.ErrInvalidEmployeeName)

		_, err = employee.SignUp("user-9", "Maria Lima", "", "hash")
		assert.ErrorIs(t, err, domainErr.ErrInvalidUsername)
	})

	t.Run("an approved user can log in", func(t *testing.T) {
		e, err := employee.SignUp("user-9", "Maria Lima", "maria", "hash")
		require.NoError(t, err)

		e.Approve()

		assert.Equal(t, employee.StatusActive, e.GetStatus())
		assert.True(t, e.CanLogIn())
	})

	t.Run("asking for a new password is recorded for the admin", func(t *testing.T) {
		e := employee.NewEmployee("user-1", "Usuário Padrão", "usuario", employee.RoleUser, "hash")

		e.RequestPasswordReset()

		assert.True(t, e.HasRequestedPasswordReset())
		assert.Equal(t, []employee.Event{employee.PasswordResetRequested{ID: "user-1", Name: "Usuário Padrão"}}, e.Events())
	})

	t.Run("asking again for a new password does not announce it twice", func(t *testing.T) {
		e := employee.Rehydrate(employee.Snapshot{ID: "user-1", Name: "Usuário Padrão", PasswordResetRequested: true})

		e.RequestPasswordReset()

		assert.True(t, e.HasRequestedPasswordReset())
		assert.Empty(t, e.Events())
	})

	t.Run("a temporary password must be changed on the next login", func(t *testing.T) {
		e := employee.NewEmployee("user-1", "Usuário Padrão", "usuario", employee.RoleUser, "hash")
		e.RequestPasswordReset()

		e.SetTemporaryPassword("temporary-hash")

		assert.Equal(t, "temporary-hash", e.GetPasswordHash())
		assert.True(t, e.MustChangePassword())
		assert.False(t, e.HasRequestedPasswordReset())
	})

	t.Run("choosing a new password clears the obligation", func(t *testing.T) {
		e := employee.NewEmployee("user-1", "Usuário Padrão", "usuario", employee.RoleUser, "hash")
		e.SetTemporaryPassword("temporary-hash")

		e.ChangePassword("new-hash")

		assert.Equal(t, "new-hash", e.GetPasswordHash())
		assert.False(t, e.MustChangePassword())
	})

	t.Run("an employee rebuilt from storage keeps its account state without events", func(t *testing.T) {
		e := employee.Rehydrate(employee.Snapshot{
			ID: "user-9", Name: "Maria Lima", Username: "maria", Role: employee.RoleUser, PasswordHash: "hash",
			Status: employee.StatusPending, MustChangePassword: true, PasswordResetRequested: true,
		})

		assert.Equal(t, employee.StatusPending, e.GetStatus())
		assert.True(t, e.MustChangePassword())
		assert.True(t, e.HasRequestedPasswordReset())
		assert.Empty(t, e.Events())
	})

	t.Run("employees created by the system are active", func(t *testing.T) {
		e, err := employee.Register("agent-1", "Ana Souza", "ana", employee.RoleSupport, "hash")
		require.NoError(t, err)

		assert.True(t, e.CanLogIn())
	})
}
