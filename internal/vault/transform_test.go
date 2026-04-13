package vault

import (
	"errors"
	"testing"
)

func TestNewTransformClient_NilInnerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewTransformClient(nil, KeyUpperTransform())
}

func TestNewTransformClient_NilFnReturnsInner(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	client := NewTransformClient(mock, nil)
	if client != mock {
		t.Fatal("expected inner client to be returned unchanged when fn is nil")
	}
}

func TestTransformClient_AppliesTransform(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"db_pass": "hunter2", "api_key": "abc123"},
	})
	client := NewTransformClient(mock, KeyUpperTransform())

	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_PASS"] != "hunter2" {
		t.Errorf("expected DB_PASS=hunter2, got %q", got["DB_PASS"])
	}
	if got["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", got["API_KEY"])
	}
	if _, ok := got["db_pass"]; ok {
		t.Error("original lower-case key should not be present")
	}
}

func TestTransformClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("read failure")
	mock := NewMockClient(nil)
	mock.SetError("secret/app", sentinel)
	client := NewTransformClient(mock, KeyUpperTransform())

	_, err := client.ReadSecrets("secret/app")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestKeyPrefixTransform(t *testing.T) {
	fn := KeyPrefixTransform("APP_")
	input := map[string]string{"token": "xyz", "secret": "abc"}
	out := fn(input)
	if out["APP_token"] != "xyz" {
		t.Errorf("expected APP_token=xyz, got %q", out["APP_token"])
	}
	if out["APP_secret"] != "abc" {
		t.Errorf("expected APP_secret=abc, got %q", out["APP_secret"])
	}
}

func TestTransformClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	var _ SecretReader = NewTransformClient(mock, KeyUpperTransform())
}
