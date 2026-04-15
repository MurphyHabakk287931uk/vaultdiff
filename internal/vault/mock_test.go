package vault

import (
	"context"
	"errors"
	"testing"
)

func TestMockClient_ReadSecrets_Found(t *testing.T) {
	m := NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "value"},
	})
	secrets, err := m.ReadSecrets(context.Background(), "secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "value" {
		t.Fatalf("expected value, got %s", secrets["key"])
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
	sentinel := errors.New("injected")
	m := NewMockClient(nil)
	m.SetError("secret/a", sentinel)
	_, err := m.ReadSecrets(context.Background(), "secret/a")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestMockClient_IsolatesInternalState(t *testing.T) {
	orig := map[string]map[string]string{"secret/a": {"k": "v"}}
	m := NewMockClient(orig)
	orig["secret/a"]["k"] = "mutated"

	secrets, err := m.ReadSecrets(context.Background(), "secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["k"] != "v" {
		t.Fatalf("expected original value 'v', got %s", secrets["k"])
	}
}

func TestMockClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewMockClient(nil)
}

func TestMockClient_SetSecrets_UpdatesLive(t *testing.T) {
	m := NewMockClient(nil)
	m.SetSecrets("secret/b", map[string]string{"x": "1"})
	secrets, err := m.ReadSecrets(context.Background(), "secret/b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["x"] != "1" {
		t.Fatalf("expected '1', got %s", secrets["x"])
	}
}
