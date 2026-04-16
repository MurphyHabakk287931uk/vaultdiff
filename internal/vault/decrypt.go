package vault

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// DecryptClient wraps a SecretReader and decrypts AES-GCM encrypted values
// produced by EncryptClient before returning them.
type DecryptClient struct {
	inner SecretReader
	gcm   cipher.AEAD
}

// NewDecryptClient returns a DecryptClient using the provided AES key.
// Panics if inner is nil or the key is invalid.
func NewDecryptClient(inner SecretReader, key []byte) *DecryptClient {
	if inner == nil {
		panic("vault: NewDecryptClient: inner must not be nil")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(fmt.Sprintf("vault: NewDecryptClient: invalid key: %v", err))
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(fmt.Sprintf("vault: NewDecryptClient: gcm init failed: %v", err))
	}
	return &DecryptClient{inner: inner, gcm: gcm}
}

// ReadSecrets reads from the inner client and decrypts each value.
func (c *DecryptClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(ctx, path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		plain, err := c.decrypt(v)
		if err != nil {
			return nil, fmt.Errorf("vault: decrypt key %q: %w", k, err)
		}
		out[k] = plain
	}
	return out, nil
}

func (c *DecryptClient) decrypt(encoded string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("base64 decode: %w", err)
	}
	ns := c.gcm.NonceSize()
	if len(data) < ns {
		return "", fmt.Errorf("ciphertext too short")
	}
	plaintext, err := c.gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", fmt.Errorf("aes-gcm open: %w", err)
	}
	return string(plaintext), nil
}
