// Package vault provides clients for reading secrets from HashiCorp Vault.
//
// # Clients
//
// Use NewClient to create a base Vault client authenticated via token.
// The SecretReader interface abstracts secret retrieval for testing and composition.
//
// # KV Engine Adapters
//
// NewKV1Client and NewKV2Client wrap a SecretReader to rewrite paths according
// to the KV v1 and KV v2 API conventions respectively, scoped to configured mounts.
//
// # Engine Selection
//
// ParseEngineType parses a user-supplied string ("kv1" or "kv2") into an EngineType.
// NewEngineClient selects and constructs the appropriate adapter automatically.
//
// # Mocking
//
// NewMockClient returns a SecretReader backed by an in-memory map, suitable for
// use in unit tests without a live Vault instance.
package vault
