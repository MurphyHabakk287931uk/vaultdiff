// Package vault provides the SplitClient, which reads secrets from two
// independent SecretReader instances and merges their output into a single
// map using configurable namespace prefixes.
//
// This is useful when you want to compare secrets from two environments
// (e.g. production vs staging) using the standard diff pipeline without
// making two separate Vault calls at the runner level.
//
// Example:
//
//	prod := vault.NewClient(prodCfg)
//	staging := vault.NewClient(stagingCfg)
//	split := vault.NewSplitClient(prod, "prod", staging, "staging")
//	secrets, err := split.ReadSecrets("secret/myapp")
//	// secrets["prod/API_KEY"]     = "..."
//	// secrets["staging/API_KEY"] = "..."
package vault
