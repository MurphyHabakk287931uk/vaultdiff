package vault

import (
	"context"
	"fmt"
	"strings"
)

// KV2Client wraps a SecretReader with KV v2 path rewriting support.
// Vault KV v2 secrets are stored under <mount>/data/<path>, but users
// typically provide <mount>/<path>. This adapter rewrites paths transparently.
type KV2Client struct {
	inner  SecretReader
	mounts []string
}

// NewKV2Client returns a KV2Client that rewrites paths for the given KV v2
// mount points. If no mounts are provided, "secret" is used as the default.
func NewKV2Client(inner SecretReader, mounts ...string) *KV2Client {
	if len(mounts) == 0 {
		mounts = []string{"secret"}
	}
	return &KV2Client{inner: inner, mounts: mounts}
}

// ReadSecrets reads secrets from the given path, rewriting it to the KV v2
// data path if it matches a known mount.
func (c *KV2Client) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	rewritten := c.rewritePath(path)
	secrets, err := c.inner.ReadSecrets(ctx, rewritten)
	if err != nil {
		return nil, fmt.Errorf("kv2: read %q (rewritten from %q): %w", rewritten, path, err)
	}
	return secrets, nil
}

// rewritePath inserts "/data/" after the mount segment if the path starts
// with a known KV v2 mount and does not already contain "/data/".
func (c *KV2Client) rewritePath(path string) string {
	path = strings.TrimPrefix(path, "/")
	for _, mount := range c.mounts {
		mount = strings.TrimSuffix(mount, "/")
		prefix := mount + "/"
		if strings.HasPrefix(path, prefix) {
			rest := strings.TrimPrefix(path, prefix)
			if strings.HasPrefix(rest, "data/") {
				// already rewritten
				return path
			}
			return mount + "/data/" + rest
		}
	}
	return path
}
