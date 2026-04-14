// Package vault provides Vault client abstractions used by vaultdiff.
//
// # FallbackClient
//
// FallbackClient wraps two SecretReader implementations — a primary and a
// secondary — and provides transparent failover between them.
//
// Typical use-cases:
//
//   - Comparing a secret path that exists in production but not yet in staging:
//     the staging client falls back to a default/shared secrets mount.
//
//   - Graceful degradation when a Vault cluster is temporarily unavailable and
//     a read-through cache is used as the secondary.
//
// The caller controls when fallback occurs via a predicate function:
//
//	client := vault.NewFallbackClient(
//		prodClient,
//		stagingClient,
//		vault.IsNotFound, // only fall back for missing paths
//	)
//
// Passing nil for the predicate falls back on any error from the primary.
//
// IsNotFound is a built-in predicate that matches errors.Is(err, ErrNotFound),
// allowing network or permission errors to propagate normally.
package vault
