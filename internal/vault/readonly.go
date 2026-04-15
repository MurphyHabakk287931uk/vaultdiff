package vault

import (
	"context"
	"errors"
)

// ErrReadOnly is returned when a write operation is attempted on a read-only client.
var ErrReadOnly = errors.New("vault: client is read-only")

// SecretWriter is the interface for writing secrets to a Vault path.
type SecretWriter interface {
	WriteSecrets(ctx context.Context, path string, data map[string]string) error
}

// ReadOnlyClient wraps a SecretReader and rejects any write attempts.
// It is useful for enforcing that diff operations never mutate Vault state.
type ReadOnlyClient struct {
	inner SecretReader
}

// NewReadOnlyClient returns a ReadOnlyClient wrapping inner.
// It panics if inner is nil.
func NewReadOnlyClient(inner SecretReader) *ReadOnlyClient {
	if inner == nil {
		panic("vault: NewReadOnlyClient requires a non-nil inner client")
	}
	return &ReadOnlyClient{inner: inner}
}

// ReadSecrets delegates to the inner client.
func (c *ReadOnlyClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	return c.inner.ReadSecrets(ctx, path)
}

// WriteSecrets always returns ErrReadOnly, preventing any mutation.
func (c *ReadOnlyClient) WriteSecrets(_ context.Context, _ string, _ map[string]string) error {
	return ErrReadOnly
}
