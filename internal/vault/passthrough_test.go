package vault

import (
	"errors"
	"testing"
)

func TestPassthroughClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewPassthroughClient(nil, nil)
}

func TestPassthroughClient_NilHookReturnsInner(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "val"},
	})
	client := NewPassthroughClient(mock, nil)
	if client != mock {
		t.Fatal("expected inner client to be returned directly when hook is nil")
	}
}

func TestPassthroughClient_CallsHookOnSuccess(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "val"},
	})
	var hookedPath string
	var hookedSecrets map[string]string
	client := NewPassthroughClient(mock, func(p string, s map[string]string) {
		hookedPath = p
		hookedSecrets = s
	})
	secrets, err := client.ReadSecrets("secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hookedPath != "secret/a" {
		t.Errorf("hook path = %q; want %q", hookedPath, "secret/a")
	}
	if hookedSecrets["key"] != "val" {
		t.Errorf("hook secrets[key] = %q; want %q", hookedSecrets["key"], "val")
	}
	if secrets["key"] != "val" {
		t.Errorf("returned secrets[key] = %q; want %q", secrets["key"], "val")
	}
}

func TestPassthroughClient_HookReceivesCopy(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"k": "v"},
	})
	var hooked map[string]string
	client := NewPassthroughClient(mock, func(_ string, s map[string]string) {
		hooked = s
		s["injected"] = "yes"
	})
	secrets, _ := client.ReadSecrets("secret/a")
	if _, ok := secrets["injected"]; ok {
		t.Error("hook mutation should not affect returned secrets")
	}
	_ = hooked
}

func TestPassthroughClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/missing", errors.New("not found"))
	called := false
	client := NewPassthroughClient(mock, func(_ string, _ map[string]string) {
		called = true
	})
	_, err := client.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if called {
		t.Error("hook should not be called on error")
	}
}

func TestPassthroughClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(nil)
	var _ SecretReader = NewPassthroughClient(mock, func(_ string, _ map[string]string) {})
}
