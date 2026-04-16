package vault

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// EncryptClient wraps a SecretReader and encrypts all secret values using
// AES-GCM before returning them. The key must be 16, 24, or 32 bytes.
type EncryptClient struct {
	inner SecretReader
	gcm   cipher.AEAD
}

// NewEncryptClient returns an EncryptClient that encrypts values with the
// provided AES key. Panics if inner is nil or the key is invalid.
func NewEncryptClient(inner SecretReader, key []byte) *EncryptClient {
	if inner == nil {
		panic("vault: NewEncryptClient: inner must not be nil")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("vault: NewEncryptClient: invalid key: %v", err))
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Sprintf("vault: NewEncryptClient: gcm init failed: %v", err))
	}
	return &EncryptClient{inner: inner, gcm: gcm}
}

// ReadSecrets reads from the inner client and encrypts each value.
func (c *EncryptClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(ctx, path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		enc, err := c.encrypt(v)
		if err != nil {
			return nil, fmt.Errorf("vault: encrypt key %q: %w", k, err)
		}
		out[k] = enc
	}
	return out, nil
}

func (c *EncryptClient) encrypt(plaintext string) (string, error) {
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := c.gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
