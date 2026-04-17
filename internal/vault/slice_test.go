package vault

import (
	"errors"
	"testing"
)

func TestSliceClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewSliceClient(nil, "", "secret/a")
}

func TestSliceClient_PanicsOnEmptyPaths(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty paths")
		}
	}()
	NewSliceClient(NewMockClient(nil), "")
}

func TestSliceClient_SinglePath(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"foo": "bar"},
	})
	c := NewSliceClient(mock, "", "secret/a")
	got, err := c.ReadSecrets("ignored")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["foo"] != "bar" {
		t.Errorf("expected foo=bar, got %q", got["foo"])
	}
}

func TestSliceClient_MergesMultiplePaths(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"key1": "v1"},
		"secret/b": {"key2": "v2"},
	})
	c := NewSliceClient(mock, "", "secret/a", "secret/b")
	got, err := c.ReadSecrets("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key1"] != "v1" || got["key2"] != "v2" {
		t.Errorf("unexpected merged result: %v", got)
	}
}

func TestSliceClient_LaterPathWinsOnConflict(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "first"},
		"secret/b": {"key": "second"},
	})
	c := NewSliceClient(mock, "", "secret/a", "secret/b")
	got, err := c.ReadSecrets("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "second" {
		t.Errorf("expected second, got %q", got["key"])
	}
}

func TestSliceClient_AppliesPrefix(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"token": "abc"},
	})
	c := NewSliceClient(mock, "db", "secret/a")
	got, err := c.ReadSecrets("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["db.token"] != "abc" {
		t.Errorf("expected db.token=abc, got %v", got)
	}
}

func TestSliceClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	mock := NewMockClient(nil)
	mock.SetError("secret/a", sentinel)
	c := NewSliceClient(mock, "", "secret/a")
	_, err := c.ReadSecrets("")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestSliceClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewSliceClient(NewMockClient(nil), "", "secret/a")
}
