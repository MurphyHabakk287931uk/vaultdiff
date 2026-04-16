package vault

import (
	"context"
	"testing"
	"time"
)

func TestJitterClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewJitterClient(nil, DefaultJitterConfig())
}

func TestJitterClient_PanicsWhenMaxLessThanMin(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when max < min")
		}
	}()
	mock := NewMockClient(map[string]map[string]string{})
	NewJitterClient(mock, JitterConfig{Min: 50 * time.Millisecond, Max: 10 * time.Millisecond})
}

func TestJitterClient_ZeroRange_NoDelay(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/data": {"key": "val"},
	})
	client := NewJitterClient(mock, JitterConfig{Min: 0, Max: 0})
	start := time.Now()
	secrets, err := client.ReadSecrets(context.Background(), "secret/data")
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "val" {
		t.Fatalf("expected val, got %s", secrets["key"])
	}
	if elapsed > 20*time.Millisecond {
		t.Fatalf("expected near-zero delay, got %v", elapsed)
	}
}

func TestJitterClient_AppliesDelay(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/cfg": {"a": "b"},
	})
	cfg := JitterConfig{Min: 20 * time.Millisecond, Max: 40 * time.Millisecond}
	client := NewJitterClient(mock, cfg)
	start := time.Now()
	_, err := client.ReadSecrets(context.Background(), "secret/cfg")
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed < cfg.Min {
		t.Fatalf("expected at least %v delay, got %v", cfg.Min, elapsed)
	}
}

func TestJitterClient_CancelledContext(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	cfg := JitterConfig{Min: 500 * time.Millisecond, Max: time.Second}
	client := NewJitterClient(mock, cfg)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := client.ReadSecrets(ctx, "secret/x")
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestJitterClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	var _ SecretReader = NewJitterClient(mock, DefaultJitterConfig())
}
