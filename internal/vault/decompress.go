package vault

import (
	"fmt"
)

// DecompressClient wraps a SecretReader and decompresses gzip+base64 secret values.
// It is the inverse of CompressClient.
type DecompressClient struct {
	inner SecretReader
}

// NewDecompressClient returns a DecompressClient wrapping inner.
// Panics if inner is nil.
func NewDecompressClient(inner SecretReader) *DecompressClient {
	if inner == nil {
		panic("vault: NewDecompressClient: inner must not be nil")
	}
	return &DecompressClient{inner: inner}
}

// ReadSecrets reads secrets from the inner client and decompresses each value.
func (c *DecompressClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		decompressed, err := decompressValue(v)
		if err != nil {
			return nil, fmt.Errorf("vault: decompress key %q: %w", k, err)
		}
		out[k] = decompressed
	}
	return out, nil
}
