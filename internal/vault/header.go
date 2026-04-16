package vault

import (
	"context"
	"fmt"
)

// HeaderClient injects custom HTTP-style metadata headers into secret results
// as synthetic keys. This is useful for tracing the origin of secrets when
// merging from multiple sources.
type HeaderClient struct {
	inner   SecretReader
	headers map[string]string
	prefix  string
}

// NewHeaderClient wraps inner and injects the provided headers into every
// successful ReadSecrets response. Keys are prefixed with prefix (e.g.
// "_meta"). Panics if inner is nil or headers is nil.
func NewHeaderClient(inner SecretReader, headers map[string]string, prefix string) *HeaderClient {
	if inner == nil {
		panic("vault: NewHeaderClient: inner must not be nil")
	}
	if headers == nil {
		panic("vault: NewHeaderClient: headers must not be nil")
	}
	if prefix == "" {
		prefix = "_header"
	}
	return &HeaderClient{inner: inner, headers: headers, prefix: prefix}
}

func (c *HeaderClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(ctx, path)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(secrets)+len(c.headers))
	for k, v := range secrets {
		out[k] = v
	}
	for k, v := range c.headers {
		out[fmt.Sprintf("%s.%s", c.prefix, k)] = v
	}
	return out, nil
}
