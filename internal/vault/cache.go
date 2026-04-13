package vault

import (
	"sync"
)

// CachedClient wraps a SecretReader and caches results in memory
// to avoid redundant Vault API calls during a single diff run.
type CachedClient struct {
	inner SecretReader
	mu    sync.RWMutex
	cache map[string]map[string]string
}

// NewCachedClient wraps the given SecretReader with an in-memory cache.
func NewCachedClient(inner SecretReader) *CachedClient {
	return &CachedClient{
		inner: inner,
		cache: make(map[string]map[string]string),
	}
}

// ReadSecrets returns cached secrets for path if available, otherwise
// delegates to the underlying client and stores the result.
func (c *CachedClient) ReadSecrets(path string) (map[string]string, error) {
	c.mu.RLock()
	if cached, ok := c.cache[path]; ok {
		c.mu.RUnlock()
		return cached, nil
	}
	c.mu.RUnlock()

	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.cache[path] = secrets
	c.mu.Unlock()

	return secrets, nil
}

// Invalidate removes a single path from the cache.
func (c *CachedClient) Invalidate(path string) {
	c.mu.Lock()
	delete(c.cache, path)
	c.mu.Unlock()
}

// InvalidateAll clears the entire cache.
func (c *CachedClient) InvalidateAll() {
	c.mu.Lock()
	c.cache = make(map[string]map[string]string)
	c.mu.Unlock()
}

// CacheSize returns the number of paths currently cached.
func (c *CachedClient) CacheSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}
