package vault

import (
	"context"
	"fmt"
	"sync"
)

// ReloadClient wraps a SecretReader and supports hot-reloading the inner client
// at runtime. Callers can swap the underlying client without restarting the
// process, which is useful when credentials or configuration change.
type ReloadClient struct {
	mu    sync.RWMutex
	inner SecretReader
}

// NewReloadClient returns a ReloadClient backed by inner.
// Panics if inner is nil.
func NewReloadClient(inner SecretReader) *ReloadClient {
	if inner == nil {
		panic("vault: NewReloadClient: inner must not be nil")
	}
	return &ReloadClient{inner: inner}
}

// ReadSecrets delegates to the current inner client.
func (c *ReloadClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	c.mu.RLock()
	inner := c.inner
	c.mu.RUnlock()
	return inner.ReadSecrets(ctx, path)
}

// Reload atomically replaces the inner client.
// Returns an error if next is nil.
func (c *ReloadClient) Reload(next SecretReader) error {
	if next == nil {
		return fmt.Errorf("vault: Reload: replacement client must not be nil")
	}
	c.mu.Lock()
	c.inner = next
	c.mu.Unlock()
	return nil
}

// Current returns the active inner client.
func (c *ReloadClient) Current() SecretReader {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.inner
}
