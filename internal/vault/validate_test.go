package vault

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewValidateClient(nil, RequireKeys("k"))
}

func TestValidateClient_PanicsOnNilFn(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil fn")
		}
	}()
	NewValidateClient(NewMockClient(nil, nil), nil)
}

func TestValidateClient_PassesValidSecrets(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"API_KEY": "abc123"},
	}, nil)
	client := NewValidateClient(mock, RequireKeys("API_KEY"))

	secrets, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["API_KEY"] != "abc123" {
		t.Errorf("unexpected value: %q", secrets["API_KEY"])
	}
}

func TestValidateClient_FailsMissingRequiredKey(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"OTHER": "value"},
	}, nil)
	client := NewValidateClient(mock, RequireKeys("API_KEY", "DB_PASS"))

	_, err := client.ReadSecrets("secret/app")
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !IsValidationError(err) {
		t.Errorf("expected *ValidationError, got %T", err)
	}
	if !strings.Contains(err.Error(), "API_KEY") {
		t.Errorf("error should mention missing key: %v", err)
	}
}

func TestValidateClient_PropagatesInnerError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := NewMockClient(nil, sentinel)
	client := NewValidateClient(mock, RequireKeys("k"))

	_, err := client.ReadSecrets("secret/app")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestRejectEmptyValues_DetectsBlankValue(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/cfg": {"KEY": "   "},
	}, nil)
	client := NewValidateClient(mock, RejectEmptyValues())

	_, err := client.ReadSecrets("secret/cfg")
	if !IsValidationError(err) {
		t.Fatalf("expected ValidationError, got %v", err)
	}
	if !strings.Contains(err.Error(), "KEY") {
		t.Errorf("error should mention offending key: %v", err)
	}
}

func TestValidateClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewValidateClient(NewMockClient(nil, nil), RequireKeys())
}

func TestValidationError_ErrorMessage(t *testing.T) {
	ve := &ValidationError{Path: "secret/x", Issues: []string{"missing key \"A\"", "missing key \"B\""}}
	msg := ve.Error()
	if !strings.Contains(msg, "secret/x") {
		t.Errorf("error should contain path: %q", msg)
	}
	if !strings.Contains(msg, "missing key") {
		t.Errorf("error should contain issue text: %q", msg)
	}
}
