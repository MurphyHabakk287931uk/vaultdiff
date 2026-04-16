// Package vault provides SchemaClient, a SecretReader decorator that
// enforces value-format constraints using regular expressions.
//
// # Usage
//
//	rules := []vault.SchemaRule{
//		{Key: "port",    Pattern: regexp.MustCompile(`^\d+$`)},
//		{Key: "api_key", Pattern: regexp.MustCompile(`^[A-Za-z0-9]{32}$`)},
//	}
//	client := vault.NewSchemaClient(inner, rules)
//
// If a secret value does not satisfy its rule's pattern, ReadSecrets
// returns an error describing the violation. Keys that are absent from
// the returned secret map are silently skipped.
//
// Passing a nil or empty rules slice returns the inner client unchanged,
// incurring no overhead at call time.
package vault
