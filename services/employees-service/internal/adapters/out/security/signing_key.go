package security

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var ErrInvalidSigningKey = errors.New("signing key must be an Ed25519 private key in PKCS8 PEM format")

func LoadSigningKey(encoded string) (key ed25519.PrivateKey, generated bool, err error) {
	if encoded == "" {
		_, key, err = ed25519.GenerateKey(rand.Reader)
		return key, true, err
	}

	key, err = parseSigningKey(encoded)
	return key, false, err
}

func parseSigningKey(encoded string) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode([]byte(encoded))
	if block == nil {
		return nil, ErrInvalidSigningKey
	}

	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if key, ok := parsed.(ed25519.PrivateKey); err == nil && ok {
		return key, nil
	}
	return nil, ErrInvalidSigningKey
}
