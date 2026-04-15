package vault

import (
	"context"
	"fmt"
)

// MockClient is an in-memory SecretReader for use in tests.
type MockClient struct {
	data   map[string]map[string]string
	errors map[string]error
}

// NewMockClient creates a MockClient pre-loaded with data.
// data maps path -> key/value secrets; nil is treated as an empty store.
func NewMockClient(data map[string]map[string]string) *MockClient {
	clone := make(map[string]map[string]string)
	for path, secrets := range data {
		secCopy := make(map[string]string, len(secrets))
		for k, v := range secrets {
			secCopy[k] = v
		}
		clone[path] = secCopy
	}
	return &MockClient{
		data:   clone,
		errors: make(map[string]error),
	}
}

// SetError registers an error to return when path is read.
func (m *MockClient) SetError(path string, err error) {
	m.errors[path] = err
}

// ReadSecrets returns the secrets registered for path.
// Returns a "not found" error when the path is absent and no explicit error is set.
func (m *MockClient) ReadSecrets(_ context.Context, path string) (map[string]string, error) {
	if err, ok := m.errors[path]; ok {
		return nil, err
	}
	secrets, ok := m.data[path]
	if !ok {
		return nil, fmt.Errorf("secret not found: %s", path)
	}
	// return a defensive copy so callers cannot mutate internal state
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return copy, nil
}

// Put adds or replaces secrets at path (useful for test setup).
func (m *MockClient) Put(path string, secrets map[string]string) {
	secCopy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		secCopy[k] = v
	}
	m.data[path] = secCopy
}
