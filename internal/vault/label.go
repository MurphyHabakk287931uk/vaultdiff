package vault

import "fmt"

// LabelClient wraps a SecretReader and attaches a human-readable label to
// every path it reads. The label is prepended to the path as a comment-style
// prefix in error messages, making it easier to identify which environment or
// cluster produced a given result in multi-source diffs.
type LabelClient struct {
	inner SecretReader
	label string
}

// NewLabelClient returns a LabelClient that decorates inner with label.
// If inner is nil the function panics. An empty label is allowed and
// effectively makes LabelClient a transparent pass-through.
func NewLabelClient(inner SecretReader, label string) *LabelClient {
	if inner == nil {
		panic("vault: NewLabelClient: inner must not be nil")
	}
	return &LabelClient{inner: inner, label: label}
}

// Label returns the label associated with this client.
func (c *LabelClient) Label() string { return c.label }

// ReadSecrets delegates to the inner client. On error the label is prepended
// to the error message so callers can identify the source of the failure.
func (c *LabelClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		if c.label != "" {
			return nil, fmt.Errorf("[%s] %w", c.label, err)
		}
		return nil, err
	}
	return secrets, nil
}
