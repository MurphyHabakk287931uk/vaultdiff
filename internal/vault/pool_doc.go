// Package vault — Pool
//
// Pool provides a named registry of SecretReader clients. It is useful
// when a single vaultdiff run must compare secrets drawn from multiple
// independent Vault clusters or authentication contexts.
//
// # Registration
//
//	 pool := vault.NewPool()
//	 pool.Register("prod",    prodClient)
//	 pool.Register("staging", stagingClient)
//
// # Dispatched reads
//
// Pool itself implements SecretReader. The path is expected to start
// with the client name followed by a slash:
//
//	 secrets, err := pool.ReadSecrets("prod/secret/myapp")
//	 // forwards to prodClient.ReadSecrets("secret/myapp")
//
// # Concurrency
//
// Register, Remove, Get, and ReadSecrets are all safe for concurrent
// use. A read-write mutex protects the internal client map.
package vault
