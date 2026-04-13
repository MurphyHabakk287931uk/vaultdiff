package vault

import (
	"testing"
)

func TestKV1Client_DefaultMount(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/myapp": {"key": "value"},
	})
	client := NewKV1Client(mock)

	got, err := client.ReadSecrets("secret/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("expected 'value', got %q", got["key"])
	}
}

func TestKV1Client_StripsAccidentalDataInfix(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/myapp": {"token": "abc123"},
	})
	client := NewKV1Client(mock)

	// Caller accidentally used a KV v2-style path
	got, err := client.ReadSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["token"] != "abc123" {
		t.Errorf("expected 'abc123', got %q", got["token"])
	}
}

func TestKV1Client_NonMatchingMount(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"kv/myapp": {"pass": "hunter2"},
	})
	client := NewKV1Client(mock, "secret")

	// Path does not match "secret" mount — passed through unchanged
	got, err := client.ReadSecrets("kv/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["pass"] != "hunter2" {
		t.Errorf("expected 'hunter2', got %q", got["pass"])
	}
}

func TestKV1Client_CustomMount(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"kv/svc/config": {"db": "postgres"},
	})
	client := NewKV1Client(mock, "kv")

	got, err := client.ReadSecrets("kv/svc/config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["db"] != "postgres" {
		t.Errorf("expected 'postgres', got %q", got["db"])
	}
}

func TestKV1Client_MultipleMounts(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"x": "1"},
		"kv/b":     {"y": "2"},
	})
	client := NewKV1Client(mock, "secret", "kv")

	if got, err := client.ReadSecrets("secret/a"); err != nil || got["x"] != "1" {
		t.Errorf("secret/a: err=%v got=%v", err, got)
	}
	if got, err := client.ReadSecrets("kv/b"); err != nil || got["y"] != "2" {
		t.Errorf("kv/b: err=%v got=%v", err, got)
	}
}

func TestKV1Client_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(nil)
	var _ SecretReader = NewKV1Client(mock)
}
