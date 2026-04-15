package vault

import (
	"fmt"
	"strings"
)

// MaskClient wraps a SecretReader and masks the values of secrets whose keys
// match any of the provided patterns. Masked values are replaced with a
// fixed placeholder so that keys are still visible in diff output without
// leaking sensitive data.
type MaskClient struct {
	inner    SecretReader
	patterns []string
	placeholder string
}

// NewMaskClient returns a MaskClient that replaces values for matching keys
// with placeholder. If placeholder is empty, "***" is used. Patterns support
// the same glob syntax as FilterClient (path.Match).
func NewMaskClient(inner SecretReader, patterns []string, placeholder string) *MaskClient {
	if inner == nil {
		panic("vault: NewMaskClient requires a non-nil inner client")
	}
	if placeholder == "" {
		placeholder = "***"
	}
	cleaned := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if t := strings.TrimSpace(p); t != "" {
			cleaned = append(cleaned, t)
		}
	}
	return &MaskClient{
		inner:       inner,
		patterns:    cleaned,
		placeholder: placeholder,
	}
}

// ReadSecrets delegates to the inner client and then masks matching keys.
func (m *MaskClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := m.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	if len(m.patterns) == 0 {
		return secrets, nil
	}
	result := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if matchesAny(k, m.patterns) {
			result[k] = m.placeholder
		} else {
			result[k] = v
		}
	}
	return result, nil
}

// String returns a human-readable description of the client.
func (m *MaskClient) String() string {
	return fmt.Sprintf("MaskClient(patterns=%v, placeholder=%q)", m.patterns, m.placeholder)
}
