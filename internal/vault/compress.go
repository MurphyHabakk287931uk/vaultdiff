package vault

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
)

// CompressClient wraps a SecretReader and compresses secret values using gzip+base64.
// This is useful when storing large secrets that exceed Vault's value size limits.
type CompressClient struct {
	inner SecretReader
}

// NewCompressClient returns a CompressClient wrapping inner.
// Panics if inner is nil.
func NewCompressClient(inner SecretReader) *CompressClient {
	if inner == nil {
		panic("vault: NewCompressClient: inner must not be nil")
	}
	return &CompressClient{inner: inner}
}

// ReadSecrets reads secrets from the inner client and compresses each value.
func (c *CompressClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		compressed, err := compressValue(v)
		if err != nil {
			return nil, fmt.Errorf("vault: compress key %q: %w", k, err)
		}
		out[k] = compressed
	}
	return out, nil
}

func compressValue(s string) (string, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := io.WriteString(w, s); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func decompressValue(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	defer r.Close()
	out, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
