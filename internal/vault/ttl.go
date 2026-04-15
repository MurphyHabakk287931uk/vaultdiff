package vault

import (
	"context"
	"fmt"
	"time"
)

// TTLClient wraps a SecretReader and attaches a maximum age to each cached
// read. If the secret was fetched longer ago than MaxAge, the next call will
// bypass any upstream cache and re-fetch from the inner client.
//
// TTLClient itself does NOT cache — pair it with CachedClient when you want
// both TTL expiry and in-memory caching.
type TTLClient struct {
	inner  SecretReader
	maxAge time.Duration
	now    func() time.Time // injectable for tests

	// lastFetch tracks when each path was last successfully read.
	lastFetch map[string]time.Time
}

// NewTTLClient creates a TTLClient that considers secrets older than maxAge
// stale. A zero or negative maxAge disables TTL enforcement and the inner
// client is always called directly.
func NewTTLClient(inner SecretReader, maxAge time.Duration) *TTLClient {
	if inner == nil {
		panic("vault: NewTTLClient: inner client must not be nil")
	}
	return &TTLClient{
		inner:     inner,
		maxAge:    maxAge,
		now:       time.Now,
		lastFetch: make(map[string]time.Time),
	}
}

// ReadSecrets delegates to the inner client and records the fetch timestamp.
// If the path was fetched within maxAge the call is still forwarded — TTLClient
// relies on the caller (e.g. CachedClient) to short-circuit repeated reads;
// its role is to signal staleness via the returned error.
func (t *TTLClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	if t.maxAge > 0 {
		if fetched, ok := t.lastFetch[path]; ok {
			if t.now().Sub(fetched) > t.maxAge {
				return nil, fmt.Errorf("vault: TTLClient: secret at %q is stale (age > %s)", path, t.maxAge)
			}
		}
	}

	secrets, err := t.inner.ReadSecrets(ctx, path)
	if err != nil {
		return nil, err
	}
	t.lastFetch[path] = t.now()
	return secrets, nil
}

// Invalidate removes the recorded fetch time for path so the next call will
// not be considered stale regardless of age.
func (t *TTLClient) Invalidate(path string) {
	delete(t.lastFetch, path)
}

// InvalidateAll clears all recorded fetch times.
func (t *TTLClient) InvalidateAll() {
	t.lastFetch = make(map[string]time.Time)
}
