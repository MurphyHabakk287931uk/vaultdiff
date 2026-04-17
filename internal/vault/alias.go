package vault

// AliasClient maps logical secret paths to real Vault paths using a
// static alias table. This is useful when callers use short, environment-
// agnostic names that differ from the actual Vault path layout.
type AliasClient struct {
	inner   SecretReader
	aliases map[string]string
}

// NewAliasClient wraps inner with a path alias table. If a requested path
// has an entry in aliases, the mapped path is used instead. Paths with no
// alias are forwarded unchanged.
//
// Panics if inner is nil or aliases is nil.
func NewAliasClient(inner SecretReader, aliases map[string]string) *AliasClient {
	if inner == nil {
		panic("vault: NewAliasClient: inner must not be nil")
	}
	if aliases == nil {
		panic("vault: NewAliasClient: aliases must not be nil")
	}
	copy := make(map[string]string, len(aliases))
	for k, v := range aliases {
		copy[k] = v
	}
	return &AliasClient{inner: inner, aliases: copy}
}

// ReadSecrets resolves path through the alias table and delegates to inner.
func (a *AliasClient) ReadSecrets(path string) (map[string]string, error) {
	if resolved, ok := a.aliases[path]; ok {
		path = resolved
	}
	return a.inner.ReadSecrets(path)
}

// Alias returns the resolved path for the given logical name, or the
// original path if no alias is registered.
func (a *AliasClient) Alias(path string) string {
	if resolved, ok := a.aliases[path]; ok {
		return resolved
	}
	return path
}
