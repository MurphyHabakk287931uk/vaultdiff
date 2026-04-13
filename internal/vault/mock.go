package vault

import "fmt"

// SecretReader is the interface satisfied by all Vault client wrappers.
type SecretReader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// MockClient is an in-memory SecretReader used in tests.
type MockClient struct {
	data   map[string]map[string]string
	errors map[string]error
}

// NewMockClient creates a MockClient pre-loaded with the given data.
// Pass nil to start with an empty store.
func NewMockClient(data map[string]map[string]string) *MockClient {
	copy := make(map[string]map[string]string)
	for k, v := range data {
		inner := make(map[string]string)
		for ik, iv := range v {
			inner[ik] = iv
		}
		copy[k] = inner
	}
	return &MockClient{
		data:   copy,
		errors: make(map[string]error),
	}
}

// SetError registers a forced error for the given path.
func (m *MockClient) SetError(path string, err error) {
	m.errors[path] = err
}

// ReadSecrets returns the secrets stored at path, or an error if one was set.
func (m *MockClient) ReadSecrets(path string) (map[string]string, error) {
	if err, ok := m.errors[path]; ok {
		return nil, err
	}
	secrets, ok := m.data[path]
	if !ok {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	result := make(map[string]string)
	for k, v := range secrets {
		result[k] = v
	}
	return result, nil
}
