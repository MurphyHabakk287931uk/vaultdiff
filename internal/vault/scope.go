package vault

import (
	"context"
	"fmt"
	"strings"
)

// ScopeClient wraps a SecretReader and restricts all reads to a fixed
// path scope. Any path that does not begin with the allowed prefix is
// rejected before the inner client is consulted.
type ScopeClient struct {
	inner  SecretReader
	scope  string
}

// ErrOutOfScope is returned when a requested path falls outside the
// configured scope.
type ErrOutOfScope struct {
	Path  string
	Scope string
}

func (e *ErrOutOfScope) Error() string {
	return fmt.Sprintf("path %q is outside allowed scope %q", e.Path, e.Scope)
}

// NewScopeClient returns a ScopeClient that only permits reads from
// paths that begin with scope. An empty scope allows all paths (the
// inner client is returned directly). Panics if inner is nil.
func NewScopeClient(inner SecretReader, scope string) SecretReader {
	if inner == nil {
		panic("vault: NewScopeClient: inner client must not be nil")
	}
	scope = strings.Trim(scope, "/")
	if scope == "" {
		return inner
	}
	return &ScopeClient{inner: inner, scope: scope}
}

// ReadSecrets enforces the scope restriction and delegates to the
// inner client when the path is permitted.
func (s *ScopeClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	norm := strings.TrimLeft(path, "/")
	prefix := s.scope + "/"
	if norm != s.scope && !strings.HasPrefix(norm, prefix) {
		return nil, &ErrOutOfScope{Path: path, Scope: s.scope}
	}
	return s.inner.ReadSecrets(ctx, path)
}
