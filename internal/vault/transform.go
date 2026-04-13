package vault

import "strings"

// TransformFunc is a function that transforms a secret map.
type TransformFunc func(map[string]string) map[string]string

// transformClient applies a TransformFunc to secrets returned by an inner SecretReader.
type transformClient struct {
	inner     SecretReader
	transform TransformFunc
}

// NewTransformClient wraps inner with a client that applies fn to every
// secret map returned by ReadSecrets. If fn is nil the inner client is
// returned unchanged.
func NewTransformClient(inner SecretReader, fn TransformFunc) SecretReader {
	if inner == nil {
		panic("vault: NewTransformClient: inner must not be nil")
	}
	if fn == nil {
		return inner
	}
	return &transformClient{inner: inner, transform: fn}
}

func (c *transformClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	return c.transform(secrets), nil
}

// KeyPrefixTransform returns a TransformFunc that prepends prefix to every key.
func KeyPrefixTransform(prefix string) TransformFunc {
	return func(m map[string]string) map[string]string {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[prefix+k] = v
		}
		return out
	}
}

// KeyUpperTransform returns a TransformFunc that upper-cases every key.
func KeyUpperTransform() TransformFunc {
	return func(m map[string]string) map[string]string {
		out := make(map[string]string, len(m))
		for k, v := range m {
			out[strings.ToUpper(k)] = v
		}
		return out
	}
}
