package vault

import (
	"context"
	"fmt"
)

// SelectClient wraps a SecretReader and returns only the keys specified in the
// allow-list. Keys not present in the allow-list are silently dropped from the
// result. If the allow-list is empty, all keys are returned unchanged.
type SelectClient struct {
	inner SecretReader
	keys  map[string]struct{}
}

// NewSelectClient creates a SelectClient that filters secrets to only the
// provided keys. Passing zero keys disables filtering (all keys pass through).
// Panics if inner is nil.
func NewSelectClient(inner SecretReader, keys ...string) SecretReader {
	if inner == nil {
		panic("vault: NewSelectClient: inner client must not be nil")
	}
	if len(keys) == 0 {
		return inner
	}
	allowed := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if k != "" {
			allowed[k] = struct{}{}
		}
	}
	if len(allowed) == 0 {
		return inner
	}
	return &SelectClient{inner: inner, keys: allowed}
}

// ReadSecrets delegates to the inner client then drops any key not present in
// the allow-list. Returns an error if the inner read fails.
func (s *SelectClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	raw, err := s.inner.ReadSecrets(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("select: %w", err)
	}
	out := make(map[string]string, len(s.keys))
	for k, v := range raw {
		if _, ok := s.keys[k]; ok {
			out[k] = v
		}
	}
	return out, nil
}
