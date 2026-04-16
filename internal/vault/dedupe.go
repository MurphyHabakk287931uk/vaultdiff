package vault

import (
	"context"
	"sync"
)

// dedupeCall tracks an in-flight or completed call for a given path.
type dedupeCall struct {
	wg  sync.WaitGroup
	val map[string]string
	err error
}

// DedupeClient wraps a SecretReader and collapses concurrent duplicate reads
// for the same path into a single upstream call. This is useful when multiple
// goroutines may request the same secret simultaneously (e.g. batch reads).
type DedupeClient struct {
	inner SecretReader
	mu    sync.Mutex
	calls map[string]*dedupeCall
}

// NewDedupeClient returns a DedupeClient wrapping inner.
// It panics if inner is nil.
func NewDedupeClient(inner SecretReader) *DedupeClient {
	if inner == nil {
		panic("vault: NewDedupeClient requires a non-nil inner client")
	}
	return &DedupeClient{
		inner: inner,
		calls: make(map[string]*dedupeCall),
	}
}

// ReadSecrets delegates to the inner client, deduplicating concurrent requests
// for the same path. The first caller triggers the real read; subsequent
// callers with the same path block until the first completes and share its result.
func (d *DedupeClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	d.mu.Lock()
	if c, ok := d.calls[path]; ok {
		d.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := &dedupeCall{}
	c.wg.Add(1)
	d.calls[path] = c
	d.mu.Unlock()

	c.val, c.err = d.inner.ReadSecrets(ctx, path)
	c.wg.Done()

	d.mu.Lock()
	delete(d.calls, path)
	d.mu.Unlock()

	return c.val, c.err
}
