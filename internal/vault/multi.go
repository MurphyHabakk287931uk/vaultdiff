package vault

import (
	"fmt"
	"strings"
)

// MultiPathClient reads secrets from multiple paths and merges them into a
// single flat map. Keys from later paths override keys from earlier paths.
type MultiPathClient struct {
	reader SecretReader
}

// NewMultiPathClient wraps an existing SecretReader to support multi-path reads.
func NewMultiPathClient(reader SecretReader) *MultiPathClient {
	return &MultiPathClient{reader: reader}
}

// ReadMerged reads secrets from all provided paths and merges the results.
// Keys from later paths take precedence over earlier ones.
// Returns an error if any individual path read fails.
func (m *MultiPathClient) ReadMerged(paths []string) (map[string]string, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("multi: at least one path is required")
	}

	merged := make(map[string]string)

	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}

		secrets, err := m.reader.ReadSecrets(p)
		if err != nil {
			return nil, fmt.Errorf("multi: reading path %q: %w", p, err)
		}

		for k, v := range secrets {
			merged[k] = v
		}
	}

	return merged, nil
}

// SplitPaths splits a comma-separated path string into individual trimmed paths.
func SplitPaths(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
