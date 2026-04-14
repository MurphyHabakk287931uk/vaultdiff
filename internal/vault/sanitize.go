package vault

import (
	"regexp"
	"strings"
)

// SanitizeFunc is a function that sanitizes a secret key or value.
type SanitizeFunc func(key, value string) string

// SanitizeClient wraps a SecretReader and applies a SanitizeFunc to all
// returned secret values before passing them to the caller.
type SanitizeClient struct {
	inner    SecretReader
	sanitize SanitizeFunc
}

// NewSanitizeClient returns a SanitizeClient that applies fn to every value
// returned by inner. If fn is nil the inner client is returned unwrapped.
func NewSanitizeClient(inner SecretReader, fn SanitizeFunc) SecretReader {
	if inner == nil {
		panic("vault: NewSanitizeClient: inner must not be nil")
	}
	if fn == nil {
		return inner
	}
	return &SanitizeClient{inner: inner, sanitize: fn}
}

// ReadSecrets implements SecretReader.
func (c *SanitizeClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = c.sanitize(k, v)
	}
	return out, nil
}

// TrimSpaceTransform returns a SanitizeFunc that trims leading and trailing
// whitespace from every secret value.
func TrimSpaceTransform() SanitizeFunc {
	return func(_, value string) string {
		return strings.TrimSpace(value)
	}
}

// RedactPatternTransform returns a SanitizeFunc that replaces any substring in
// a value that matches pattern with replacement.
func RedactPatternTransform(pattern, replacement string) SanitizeFunc {
	re := regexp.MustCompile(pattern)
	return func(_, value string) string {
		return re.ReplaceAllString(value, replacement)
	}
}
