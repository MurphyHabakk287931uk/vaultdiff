package vault

import "fmt"

// MockClient is an in-memory SecretReader for use in tests.
type MockClient struct {
	data   map[string]map[string]string
	errors map[string]error
}

// NewMockClient returns a MockClient pre-loaded with data.
// data maps path -> key/value secrets. A nil map is valid (empty store).
func NewMockClient(data map[string]map[string]string) *MockClient {
	// deep-copy to prevent callers mutating internal state
	copy := make(map[string]map[string]string, len(data))
	for path, secrets := range data {
		inner := make(map[string]string, len(secrets))
		for k, v := range secrets {
			inner[k] = v
		}
		copy[path] = inner
	}
	return &MockClient{
		data:   copy,
		errors: make(map[string]error),
	}
}

// SetError registers a fixed error to return for path.
func (m *MockClient) SetError(path string, err error) {
	m.errors[path] = err
}

// ReadSecrets returns the secrets stored at path, or an error if one is set.
func (m *MockClient) ReadSecrets(path string) (map[string]string, error) {
	if err, ok := m.errors[path]; ok {
		return nil, err
	}
	secrets, ok := m.data[path]
	if !ok {
		return nil, fmt.Errorf("path not found: %s", path)
	}
	// return a copy so callers cannot mutate internal state
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	return out, nil
}
