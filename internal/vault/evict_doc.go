// Package vault provides the EvictClient decorator.
//
// # EvictClient
//
// EvictClient wraps any SecretReader and records which paths were read
// that match one or more glob patterns. This is useful when you need to
// track which cache entries should be invalidated after an external
// mutation, without coupling the eviction logic to the cache itself.
//
// Usage:
//
//	base := vault.NewClient(cfg)
//	cached := vault.NewCachedClient(base)
//	evict := vault.NewEvictClient(cached, "secret/app/*")
//
//	// After a write elsewhere:
//	for _, path := range evict.Evicted() {
//	    cached.Invalidate(path)
//	}
//	evict.Reset()
//
// Patterns follow the syntax of path.Match from the standard library.
package vault
