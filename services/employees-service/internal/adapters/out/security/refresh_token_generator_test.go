package security_test

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/security"
)

func TestRefreshTokenGenerator(t *testing.T) {
	generator := security.NewRefreshTokenGenerator()

	t.Run("should generate unpredictable 256-bit URL-safe tokens", func(t *testing.T) {
		first, err := generator.Generate()
		require.NoError(t, err)
		second, err := generator.Generate()
		require.NoError(t, err)

		raw, err := base64.RawURLEncoding.DecodeString(first)
		require.NoError(t, err)
		assert.Len(t, raw, 32)
		assert.NotEqual(t, first, second)
	})

	t.Run("should hash deterministically without exposing the token", func(t *testing.T) {
		token, err := generator.Generate()
		require.NoError(t, err)

		assert.Equal(t, generator.Hash(token), generator.Hash(token))
		assert.NotEqual(t, token, generator.Hash(token))
		assert.NotEqual(t, generator.Hash(token), generator.Hash(token+"x"))
	})
}
