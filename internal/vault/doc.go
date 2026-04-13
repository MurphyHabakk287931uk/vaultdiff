// Package vault provides clients for reading secrets from HashiCorp Vault.
//
// # Clients
//
// [NewClient] creates a standard Vault API client backed by the official SDK.
// It reads VAULT_TOKEN from the environment or accepts an explicit token.
//
// [NewKV2Client] wraps any [SecretReader] and rewrites paths for KV v2 secrets
// engines by injecting the "/data/" infix required by the Vault HTTP API.
//
// [NewKV1Client] wraps any [SecretReader] for KV v1 secrets engines where
// paths are used directly without any infix. It also strips accidental "/data/"
// segments that callers may include by mistake.
//
// # Mocking
//
// [NewMockClient] provides an in-memory implementation of [SecretReader]
// suitable for use in tests. Secrets are pre-loaded from a map at construction
// time and the client's internal state is never mutated after creation.
package vault
