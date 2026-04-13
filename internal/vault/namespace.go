package vault

import (
	"fmt"
	"strings"
)

// NamespaceClient wraps a SecretReader and prepends a Vault namespace
// (Enterprise feature) to every path before delegating the read.
type NamespaceClient struct {
	inner     SecretReader
	namespace string
}

// NewNamespaceClient returns a NamespaceClient that prefixes every secret
// path with the given namespace segment. The namespace must be non-empty.
func NewNamespaceClient(inner SecretReader, namespace string) (*NamespaceClient, error) {
	namespace = strings.Trim(namespace, "/")
	if namespace == "" {
		return nil, fmt.Errorf("vault: namespace must not be empty")
	}
	return &NamespaceClient{inner: inner, namespace: namespace}, nil
}

// ReadSecrets prepends the configured namespace to path and delegates to the
// inner SecretReader.
func (n *NamespaceClient) ReadSecrets(path string) (map[string]string, error) {
	path = strings.TrimPrefix(path, "/")
	namespacedPath := n.namespace + "/" + path
	return n.inner.ReadSecrets(namespacedPath)
}

// Namespace returns the configured namespace string.
func (n *NamespaceClient) Namespace() string {
	return n.namespace
}
