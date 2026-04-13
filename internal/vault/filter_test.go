package vault

import (
	"errors"
	"testing"
)

func TestFilterClient_NoPatterns_ReturnsAll(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t", "API_KEY": "abc123"},
	})
	client := NewFilterClient(mock, nil)
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestFilterClient_ExactPattern(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t", "API_KEY": "abc123", "PORT": "5432"},
	})
	client := NewFilterClient(mock, []string{"DB_PASS"})
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d", len(got))
	}
	if got["DB_PASS"] != "s3cr3t" {
		t.Errorf("unexpected value for DB_PASS: %q", got["DB_PASS"])
	}
}

func TestFilterClient_GlobPattern(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t", "DB_HOST": "localhost", "API_KEY": "abc123"},
	})
	client := NewFilterClient(mock, []string{"DB_*"})
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
	if _, ok := got["API_KEY"]; ok {
		t.Error("API_KEY should have been filtered out")
	}
}

func TestFilterClient_MultiplePatterns(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t", "API_KEY": "abc123", "PORT": "5432"},
	})
	client := NewFilterClient(mock, []string{"DB_PASS", "PORT"})
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestFilterClient_EmptyPatternsStripped(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "s3cr3t"},
	})
	client := NewFilterClient(mock, []string{"  ", "", "DB_PASS"})
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("expected 1 key, got %d", len(got))
	}
}

func TestFilterClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/app", errors.New("permission denied"))
	client := NewFilterClient(mock, []string{"DB_*"})
	_, err := client.ReadSecrets("secret/app")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFilterClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewFilterClient(NewMockClient(nil), nil)
}
