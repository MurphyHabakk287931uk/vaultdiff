// Package vault provides EnrichClient, a decorator that injects static
// key/value pairs into every secret read response.
//
// Use cases:
//
//   - Attach environment labels (e.g. "env": "production") to all secrets
//     so downstream consumers always have context.
//   - Inject build metadata or deployment identifiers without modifying
//     the underlying Vault data.
//
// Example:
//
//	extra := map[string]string{"env": "staging", "region": "eu-west-1"}
//	client := vault.NewEnrichClient(inner, extra)
//	secrets, err := client.ReadSecrets(ctx, "secret/myapp")
//	// secrets will always contain "env" and "region" keys.
//
// Extra keys always win over values returned by the inner client.
package vault
