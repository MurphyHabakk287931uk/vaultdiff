package vault

import (
	"fmt"
	"strings"
)

// KV1Client wraps a SecretReader and rewrites paths for KV v1 secrets engines.
// KV v1 paths are read directly without the "/data/" infix used by KV v2.
type KV1Client struct {
	inner  SecretReader
	mounts []string
}

// NewKV1Client returns a KV1Client that rewrites paths for the given mount
// prefixes. If no mounts are provided, "secret" is used as the default.
func NewKV1Client(inner SecretReader, mounts ...string) *KV1Client {
	if len(mounts) == 0 {
		mounts = []string{"secret"}
	}
	return &KV1Client{inner: inner, mounts: mounts}
}

// ReadSecrets reads secrets from a KV v1 path. Unlike KV v2, no path rewriting
// is needed — the path is used as-is relative to the mount.
func (c *KV1Client) ReadSecrets(path string) (map[string]string, error) {
	normalized := c.normalizePath(path)
	return c.inner.ReadSecrets(normalized)
}

// normalizePath ensures the path does not accidentally contain a "/data/"
// segment that would be appropriate only for KV v2 mounts.
func (c *KV1Client) normalizePath(path string) string {
	for _, mount := range c.mounts {
		prefix := strings.TrimRight(mount, "/") + "/"
		if strings.HasPrefix(path, prefix) {
			rest := strings.TrimPrefix(path, prefix)
			// Strip accidental /data/ infix if present
			rest = strings.TrimPrefix(rest, "data/")
			return fmt.Sprintf("%s%s", prefix, rest)
		}
	}
	return path
}
