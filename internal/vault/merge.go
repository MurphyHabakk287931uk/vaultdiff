package vault

// MergeClient combines secrets from multiple readers under a single path,
// with later clients overriding keys from earlier ones.
type MergeClient struct {
	clients []SecretReader
}

// NewMergeClient returns a MergeClient that merges results from all provided
// clients. Later clients take precedence over earlier ones for duplicate keys.
// Panics if fewer than two clients are provided.
func NewMergeClient(clients ...SecretReader) *MergeClient {
	if len(clients) < 2 {
		panic("vault: NewMergeClient requires at least two clients")
	}
	for i, c := range clients {
		if c == nil {
			panic(fmt.Sprintf("vault: NewMergeClient: client at index %d is nil", i))
		}
	}
	return &MergeClient{clients: clients}
}

// ReadSecrets reads the given path from every inner client and merges the
// resulting maps. If all clients return a not-found error the call returns
// that error. Any other error is returned immediately.
func (m *MergeClient) ReadSecrets(path string) (map[string]string, error) {
	merged := make(map[string]string)
	found := false

	for _, c := range m.clients {
		secrets, err := c.ReadSecrets(path)
		if err != nil {
			if IsNotFound(err) {
				continue
			}
			return nil, err
		}
		found = true
		for k, v := range secrets {
			merged[k] = v
		}
	}

	if !found {
		return nil, fmt.Errorf("secret not found: %s", path)
	}
	return merged, nil
}
