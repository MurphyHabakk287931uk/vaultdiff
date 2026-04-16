package vault

import (
	"errors"
	"testing"
)

func TestCachedClient_ReturnsCachedResult(t *testing.T) {
	callCount := 0
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	_ = callCount

	cached := NewCachedClient(mock)

	first, err := cached.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := cached.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if first["key"] != second["key"] {
		t.Errorf("expected same result, got %v and %v", first, second)
	}
	if cached.CacheSize() != 1 {
		t.Errorf("expected cache size 1, got %d", cached.CacheSize())
	}
	// Verify the underlying client was only called once
	if mock.CallCount("secret/app") != 1 {
		t.Errorf("expected underlying client called once, got %d", mock.CallCount("secret/app"))
	}
}

func TestCachedClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/missing", errors.New("not found"))

	cached := NewCachedClient(mock)
	_, err := cached.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if cached.CacheSize() != 0 {
		t.Errorf("expected cache size 0 after error, got %d", cached.CacheSize())
	}
}

func TestCachedClient_Invalidate(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	cached := NewCachedClient(mock)

	_, _ = cached.ReadSecrets("secret/app")
	if cached.CacheSize() != 1 {
		t.Fatalf("expected cache size 1 before invalidate")
	}

	cached.Invalidate("secret/app")
	if cached.CacheSize() != 0 {
		t.Errorf("expected cache size 0 after invalidate, got %d", cached.CacheSize())
	}
}

func TestCachedClient_InvalidateAll(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"k": "v"},
		"secret/b": {"k": "v"},
	})
	cached := NewCachedClient(mock)

	_, _ = cached.ReadSecrets("secret/a")
	_, _ = cached.ReadSecrets("secret/b")
	if cached.CacheSize() != 2 {
		t.Fatalf("expected cache size 2 before clear")
	}

	cached.InvalidateAll()
	if cached.CacheSize() != 0 {
		t.Errorf("expected cache size 0 after InvalidateAll, got %d", cached.CacheSize())
	}
}

func TestCachedClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(nil)
	var _ SecretReader = NewCachedClient(mock)
}
