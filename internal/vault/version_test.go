package vault

import (
	"errors"
	"testing"
)

func TestVersionedClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewVersionedClient(nil, "v1")
}

func TestVersionedClient_PanicsOnEmptyVersion(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewVersionedClient(NewMockClient(nil, nil), "")
}

func TestVersionedClient_AppendsVersion(t *testing.T) {
	var got string
	mock := &captureMock{fn: func(path string) (map[string]string, error) {
		got = path
		return map[string]string{"key": "val"}, nil
	}}
	c := NewVersionedClient(mock, "v2")
	secrets, err := c.ReadSecrets("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret/myapp/v2" {
		t.Errorf("got path %q, want %q", got, "secret/myapp/v2")
	}
	if secrets["key"] != "val" {
		t.Errorf("unexpected secrets: %v", secrets)
	}
}

func TestVersionedClient_StripsTrailingSlash(t *testing.T) {
	var got string
	mock := &captureMock{fn: func(path string) (map[string]string, error) {
		got = path
		return nil, nil
	}}
	c := NewVersionedClient(mock, "v1")
	_, _ = c.ReadSecrets("secret/myapp/")
	if got != "secret/myapp/v1" {
		t.Errorf("got path %q, want %q", got, "secret/myapp/v1")
	}
}

func TestVersionedClient_PropagatesError(t *testing.T) {
	want := errors.New("boom")
	mock := &captureMock{fn: func(path string) (map[string]string, error) {
		return nil, want
	}}
	c := NewVersionedClient(mock, "v1")
	_, err := c.ReadSecrets("secret/myapp")
	if err != want {
		t.Errorf("got %v, want %v", err, want)
	}
}

func TestVersionedClient_Version(t *testing.T) {
	c := NewVersionedClient(NewMockClient(nil, nil), "v3")
	if c.Version() != "v3" {
		t.Errorf("got %q, want v3", c.Version())
	}
}

func TestVersionedClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewVersionedClient(NewMockClient(nil, nil), "v1")
}

// captureMock is a minimal SecretReader that delegates to a function.
type captureMock struct{ fn func(string) (map[string]string, error) }

func (c *captureMock) ReadSecrets(path string) (map[string]string, error) { return c.fn(path) }
