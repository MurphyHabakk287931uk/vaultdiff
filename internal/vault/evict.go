package vault

import (
	"fmt"
	"path"
	"sync"
)

// EvictClient wraps a SecretReader and evicts cached entries matching
// a glob pattern whenever a read is performed on a matching path.
// It is useful when paired with CachedClient to force re-reads of
// related secrets after a known mutation.
type EvictClient struct {
	mu       sync.Mutex
	inner    SecretReader
	patterns []string
	evicted  map[string]struct{}
}

// NewEvictClient returns an EvictClient that records paths matching
// any of the provided glob patterns as evicted. Panics if inner is nil.
func NewEvictClient(inner SecretReader, patterns ...string) *EvictClient {
	if inner == nil {
		panic("evict: inner client must not be nil")
	}
	clean := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if p != "" {
			clean = append(clean, p)
		}
	}
	return &EvictClient{inner: inner, patterns: clean, evicted: make(map[string]struct{})}
}

// ReadSecrets delegates to the inner client. If the requested path
// matches one of the eviction patterns the path is recorded as evicted.
func (c *EvictClient) ReadSecrets(p string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(p)
	if err != nil {
		return nil, err
	}
	if c.matches(p) {
		c.mu.Lock()
		c.evicted[p] = struct{}{}
		c.mu.Unlock()
	}
	return secrets, nil
}

// Evicted returns a sorted snapshot of all paths that have been evicted.
func (c *EvictClient) Evicted() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]string, 0, len(c.evicted))
	for k := range c.evicted {
		out = append(out, k)
	}
	return out
}

// Reset clears the eviction record.
func (c *EvictClient) Reset() {
	c.mu.Lock()
	c.evicted = make(map[string]struct{})
	c.mu.Unlock()
}

func (c *EvictClient) matches(p string) bool {
	for _, pat := range c.patterns {
		matched, err := path.Match(pat, p)
		if err != nil {
			_ = fmt.Sprintf("evict: bad pattern %q: %v", pat, err)
			continue
		}
		if matched {
			return true
		}
	}
	return false
}
