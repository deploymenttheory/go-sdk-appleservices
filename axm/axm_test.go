package axm

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebhookServiceConstructors(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	der, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	encoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	path := filepath.Join(t.TempDir(), "test.p8")
	require.NoError(t, os.WriteFile(path, encoded, 0600))
	t.Setenv("APPLE_KEY_ID", "test")
	t.Setenv("APPLE_ISSUER_ID", "test")
	t.Setenv("APPLE_PRIVATE_KEY_PEM", string(encoded))
	for _, tc := range []struct {
		name string
		new  func() (*Client, error)
	}{
		{"direct", func() (*Client, error) { return NewClient("test", "test", key) }},
		{"file", func() (*Client, error) { return NewClientFromFile("test", "test", path) }},
		{"env", func() (*Client, error) { return NewClientFromEnv() }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c, err := tc.new()
			require.NoError(t, err)
			defer c.transport.Close()
			require.NotNil(t, c.AXMAPI.Webhooks)
		})
	}
}
