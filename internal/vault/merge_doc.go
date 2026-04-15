// Package vault — MergeClient
//
// MergeClient fans out a single ReadSecrets call to multiple inner clients and
// merges their results into one map. This is useful when secrets for a logical
// path are spread across several Vault mounts or namespaces and you want a
// unified view for diffing.
//
// Precedence
//
// Clients are merged in the order they are passed to NewMergeClient. If two
// clients return the same key, the value from the *later* client wins. This
// mirrors the common "override" pattern where a base environment is overridden
// by a more specific one.
//
// Error handling
//
// A not-found error from an individual client is silently skipped — it simply
// means that client has no data for that path. Any other error is returned
// immediately and no merged result is produced. If *all* clients return
// not-found the MergeClient itself returns a not-found error.
//
// Example
//
//	 base := vault.NewKV2Client(raw, "secret")
//	 overrides := vault.NewKV2Client(raw, "overrides")
//	 merged := vault.NewMergeClient(base, overrides)
//	 secrets, err := merged.ReadSecrets("myapp/config")
package vault
