// Package vault provides Vault client wrappers used by vaultdiff.
//
// # ScopeClient
//
// ScopeClient restricts a SecretReader to a fixed path prefix, rejecting
// any read attempt that targets a path outside the configured scope. This
// is useful when composing clients that should only ever touch a specific
// environment subtree (e.g. "secret/prod").
//
// Usage:
//
//	client := vault.NewScopeClient(inner, "secret/prod")
//	// Only paths under secret/prod are forwarded to inner.
//	// All other paths return *vault.ErrOutOfScope.
//
// An empty scope string is a no-op: the inner client is returned as-is.
// A nil inner client panics at construction time.
package vault
