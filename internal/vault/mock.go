package vault

import (
	"context"
	"fmt"
	"sync"
)

// MockClient is an in-memory SecretReader for use in tests.
type MockClient struct {
	mu      sync.RWMutex
	data    map[string]map[string]string
	errors  map[string]error
}

// NewMockClient returns a MockClient pre-populated with the given data.
// The data map is deep-copied so callers retain ownership of the original.
func NewMockClient(data map[string]map[string]string) *MockClient {
	c := &MockClient{
		data:   make(map[string]map[string]string, len(data)),
		errors: make(map[string]error),
	}
	for path, secrets := range data {
		copy := make(map[string]string, len(secrets))
		for k, v := range secrets {
			copy[k] = v
		}
		c.data[path] = copy
	}
	return c
}

// SetSecrets replaces the secrets stored at path.
func (m *MockClient) SetSecrets(path string, secrets map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	m.data[path] = copy
}

// SetError configures ReadSecrets to return err for the given path.
func (m *MockClient) SetError(path string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors[path] = err
}

// ReadSecrets implements SecretReader.
func (m *MockClient) ReadSecrets(_ context.Context, path string) (map[string]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if err, ok := m.errors[path]; ok {
		return nil, err
	}
	secrets, ok := m.data[path]
	if !ok {
		return nil, fmt.Errorf("secret not found: %s", path)
	}
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return copy, nil
}
