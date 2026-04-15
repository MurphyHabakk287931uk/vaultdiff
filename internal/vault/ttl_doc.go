// Package vault — TTLClient
//
// TTLClient enforces a maximum age on secrets read from an upstream
// SecretReader. Once a secret has been held longer than the configured
// MaxAge the client returns a stale-error on the next access, signalling
// to the caller that the value should be re-fetched.
//
// Typical usage — pair TTLClient with CachedClient so that the cache
// serves fast in-memory reads while TTLClient gates how long any single
// entry may live:
//
//	inner  := vault.NewClient(cfg)
//	cached := vault.NewCachedClient(inner)
//	ttl    := vault.NewTTLClient(cached, 5*time.Minute)
//
// When ttl.ReadSecrets returns a stale error the caller should call
// cached.Invalidate(path) and retry; TTLClient.Invalidate is a
// convenience helper that does the same on its own lastFetch map.
//
// A zero or negative MaxAge disables TTL enforcement entirely — all
// reads are forwarded to the inner client unconditionally.
package vault
