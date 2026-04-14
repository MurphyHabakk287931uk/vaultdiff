package vault

import (
	"errors"
	"testing"
)

func TestSnapshotClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewSnapshotClient(nil)
}

func TestSnapshotClient_RecordsSnapshot(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	client := NewSnapshotClient(mock)

	_, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	snaps := client.Snapshots()
	if len(snaps) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(snaps))
	}
	if snaps[0].Path != "secret/app" {
		t.Errorf("expected path %q, got %q", "secret/app", snaps[0].Path)
	}
	if snaps[0].Secrets["key"] != "value" {
		t.Errorf("expected secret value %q, got %q", "value", snaps[0].Secrets["key"])
	}
}

func TestSnapshotClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/missing", errors.New("not found"))
	client := NewSnapshotClient(mock)

	_, err := client.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if len(client.Snapshots()) != 0 {
		t.Error("expected no snapshots on error")
	}
}

func TestSnapshotClient_Latest(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"a": "1"},
	})
	client := NewSnapshotClient(mock)

	client.ReadSecrets("secret/app") //nolint:errcheck
	client.ReadSecrets("secret/app") //nolint:errcheck

	snap, err := client.Latest("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap == nil {
		t.Fatal("expected a snapshot")
	}
}

func TestSnapshotClient_Latest_NotFound(t *testing.T) {
	client := NewSnapshotClient(NewMockClient(nil))
	_, err := client.Latest("secret/nope")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestSnapshot_SortedKeys(t *testing.T) {
	snap := &Snapshot{Secrets: map[string]string{"z": "1", "a": "2", "m": "3"}}
	keys := snap.SortedKeys()
	expected := []string{"a", "m", "z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestSnapshotClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewSnapshotClient(NewMockClient(nil))
}
