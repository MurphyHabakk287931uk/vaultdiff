package vault

// ChainClient tries each inner client in order, returning the first
// successful result. Unlike FallbackClient, ChainClient accepts an
// arbitrary number of readers and is useful when building layered
// resolution strategies (e.g. cache → primary → replica).
type ChainClient struct {
	clients []SecretReader
}

// NewChainClient returns a ChainClient that delegates to each reader
// in the supplied slice in order. Panics if no readers are provided.
func NewChainClient(clients ...SecretReader) *ChainClient {
	if len(clients) == 0 {
		panic("vault: NewChainClient requires at least one client")
	}
	return &ChainClient{clients: clients}
}

// ReadSecrets iterates through the chain and returns the first
// successful response. If every client returns an error the last
// error is returned to the caller.
func (c *ChainClient) ReadSecrets(path string) (map[string]string, error) {
	var lastErr error
	for _, r := range c.clients {
		secrets, err := r.ReadSecrets(path)
		if err == nil {
			return secrets, nil
		}
		lastErr = err
	}
	return nil, lastErr
}
