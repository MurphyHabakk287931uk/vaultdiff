package vault

import (
	"errors"
	"testing"
)

func TestRewriteClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewRewriteClient(nil, []RewriteRule{{From: "a", To: "b"}})
}

func TestRewriteClient_NoRules_ReturnsInner(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"k": "v"},
	})
	client := NewRewriteClient(mock, nil)
	if client != mock {
		t.Fatal("expected inner client returned directly")
	}
}

func TestRewriteClient_RewritesMatchingPrefix(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"prod/db": {"password": "secret"},
	})
	client := NewRewriteClient(mock, []RewriteRule{
		{From: "staging/", To: "prod/"},
	})
	got, err := client.ReadSecrets("staging/db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["password"] != "secret" {
		t.Fatalf("expected secret, got %q", got["password"])
	}
}

func TestRewriteClient_NoMatchUsesOriginalPath(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"other/path": {"key": "val"},
	})
	client := NewRewriteClient(mock, []RewriteRule{
		{From: "staging/", To: "prod/"},
	})
	got, err := client.ReadSecrets("other/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "val" {
		t.Fatalf("expected val, got %q", got["key"])
	}
}

func TestRewriteClient_FirstRuleWins(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"prod/svc": {"x": "1"},
		"canary/svc": {"x": "2"},
	})
	client := NewRewriteClient(mock, []RewriteRule{
		{From: "dev/", To: "prod/"},
		{From: "dev/", To: "canary/"},
	})
	got, err := client.ReadSecrets("dev/svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["x"] != "1" {
		t.Fatalf("expected first rule to win, got %q", got["x"])
	}
}

func TestRewriteClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("prod/db", errors.New("vault unavailable"))
	client := NewRewriteClient(mock, []RewriteRule{
		{From: "staging/", To: "prod/"},
	})
	_, err := client.ReadSecrets("staging/db")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRewriteClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(nil)
	var _ SecretReader = NewRewriteClient(mock, []RewriteRule{{From: "a", To: "b"}})
}
