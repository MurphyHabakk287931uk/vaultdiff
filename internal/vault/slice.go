package vault

import "fmt"

// SliceClient fans out a single ReadSecrets call across multiple fixed paths
// and returns a merged view keyed by a configurable prefix per path.
//
// Unlike BatchClient (which accepts paths at call time), SliceClient is
// configured once at construction and always reads the same set of paths.
type SliceClient struct {
	inner   SecretReader
	paths   []string
	prefix  string
}

// NewSliceClient returns a SliceClient that reads from each path in paths and
// merges the results. Keys from later paths overwrite earlier ones on conflict.
// If prefix is non-empty it is prepended to every key as "<prefix>.<key>".
// Panics if inner is nil or paths is empty.
func NewSliceClient(inner SecretReader, prefix string, paths ...string) *SliceClient {
	if inner == nil {
		panic("vault: NewSliceClient: inner must not be nil")
	}
	if len(paths) == 0 {
		panic("vault: NewSliceClient: at least one path is required")
	}
	return &SliceClient{inner: inner, paths: paths, prefix: prefix}
}

// ReadSecrets reads each configured path in order and merges the results.
func (c *SliceClient) ReadSecrets(_ string) (map[string]string, error) {
	merged := make(map[string]string)
	for _, p := range c.paths {
		secrets, err := c.inner.ReadSecrets(p)
		if err != nil {
			return nil, fmt.Errorf("vault: slice read %q: %w", p, err)
		}
		for k, v := range secrets {
			key := k
			if c.prefix != "" {
				key = c.prefix + "." + k
			}
			merged[key] = v
		}
	}
	return merged, nil
}
