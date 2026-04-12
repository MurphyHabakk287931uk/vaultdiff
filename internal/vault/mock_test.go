package vault

import (
	"errors"
	"testing"
)

func TestMockClient_ReadSecrets_Found(t *testing.T) {
	m := NewMockClient()
	m.SetSecret("secret/data/app", map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	})

	got, err := m.ReadSecrets("secret/data/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_PASS"] != "hunter2" {
		t.Errorf("DB_PASS: expected %q, got %q", "hunter2", got["DB_PASS"])
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("API_KEY: expected %q, got %q", "abc123", got["API_KEY"])
	}
}

func TestMockClient_ReadSecrets_NotFound(t *testing.T) {
	m := NewMockClient()

	_, err := m.ReadSecrets("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for missing path, got nil")
	}
}

func TestMockClient_ReadSecrets_Error(t *testing.T) {
	m := NewMockClient()
	sentinel := errors.New("permission denied")
	m.SetError("secret/data/restricted", sentinel)

	_, err := m.ReadSecrets("secret/data/restricted")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestMockClient_IsolatesInternalState(t *testing.T) {
	m := NewMockClient()
	m.SetSecret("secret/data/app", map[string]string{"KEY": "original"})

	got, _ := m.ReadSecrets("secret/data/app")
	got["KEY"] = "mutated"

	again, _ := m.ReadSecrets("secret/data/app")
	if again["KEY"] != "original" {
		t.Errorf("internal state mutated: expected %q, got %q", "original", again["KEY"])
	}
}

func TestMockClient_ImplementsSecretReader(t *testing.T) {
	// Compile-time check that MockClient satisfies the SecretReader interface.
	var _ SecretReader = (*MockClient)(nil)
	var _ SecretReader = (*Client)(nil)
}
