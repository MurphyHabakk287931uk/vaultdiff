package vault

import "fmt"

// CloneClient copies all secrets from one path to another using a source
// SecretReader and a destination SecretWriter. It reads the source path,
// applies an optional transform function to each key, and writes the result
// to the destination path.
//
// This is useful for promoting secrets between environments or namespaces.

// SecretWriter is the write side of a Vault client.
type SecretWriter interface {
	WriteSecrets(path string, data map[string]string) error
}

// CloneClient combines reading from a source and writing to a destination.
type CloneClient struct {
	src SecretReader
	dst SecretWriter
	transformFn func(key string) string
}

// NewCloneClient creates a CloneClient. If transformFn is nil, keys are
// copied verbatim.
func NewCloneClient(src SecretReader, dst SecretWriter, transformFn func(string) string) *CloneClient {
	if src == nil {
		panic("vault: NewCloneClient: src must not be nil")
	}
	if dst == nil {
		panic("vault: NewCloneClient: dst must not be nil")
	}
	return &CloneClient{src: src, dst: dst, transformFn: transformFn}
}

// Clone reads secrets from srcPath and writes them to dstPath.
// Returns an error if the read or write fails.
func (c *CloneClient) Clone(srcPath, dstPath string) error {
	secrets, err := c.src.ReadSecrets(srcPath)
	if err != nil {
		return fmt.Errorf("clone: read from %q: %w", srcPath, err)
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		key := k
		if c.transformFn != nil {
			key = c.transformFn(k)
		}
		out[key] = v
	}

	if err := c.dst.WriteSecrets(dstPath, out); err != nil {
		return fmt.Errorf("clone: write to %q: %w", dstPath, err)
	}
	return nil
}
