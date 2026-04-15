// Package vault provides a collection of composable SecretReader decorators.
//
// # SelectClient
//
// SelectClient is a decorator that restricts the keys returned by a
// SecretReader to an explicit allow-list supplied at construction time.
//
// Use it when you want to expose only a known subset of keys from a broader
// secret path — for example, to surface only DATABASE_URL and API_KEY from a
// path that contains many additional fields.
//
// Example:
//
//	client := vault.NewSelectClient(base, "DATABASE_URL", "API_KEY")
//	secrets, err := client.ReadSecrets(ctx, "myapp/config")
//	// secrets contains at most {"DATABASE_URL": "...", "API_KEY": "..."}
//
// Passing no keys (or only empty strings) disables filtering and returns the
// inner client directly, so there is no overhead in the zero-key case.
package vault
