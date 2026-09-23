package security_test

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/security"
)

func TestTemporaryPasswordGenerator(t *testing.T) {
	t.Run("generates different passwords long enough to be accepted", func(t *testing.T) {
		generator := security.NewTemporaryPasswordGenerator()

		first, err := generator.Generate()
		require.NoError(t, err)
		second, err := generator.Generate()
		require.NoError(t, err)

		assert.GreaterOrEqual(t, utf8.RuneCountInString(first), 8)
		assert.NotEqual(t, first, second)
	})
}
