package vault

import (
	"errors"
	"testing"
)

func TestCheckpointStore_RecordAndGet(t *testing.T) {
	client := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	store := NewCheckpointStore()
	if err := store.Record("baseline", "secret/app", client); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp, err := store.Get("baseline")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp.Secrets["key"] != "value" {
		t.Errorf("expected value, got %q", cp.Secrets["key"])
	}
	if cp.Path != "secret/app" {
		t.Errorf("expected path secret/app, got %q", cp.Path)
	}
}

func TestCheckpointStore_EmptyName(t *testing.T) {
	client := NewMockClient(nil)
	store := NewCheckpointStore()
	err := store.Record("", "secret/app", client)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCheckpointStore_GetNotFound(t *testing.T) {
	store := NewCheckpointStore()
	_, err := store.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing checkpoint")
	}
}

func TestCheckpointStore_PropagatesReadError(t *testing.T) {
	client := NewMockClient(nil)
	client.SetError("secret/app", errors.New("vault unavailable"))
	store := NewCheckpointStore()
	err := store.Record("snap", "secret/app", client)
	if err == nil {
		t.Fatal("expected error from inner client")
	}
}

func TestCheckpointStore_Names(t *testing.T) {
	client := NewMockClient(map[string]map[string]string{
		"a": {"x": "1"},
		"b": {"y": "2"},
	})
	store := NewCheckpointStore()
	_ = store.Record("first", "a", client)
	_ = store.Record("second", "b", client)
	names := store.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestCheckpointStore_Delete(t *testing.T) {
	client := NewMockClient(map[string]map[string]string{
		"secret/x": {"k": "v"},
	})
	store := NewCheckpointStore()
	_ = store.Record("snap", "secret/x", client)
	store.Delete("snap")
	_, err := store.Get("snap")
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}
