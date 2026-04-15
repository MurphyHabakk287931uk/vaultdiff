package vault

import (
	"errors"
	"testing"
)

func TestMaskClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewMaskClient(nil, nil, "")
}

func TestMaskClient_NoPatterns_ReturnsUnmasked(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"password": "s3cr3t", "user": "admin"},
	})
	client := NewMaskClient(mock, nil, "")
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["password"] != "s3cr3t" {
		t.Errorf("expected original value, got %q", got["password"])
	}
}

func TestMaskClient_MasksMatchingKey(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"password": "s3cr3t", "user": "admin"},
	})
	client := NewMaskClient(mock, []string{"password"}, "***")
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["password"] != "***" {
		t.Errorf("expected masked value, got %q", got["password"])
	}
	if got["user"] != "admin" {
		t.Errorf("expected unmasked user, got %q", got["user"])
	}
}

func TestMaskClient_GlobPattern(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"db_password": "abc", "db_user": "root", "api_key": "xyz"},
	})
	client := NewMaskClient(mock, []string{"*_password", "api_*"}, "REDACTED")
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["db_password"] != "REDACTED" {
		t.Errorf("db_password should be masked, got %q", got["db_password"])
	}
	if got["api_key"] != "REDACTED" {
		t.Errorf("api_key should be masked, got %q", got["api_key"])
	}
	if got["db_user"] != "root" {
		t.Errorf("db_user should be unmasked, got %q", got["db_user"])
	}
}

func TestMaskClient_DefaultPlaceholder(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"token": "abc123"},
	})
	client := NewMaskClient(mock, []string{"token"}, "")
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["token"] != "***" {
		t.Errorf("expected default placeholder ***, got %q", got["token"])
	}
}

func TestMaskClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/app", errors.New("vault unavailable"))
	client := NewMaskClient(mock, []string{"password"}, "")
	_, err := client.ReadSecrets("secret/app")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestMaskClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewMaskClient(NewMockClient(nil), nil, "")
}
