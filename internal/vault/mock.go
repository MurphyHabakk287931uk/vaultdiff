package vault

import (
	"context"
	"fmt"
)

// MockClient is an in-memory SecretReader used in tests.
type MockClient struct {
	data   map[string]map[string]string
	errors map[string]error
}

// NewMockClient creates a MockClient pre-loaded with data.
// data maps path → key/value pairs; a nil map is valid (empty store).
func NewMockClient(data map[string]map[string]string) *MockClient {
	copy := make(map[string]map[string]string)
	for k, v := range data {
		row := make(map[string]string, len(v))
		for kk, vv := range v {
			row[kk] = vv
		}
		copy[k] = row
	}
	return &MockClient{
		data:   copy,
		errors: make(map[string]error),
	}
}

// SetError registers a fixed error to be returned for path.
func (m *MockClient) SetError(path string, err error) {
	m.errors[path] = err
}

// ReadSecrets returns the secrets stored at path, or an error if one was
// registered via SetError, or a not-found error if path is unknown.
func (m *MockClient) ReadSecrets(_ context.Context, path string) (map[string]string, error) {
	if err, ok := m.errors[path]; ok {
		return nil, err
	}
	secrets, ok := m.data[path]
	if !ok {
		return nil, fmt.Errorf("path %q not found", path)
	}
	// Return a defensive copy so callers cannot mutate internal state.
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	return out, nil
}

// Put stores or replaces secrets at path (useful for test setup).
func (m *MockClient) Put(path string, secrets map[string]string) {
	row := make(map[string]string, len(secrets))
	for k, v := range secrets {
		row[k] = v
	}
	m.data[path] = row
}
