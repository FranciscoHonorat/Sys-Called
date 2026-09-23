package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/security"
)

func TestBcryptHasher(t *testing.T) {
	hasher := security.NewBcryptHasher()

	t.Run("should never store the password as is", func(t *testing.T) {
		hash, err := hasher.Hash("secret")

		require.NoError(t, err)
		assert.NotEqual(t, "secret", hash)
	})

	t.Run("should match the right password", func(t *testing.T) {
		hash, err := hasher.Hash("secret")
		require.NoError(t, err)

		assert.True(t, hasher.Matches(hash, "secret"))
	})

	t.Run("should not match a wrong password", func(t *testing.T) {
		hash, err := hasher.Hash("secret")
		require.NoError(t, err)

		assert.False(t, hasher.Matches(hash, "wrong"))
	})

	t.Run("should not match a malformed hash", func(t *testing.T) {
		assert.False(t, hasher.Matches("", "secret"))
	})
}
