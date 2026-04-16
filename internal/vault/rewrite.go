package vault

import "strings"

// RewriteRule maps an input path prefix to a replacement prefix.
type RewriteRule struct {
	From string
	To   string
}

// rewriteClient rewrites secret paths before delegating to the inner client.
type rewriteClient struct {
	inner SecretReader
	rules []RewriteRule
}

// NewRewriteClient returns a SecretReader that rewrites paths using the given
// rules before forwarding to inner. Rules are applied in order; the first
// matching rule wins.
func NewRewriteClient(inner SecretReader, rules []RewriteRule) SecretReader {
	if inner == nil {
		panic("rewriteClient: inner must not be nil")
	}
	if len(rules) == 0 {
		return inner
	}
	return &rewriteClient{inner: inner, rules: rules}
}

func (c *rewriteClient) ReadSecrets(path string) (map[string]string, error) {
	return c.inner.ReadSecrets(c.rewrite(path))
}

func (c *rewriteClient) rewrite(path string) string {
	for _, r := range c.rules {
		if strings.HasPrefix(path, r.From) {
			return r.To + strings.TrimPrefix(path, r.From)
		}
	}
	return path
}
