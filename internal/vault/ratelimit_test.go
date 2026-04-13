package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRateLimitClient_AllowsRequestWithinBurst(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "value"},
	})
	client := NewRateLimitClient(mock, RateLimitConfig{RequestsPerSecond: 100, Burst: 10})

	ctx := context.Background()
	secrets, err := client.ReadSecrets(ctx, "secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "value" {
		t.Errorf("expected 'value', got %q", secrets["key"])
	}
}

func TestRateLimitClient_PropagatesInnerError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/missing", errors.New("not found"))
	client := NewRateLimitClient(mock, DefaultRateLimitConfig())

	_, err := client.ReadSecrets(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRateLimitClient_CancelledContextReturnsError(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	// Very low rate so the limiter will need to wait
	client := NewRateLimitClient(mock, RateLimitConfig{RequestsPerSecond: 0.001, Burst: 1})

	// Drain the single burst token first
	_ , _ = client.ReadSecrets(context.Background(), "any")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	_, err := client.ReadSecrets(ctx, "any")
	if err == nil {
		t.Fatal("expected context deadline error, got nil")
	}
}

func TestRateLimitClient_ZeroRateReturnsInnerDirectly(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/b": {"x": "y"},
	})
	client := NewRateLimitClient(mock, RateLimitConfig{RequestsPerSecond: 0})

	// Should be the raw mock (no wrapping)
	if _, ok := client.(*RateLimitClient); ok {
		t.Error("expected unwrapped client when rate <= 0")
	}
	secrets, err := client.ReadSecrets(context.Background(), "secret/b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["x"] != "y" {
		t.Errorf("expected 'y', got %q", secrets["x"])
	}
}

func TestRateLimitClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(nil)
	var _ SecretReader = NewRateLimitClient(mock, DefaultRateLimitConfig())
}

func TestDefaultRateLimitConfig(t *testing.T) {
	cfg := DefaultRateLimitConfig()
	if cfg.RequestsPerSecond != 10 {
		t.Errorf("expected 10 RPS, got %v", cfg.RequestsPerSecond)
	}
	if cfg.Burst != 5 {
		t.Errorf("expected burst 5, got %d", cfg.Burst)
	}
}
