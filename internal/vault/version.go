package vault

import (
	"fmt"
	"strings"
)

// VersionedClient wraps a SecretReader and appends a version suffix to every
// path before delegating, enabling side-by-side comparison of secret versions
// stored under paths like "secret/myapp/v1" vs "secret/myapp/v2".
type VersionedClient struct {
	inner   SecretReader
	version string
}

// NewVersionedClient returns a VersionedClient that appends version to each
// path. version must be non-empty.
func NewVersionedClient(inner SecretReader, version string) *VersionedClient {
	if inner == nil {
		panic("vault: NewVersionedClient: inner must not be nil")
	}
	if version == "" {
		panic("vault: NewVersionedClient: version must not be empty")
	}
	return &VersionedClient{inner: inner, version: version}
}

// ReadSecrets appends the version as a trailing path segment then delegates.
func (v *VersionedClient) ReadSecrets(path string) (map[string]string, error) {
	versioned := fmt.Sprintf("%s/%s", strings.TrimRight(path, "/"), v.version)
	return v.inner.ReadSecrets(versioned)
}

// Version returns the version label used by this client.
func (v *VersionedClient) Version() string { return v.version }
