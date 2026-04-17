package vault

// PassthroughClient returns secrets unchanged but calls an optional hook
// after every successful read. It is useful for side-effects such as
// logging, metrics emission, or test assertions without altering the
// secret data.
type PassthroughClient struct {
	inner  SecretReader
	onRead func(path string, secrets map[string]string)
}

// NewPassthroughClient wraps inner and calls hook after every successful
// ReadSecrets call. If hook is nil the client behaves identically to inner.
// Panics if inner is nil.
func NewPassthroughClient(inner SecretReader, hook func(string, map[string]string)) SecretReader {
	if inner == nil {
		panic("passthrough: inner client must not be nil")
	}
	if hook == nil {
		return inner
	}
	return &PassthroughClient{inner: inner, onRead: hook}
}

func (c *PassthroughClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	c.onRead(path, copy)
	return secrets, nil
}
