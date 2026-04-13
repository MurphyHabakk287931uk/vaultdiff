package vault

import (
	"errors"
	"testing"
)

func TestMockClient_ReadSecrets_Found(t *testing.T) {
	m := NewMockClient(map[string]map[string]string{
		"secret/app": {"db_pass": "s3cr3t"},
	})
	got, err := m.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["db_pass"] != "s3cr3t" {
		t.Errorf("expected s3cr3t, got %s", got["db_pass"])
	}
}

func TestMockClient_ReadSecrets_NotFound(t *testing.T) {
	m := NewMockClient(nil)
	_, err := m.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestMockClient_ReadSecrets_Error(t *testing.T) {
	m := NewMockClient(nil)
	m.SetError("secret/forbidden", errors.New("403 forbidden"))
	_, err := m.ReadSecrets("secret/forbidden")
	if err == nil || err.Error() != "403 forbidden" {
		t.Fatalf("expected 403 forbidden, got %v", err)
	}
}

func TestMockClient_IsolatesInternalState(t *testing.T) {
	original := map[string]map[string]string{
		"secret/app": {"key": "original"},
	}
	m := NewMockClient(original)
	// mutate the source map after construction
	original["secret/app"]["key"] = "mutated"

	got, err := m.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "original" {
		t.Errorf("internal state was mutated: got %s", got["key"])
	}
}

func TestMockClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = (*MockClient)(nil)
}
