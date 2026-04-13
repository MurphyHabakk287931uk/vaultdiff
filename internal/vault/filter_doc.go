// Package vault provides Vault client abstractions used throughout vaultdiff.
//
// # FilterClient
//
// FilterClient is a SecretReader decorator that narrows the set of secret keys
// returned by an underlying client to only those matching one or more glob
// patterns.
//
// Patterns follow the same syntax as [path.Match]:
//
//	"DB_*"       — matches any key starting with "DB_"
//	"*_KEY"      — matches any key ending with "_KEY"
//	"API_KEY"    — exact match
//
// If no patterns are supplied (or all patterns are blank) every key is
// returned, preserving the original behaviour of the wrapped client.
//
// Example usage:
//
//	base := vault.NewClient(cfg)
//	filtered := vault.NewFilterClient(base, []string{"DB_*", "API_KEY"})
//	secrets, err := filtered.ReadSecrets("secret/myapp")
package vault
