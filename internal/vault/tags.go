package vault

import (
	"fmt"
	"strings"
)

// TaggedSecrets wraps a secrets map with an associated set of string tags.
type TaggedSecrets struct {
	Secrets map[string]string
	Tags    map[string]string
}

// tagsClient wraps a SecretReader and injects static metadata tags into
// every result returned. Tags are merged with the secret values under a
// configurable prefix (default "_tag.").
type tagsClient struct {
	inner  SecretReader
	tags   map[string]string
	prefix string
}

// NewTagsClient returns a SecretReader that annotates every successful read
// with the provided static tags. If prefix is empty, "_tag." is used.
// Panics if inner is nil or tags is nil.
func NewTagsClient(inner SecretReader, tags map[string]string, prefix string) SecretReader {
	if inner == nil {
		panic("vault: NewTagsClient: inner must not be nil")
	}
	if tags == nil {
		panic("vault: NewTagsClient: tags must not be nil")
	}
	if prefix == "" {
		prefix = "_tag."
	}
	// Normalise prefix to end with a dot.
	if !strings.HasSuffix(prefix, ".") {
		prefix += "."
	}
	copy := make(map[string]string, len(tags))
	for k, v := range tags {
		copy[k] = v
	}
	return &tagsClient{inner: inner, tags: copy, prefix: prefix}
}

func (c *tagsClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(secrets)+len(c.tags))
	for k, v := range secrets {
		result[k] = v
	}
	for k, v := range c.tags {
		key := fmt.Sprintf("%s%s", c.prefix, k)
		result[key] = v
	}
	return result, nil
}
