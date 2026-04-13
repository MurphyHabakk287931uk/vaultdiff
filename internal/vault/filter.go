package vault

import (
	"path"
	"strings"
)

// FilterClient wraps a SecretReader and filters keys from the returned secrets
// based on a list of glob patterns. Only keys matching at least one pattern are
// included. If no patterns are provided all keys are returned unchanged.
type FilterClient struct {
	inner    SecretReader
	patterns []string
}

// NewFilterClient returns a FilterClient that delegates reads to inner and
// retains only secret keys whose names match one of the provided glob patterns.
func NewFilterClient(inner SecretReader, patterns []string) *FilterClient {
	cleaned := make([]string, 0, len(patterns))
	for _, p := range patterns {
		p = strings.TrimSpace(p)
		if p != "" {
			cleaned = append(cleaned, p)
		}
	}
	return &FilterClient{inner: inner, patterns: cleaned}
}

// ReadSecrets reads secrets from the underlying client and removes any keys
// that do not match the configured patterns.
func (f *FilterClient) ReadSecrets(secretPath string) (map[string]string, error) {
	secrets, err := f.inner.ReadSecrets(secretPath)
	if err != nil {
		return nil, err
	}
	if len(f.patterns) == 0 {
		return secrets, nil
	}
	filtered := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if matchesAny(k, f.patterns) {
			filtered[k] = v
		}
	}
	return filtered, nil
}

// matchesAny reports whether key matches at least one of the provided glob
// patterns using path.Match semantics.
func matchesAny(key string, patterns []string) bool {
	for _, p := range patterns {
		matched, err := path.Match(p, key)
		if err == nil && matched {
			return true
		}
	}
	return false
}
