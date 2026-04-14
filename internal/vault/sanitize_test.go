package vault

import (
	"errors"
	"testing"
)

func TestNewSanitizeClient_NilInnerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewSanitizeClient(nil, TrimSpaceTransform())
}

func TestNewSanitizeClient_NilFnReturnsInner(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	got := NewSanitizeClient(mock, nil)
	if got != mock {
		t.Fatal("expected inner client to be returned unchanged")
	}
}

func TestSanitizeClient_AppliesTrimSpace(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "  hello world  "},
	})
	client := NewSanitizeClient(mock, TrimSpaceTransform())
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "hello world" {
		t.Errorf("expected trimmed value, got %q", got["key"])
	}
}

func TestSanitizeClient_AppliesRedactPattern(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/db": {"password": "super-secret-123"},
	})
	client := NewSanitizeClient(mock, RedactPatternTransform(`\d+`, "***"))
	got, err := client.ReadSecrets("secret/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["password"] != "super-secret-***" {
		t.Errorf("unexpected value: %q", got["password"])
	}
}

func TestSanitizeClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/fail", errors.New("read error"))
	client := NewSanitizeClient(mock, TrimSpaceTransform())
	_, err := client.ReadSecrets("secret/fail")
	if err == nil {
		t.Fatal("expected error to propagate")
	}
}

func TestSanitizeClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	var _ SecretReader = NewSanitizeClient(mock, TrimSpaceTransform())
}

func TestTrimSpaceTransform_LeavesCleanValues(t *testing.T) {
	fn := TrimSpaceTransform()
	if got := fn("k", "clean"); got != "clean" {
		t.Errorf("expected clean, got %q", got)
	}
}

func TestRedactPatternTransform_NoMatch(t *testing.T) {
	fn := RedactPatternTransform(`\d+`, "***")
	if got := fn("k", "no-digits-here"); got != "no-digits-here" {
		t.Errorf("unexpected modification: %q", got)
	}
}
