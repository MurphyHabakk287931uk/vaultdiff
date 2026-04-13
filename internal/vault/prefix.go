package vault

import (
	"fmt"
	"strings"
)

// PrefixClient wraps a SecretReader and prepends a fixed prefix to every path
// before delegating to the inner client. This is useful when comparing secrets
// across different base paths without modifying call sites.
type PrefixClient struct {
	inner  SecretReader
	prefix string
}

// NewPrefixClient returns a PrefixClient that prepends prefix to every path
// passed to ReadSecrets. The prefix is normalised: leading/trailing slashes are
// stripped and a single separating slash is inserted between prefix and path.
// An empty prefix returns the inner client unwrapped.
func NewPrefixClient(inner SecretReader, prefix string) (SecretReader, error) {
	if inner == nil {
		return nil, fmt.Errorf("vault: NewPrefixClient: inner client must not be nil")
	}
	prefix = strings.Trim(prefix, "/")
	if prefix == "" {
		return inner, nil
	}
	return &PrefixClient{inner: inner, prefix: prefix}, nil
}

// ReadSecrets prepends the configured prefix to path and delegates to the
// inner SecretReader.
func (p *PrefixClient) ReadSecrets(path string) (map[string]string, error) {
	path = strings.TrimPrefix(path, "/")
	full := p.prefix + "/" + path
	return p.inner.ReadSecrets(full)
}

// Prefix returns the normalised prefix string used by this client.
func (p *PrefixClient) Prefix() string {
	return p.prefix
}
