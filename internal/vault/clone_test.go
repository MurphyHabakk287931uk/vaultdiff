package vault

import (
	"errors"
	"strings"
	"testing"
)

// mockWriter is a simple in-memory SecretWriter for tests.
type mockWriter struct {
	data    map[string]map[string]string
	wantErr error
}

func (m *mockWriter) WriteSecrets(path string, data map[string]string) error {
	if m.wantErr != nil {
		return m.wantErr
	}
	if m.data == nil {
		m.data = make(map[string]map[string]string)
	}
	m.data[path] = data
	return nil
}

func TestCloneClient_CopiesSecrets(t *testing.T) {
	src := NewMockClient(map[string]map[string]string{
		"secret/src": {"key": "value", "token": "abc"},
	})
	dst := &mockWriter{}

	c := NewCloneClient(src, dst, nil)
	if err := c.Clone("secret/src", "secret/dst"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := dst.data["secret/dst"]
	if got["key"] != "value" || got["token"] != "abc" {
		t.Errorf("unexpected dst data: %v", got)
	}
}

func TestCloneClient_AppliesTransform(t *testing.T) {
	src := NewMockClient(map[string]map[string]string{
		"secret/src": {"db_pass": "s3cr3t"},
	})
	dst := &mockWriter{}

	c := NewCloneClient(src, dst, strings.ToUpper)
	if err := c.Clone("secret/src", "secret/dst"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := dst.data["secret/dst"]
	if got["DB_PASS"] != "s3cr3t" {
		t.Errorf("expected DB_PASS key, got: %v", got)
	}
}

func TestCloneClient_PropagatesReadError(t *testing.T) {
	src := NewMockClient(nil)
	src.SetError("secret/missing", errors.New("not found"))
	dst := &mockWriter{}

	c := NewCloneClient(src, dst, nil)
	err := c.Clone("secret/missing", "secret/dst")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "read from") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestCloneClient_PropagatesWriteError(t *testing.T) {
	src := NewMockClient(map[string]map[string]string{
		"secret/src": {"k": "v"},
	})
	dst := &mockWriter{wantErr: errors.New("permission denied")}

	c := NewCloneClient(src, dst, nil)
	err := c.Clone("secret/src", "secret/dst")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "write to") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestNewCloneClient_PanicsOnNilSrc(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on nil src")
		}
	}()
	NewCloneClient(nil, &mockWriter{}, nil)
}

func TestNewCloneClient_PanicsOnNilDst(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on nil dst")
		}
	}()
	src := NewMockClient(nil)
	NewCloneClient(src, nil, nil)
}
