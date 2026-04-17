// Package vault — SliceClient
//
// SliceClient reads a fixed set of Vault paths and merges their secrets into a
// single flat map. It is useful when a logical "environment" is spread across
// several Vault paths that should be treated as one unit for diffing purposes.
//
// Construction
//
//	client := vault.NewSliceClient(
//		inner,
//		"db",              // optional key prefix
//		"secret/common",
//		"secret/app/prod",
//	)
//
// The path argument passed to ReadSecrets is intentionally ignored; the paths
// are fixed at construction time. This makes SliceClient compatible with any
// wrapper that calls ReadSecrets with a single path string.
//
// Conflict resolution
//
// When the same key appears in more than one path the value from the last path
// in the list wins. Use a consistent ordering to keep behaviour deterministic.
package vault
