package vault

import (
	"testing"
)

func TestNewClient_MissingToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")

	_, err := NewClient(Config{})
	if err == nil {
		t.Fatal("expected error when no token provided, got nil")
	}
}

func TestNewClient_WithToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "")

	c, err := NewClient(Config{
		Address: "http://127.0.0.1:8200",
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_TokenFromEnv(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "env-token")
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")

	c, err := NewClient(Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestFlattenData_StringValues(t *testing.T) {
	raw := map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}

	got, err := flattenData(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key1"] != "value1" {
		t.Errorf("key1: expected %q, got %q", "value1", got["key1"])
	}
	if got["key2"] != "value2" {
		t.Errorf("key2: expected %q, got %q", "value2", got["key2"])
	}
}

func TestFlattenData_NilValue(t *testing.T) {
	raw := map[string]interface{}{
		"empty": nil,
	}

	got, err := flattenData(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["empty"] != "" {
		t.Errorf("expected empty string for nil value, got %q", got["empty"])
	}
}

func TestFlattenData_NonStringValue(t *testing.T) {
	raw := map[string]interface{}{
		"count": 42,
	}

	got, err := flattenData(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["count"] != "42" {
		t.Errorf("expected \"42\", got %q", got["count"])
	}
}
