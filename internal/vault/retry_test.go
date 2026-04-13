package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

func fastRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Millisecond,
		MaxDelay:    5 * time.Millisecond,
	}
}

func TestRetryClient_SucceedsFirstAttempt(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"key": "value"},
	})
	client := NewRetryClient(mock, fastRetryConfig())

	result, err := client.ReadSecrets(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected value, got %q", result["key"])
	}
}

func TestRetryClient_RetriesTransientError(t *testing.T) {
	attempts := 0
	transientErr := &TransientError{Cause: errors.New("connection reset")}

	stub := &callbackReader{fn: func(_ context.Context, _ string) (map[string]string, error) {
		attempts++
		if attempts < 3 {
			return nil, transientErr
		}
		return map[string]string{"k": "v"}, nil
	}}

	client := NewRetryClient(stub, fastRetryConfig())
	result, err := client.ReadSecrets(context.Background(), "secret/x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
	if result["k"] != "v" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestRetryClient_DoesNotRetryPermanentError(t *testing.T) {
	attempts := 0
	permanent := errors.New("permission denied")

	stub := &callbackReader{fn: func(_ context.Context, _ string) (map[string]string, error) {
		attempts++
		return nil, permanent
	}}

	client := NewRetryClient(stub, fastRetryConfig())
	_, err := client.ReadSecrets(context.Background(), "secret/x")
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestRetryClient_ExhaustsAllAttempts(t *testing.T) {
	attempts := 0
	stub := &callbackReader{fn: func(_ context.Context, _ string) (map[string]string, error) {
		attempts++
		return nil, &TransientError{Cause: errors.New("timeout")}
	}}

	cfg := fastRetryConfig()
	client := NewRetryClient(stub, cfg)
	_, err := client.ReadSecrets(context.Background(), "secret/x")
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if attempts != cfg.MaxAttempts {
		t.Errorf("expected %d attempts, got %d", cfg.MaxAttempts, attempts)
	}
}

func TestRetryClient_RespectsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	stub := &callbackReader{fn: func(_ context.Context, _ string) (map[string]string, error) {
		return nil, &TransientError{Cause: errors.New("timeout")}
	}}

	client := NewRetryClient(stub, fastRetryConfig())
	_, err := client.ReadSecrets(ctx, "secret/x")
	if err == nil {
		t.Fatal("expected error due to cancelled context")
	}
}

// callbackReader is a test helper that delegates to a function.
type callbackReader struct {
	fn func(ctx context.Context, path string) (map[string]string, error)
}

func (c *callbackReader) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	return c.fn(ctx, path)
}
