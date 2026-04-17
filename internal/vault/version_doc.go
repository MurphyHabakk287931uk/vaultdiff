// Package vault — VersionedClient
//
// VersionedClient appends a fixed version label to every secret path before
// delegating to an inner SecretReader. This is useful when comparing two
// versions of the same secret tree:
//
//	src := vault.NewVersionedClient(base, "v1")
//	dst := vault.NewVersionedClient(base, "v2")
//	// runner.Run("secret/myapp", "secret/myapp") will read
//	// secret/myapp/v1 and secret/myapp/v2 respectively.
//
// The version string is appended as a trailing path segment after any
// trailing slash is stripped from the caller-supplied path.
package vault
