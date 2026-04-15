package vault

import (
	"errors"
	"testing"
)

func TestMergeClient_PanicsOnSingleClient(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for single client")
		}
	}()
	NewMergeClient(NewMockClient(nil))
}

func TestMergeClient_PanicsOnNilClient(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil client")
		}
	}()
	NewMergeClient(NewMockClient(nil), nil)
}

func TestMergeClient_MergesDisjointKeys(t *testing.T) {
	a := NewMockClient(map[string]map[string]string{
		"app/cfg": {"host": "localhost"},
	})
	b := NewMockClient(map[string]map[string]string{
		"app/cfg": {"port": "8080"},
	})

	mc := NewMergeClient(a, b)
	got, err := mc.ReadSecrets("app/cfg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["host"] != "localhost" || got["port"] != "8080" {
		t.Errorf("unexpected merged result: %v", got)
	}
}

func TestMergeClient_LaterClientWins(t *testing.T) {
	a := NewMockClient(map[string]map[string]string{
		"app/cfg": {"key": "base"},
	})
	b := NewMockClient(map[string]map[string]string{
		"app/cfg": {"key": "override"},
	})

	mc := NewMergeClient(a, b)
	got, err := mc.ReadSecrets("app/cfg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "override" {
		t.Errorf("expected 'override', got %q", got["key"])
	}
}

func TestMergeClient_SkipsNotFound(t *testing.T) {
	a := NewMockClient(nil) // returns not-found for everything
	b := NewMockClient(map[string]map[string]string{
		"app/cfg": {"key": "value"},
	})

	mc := NewMergeClient(a, b)
	got, err := mc.ReadSecrets("app/cfg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("expected 'value', got %q", got["key"])
	}
}

func TestMergeClient_AllNotFound(t *testing.T) {
	a := NewMockClient(nil)
	b := NewMockClient(nil)

	mc := NewMergeClient(a, b)
	_, err := mc.ReadSecrets("missing/path")
	if err == nil {
		t.Fatal("expected error when all clients return not-found")
	}
}

func TestMergeClient_PropagatesHardError(t *testing.T) {
	hardErr := errors.New("vault unavailable")
	a := NewMockClient(nil)
	mc_err := NewMockClient(nil)
	mc_err.SetError("app/cfg", hardErr)

	mc := NewMergeClient(a, mc_err)
	_, err := mc.ReadSecrets("app/cfg")
	if !errors.Is(err, hardErr) {
		t.Errorf("expected hard error, got: %v", err)
	}
}

func TestMergeClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewMergeClient(NewMockClient(nil), NewMockClient(nil))
}
