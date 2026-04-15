// Package vault – ChainClient
//
// ChainClient composes an ordered list of SecretReader implementations
// and resolves secrets by trying each one in sequence.
//
// Typical usage
//
//	cache  := vault.NewCachedClient(primary)
//	replica := vault.NewClient(replicaConfig)
//
//	chain := vault.NewChainClient(cache, primary, replica)
//	secrets, err := chain.ReadSecrets("secret/myapp")
//
// Behaviour
//
//   - Clients are tried in the order they are provided.
//   - The first successful read is returned immediately.
//   - If all clients fail, the error from the last client is returned.
//   - Panics at construction time if the client list is empty.
//
// ChainClient implements the SecretReader interface and can itself be
// wrapped by any other decorator in this package.
package vault
