package security_test

import (
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/franciscoHonorat/Sys-Called/services/employees-service/internal/adapters/out/security"
)

func TestLoadSigningKey(t *testing.T) {
	t.Run("should load an Ed25519 key from PEM", func(t *testing.T) {
		_, private := newKeyPair(t)
		der, err := x509.MarshalPKCS8PrivateKey(private)
		require.NoError(t, err)
		encoded := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))

		key, generated, err := security.LoadSigningKey(encoded)

		require.NoError(t, err)
		assert.False(t, generated)
		assert.Equal(t, private, key)
	})

	t.Run("should generate an ephemeral key when none is configured", func(t *testing.T) {
		key, generated, err := security.LoadSigningKey("")

		require.NoError(t, err)
		assert.True(t, generated)
		assert.NotEmpty(t, key)
	})

	t.Run("should fail for an invalid PEM instead of running without a usable key", func(t *testing.T) {
		_, _, err := security.LoadSigningKey("not a pem")

		assert.Error(t, err)
	})
}
