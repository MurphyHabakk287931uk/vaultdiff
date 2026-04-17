package vault

import (
	"errors"
	"sort"
	"testing"
)

func TestEvictClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewEvictClient(nil, "secret/*")
}

func TestEvictClient_NoPatterns_NoEviction(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"k": "v"},
	})
	c := NewEvictClient(mock)
	_, err := c.ReadSecrets("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Evicted()) != 0 {
		t.Fatalf("expected no evictions, got %v", c.Evicted())
	}
}

func TestEvictClient_RecordsMatchingPath(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"k": "v"},
	})
	c := NewEvictClient(mock, "secret/*")
	_, err := c.ReadSecrets("secret/foo")
	if err != nil {
		t.Fatal(err)
	}
	evicted := c.Evicted()
	if len(evicted) != 1 || evicted[0] != "secret/foo" {
		t.Fatalf("unexpected evictions: %v", evicted)
	}
}

func TestEvictClient_DoesNotRecordNonMatchingPath(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"other/bar": {"k": "v"},
	})
	c := NewEvictClient(mock, "secret/*")
	_, _ = c.ReadSecrets("other/bar")
	if len(c.Evicted()) != 0 {
		t.Fatalf("expected no evictions")
	}
}

func TestEvictClient_MultiplePatterns(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"a": "1"},
		"config/bar": {"b": "2"},
	})
	c := NewEvictClient(mock, "secret/*", "config/*")
	_, _ = c.ReadSecrets("secret/foo")
	_, _ = c.ReadSecrets("config/bar")
	evicted := c.Evicted()
	sort.Strings(evicted)
	if len(evicted) != 2 {
		t.Fatalf("expected 2 evictions, got %v", evicted)
	}
}

func TestEvictClient_Reset(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"k": "v"},
	})
	c := NewEvictClient(mock, "secret/*")
	_, _ = c.ReadSecrets("secret/foo")
	c.Reset()
	if len(c.Evicted()) != 0 {
		t.Fatal("expected empty after reset")
	}
}

func TestEvictClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/foo", errors.New("vault down"))
	c := NewEvictClient(mock, "secret/*")
	_, err := c.ReadSecrets("secret/foo")
	if err == nil {
		t.Fatal("expected error")
	}
	if len(c.Evicted()) != 0 {
		t.Fatal("should not record eviction on error")
	}
}

func TestEvictClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(nil)
	var _ SecretReader = NewEvictClient(mock, "*")
}
