package vault

import "fmt"

// MockClient implements a fake Vault backend for use in tests.
// It stores secrets as path -> key/value maps.
type MockClient struct {
	Secrets map[string]map[string]string
	Errors  map[string]error
}

// NewMockClient returns an initialised MockClient.
func NewMockClient() *MockClient {
	return &MockClient{
		Secrets: make(map[string]map[string]string),
		Errors:  make(map[string]error),
	}
}

// SetSecret registers a set of key/value pairs at the given path.
func (m *MockClient) SetSecret(path string, data map[string]string) {
	m.Secrets[path] = data
}

// SetError causes ReadSecrets to return an error for the given path.
func (m *MockClient) SetError(path string, err error) {
	m.Errors[path] = err
}

// ReadSecrets returns the secrets stored at path, mimicking Client.ReadSecrets.
func (m *MockClient) ReadSecrets(path string) (map[string]string, error) {
	if err, ok := m.Errors[path]; ok {
		return nil, err
	}
	data, ok := m.Secrets[path]
	if !ok {
		return nil, fmt.Errorf("vault: path %q not found or empty", path)
	}
	// Return a copy so callers cannot mutate internal state.
	out := make(map[string]string, len(data))
	for k, v := range data {
		out[k] = v
	}
	return out, nil
}

// SecretReader is the interface satisfied by both Client and MockClient.
type SecretReader interface {
	ReadSecrets(path string) (map[string]string, error)
}
