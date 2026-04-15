package vault

import (
	"context"
	"errors"
	"testing"
)

func TestMockClient_ReadSecrets_Found(t *testing.T) {
	m := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	secrets, err := m.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "value" {
		t.Errorf("expected 'value', got %q", secrets["key"])
	}
}

func TestMockClient_ReadSecrets_NotFound(t *testing.T) {
	m := NewMockClient(nil)
	_, err := m.ReadSecrets(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestMockClient_ReadSecrets_Error(t *testing.T) {
	m := NewMockClient(nil)
	sentinel := errors.New("vault unavailable")
	m.SetError("secret/app", sentinel)
	_, err := m.ReadSecrets(context.Background(), "secret/app")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestMockClient_IsolatesInternalState(t *testing.T) {
	m := NewMockClient(map[string]map[string]string{
		"secret/app": {"k": "original"},
	})
	secrets, _ := m.ReadSecrets(context.Background(), "secret/app")
	secrets["k"] = "mutated"

	secrets2, _ := m.ReadSecrets(context.Background(), "secret/app")
	if secrets2["k"] != "original" {
		t.Errorf("internal state was mutated; expected 'original', got %q", secrets2["k"])
	}
}

func TestMockClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewMockClient(nil)
}

func TestMockClient_Put(t *testing.T) {
	m := NewMockClient(nil)
	m.Put("secret/new", map[string]string{"a": "1"})
	secrets, err := m.ReadSecrets(context.Background(), "secret/new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["a"] != "1" {
		t.Errorf("expected '1', got %q", secrets["a"])
	}
}
