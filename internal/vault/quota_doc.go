// Package vault — QuotaClient
//
// NewQuotaClient wraps any SecretReader and enforces an upper bound on the
// total number of ReadSecrets calls that may be made through it.
//
// # Usage
//
//	cfg := vault.QuotaConfig{MaxReads: 100}
//	client := vault.NewQuotaClient(inner, cfg)
//
// Once the limit is reached every subsequent call returns ErrQuotaExceeded.
// Set MaxReads to 0 to disable the quota (the default).
//
// # Thread safety
//
// The internal counter uses sync/atomic and is safe for concurrent use.
//
// # Typical use-cases
//
//   - Preventing runaway diff operations from hammering Vault in CI.
//   - Enforcing a budget when fanning out across many paths via BatchClient.
//   - Smoke-testing that a code path reads only an expected number of secrets.
package vault
