package vault_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestMultiPathClient_SinglePath(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/a": {"key1": "val1"},
	})
	client := vault.NewMultiPathClient(mock)

	result, err := client.ReadMerged([]string{"secret/a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key1"] != "val1" {
		t.Errorf("expected key1=val1, got %q", result["key1"])
	}
}

func TestMultiPathClient_MergePaths(t *testing.T) {
	mock := vault.NewMockClient(map[string]map[string]string{
		"secret/a": {"key1": "val1", "shared": "from-a"},
		"secret/b": {"key2": "val2", "shared": "from-b"},
	})
	client := vault.NewMultiPathClient(mock)

	result, err := client.ReadMerged([]string{"secret/a", "secret/b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key1"] != "val1" {
		t.Errorf("expected key1=val1, got %q", result["key1"])
	}
	if result["key2"] != "val2" {
		t.Errorf("expected key2=val2, got %q", result["key2"])
	}
	// later path wins
	if result["shared"] != "from-b" {
		t.Errorf("expected shared=from-b (later path wins), got %q", result["shared"])
	}
}

func TestMultiPathClient_EmptyPaths(t *testing.T) {
	mock := vault.NewMockClient(nil)
	client := vault.NewMultiPathClient(mock)

	_, err := client.ReadMerged([]string{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestMultiPathClient_PropagatesError(t *testing.T) {
	mock := vault.NewMockClient(nil)
	mock.SetError("secret/bad", errors.New("permission denied"))
	client := vault.NewMultiPathClient(mock)

	_, err := client.ReadMerged([]string{"secret/bad"})
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestSplitPaths_CommaSeparated(t *testing.T) {
	paths := vault.SplitPaths("secret/a, secret/b ,secret/c")
	if len(paths) != 3 {
		t.Fatalf("expected 3 paths, got %d", len(paths))
	}
	if paths[0] != "secret/a" || paths[1] != "secret/b" || paths[2] != "secret/c" {
		t.Errorf("unexpected paths: %v", paths)
	}
}

func TestSplitPaths_Empty(t *testing.T) {
	paths := vault.SplitPaths("")
	if len(paths) != 0 {
		t.Errorf("expected 0 paths, got %d", len(paths))
	}
}

func TestSplitPaths_SkipsBlanks(t *testing.T) {
	paths := vault.SplitPaths("secret/a,,secret/b")
	if len(paths) != 2 {
		t.Errorf("expected 2 paths, got %d: %v", len(paths), paths)
	}
}
