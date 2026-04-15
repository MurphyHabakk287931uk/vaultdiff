package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time { return func() time.Time { return t } }

func TestTTLClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewTTLClient(nil, time.Minute)
}

func TestTTLClient_FirstReadSucceeds(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"key": "value"},
	})
	client := NewTTLClient(mock, time.Minute)

	got, err := client.ReadSecrets(context.Background(), "secret/foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("expected 'value', got %q", got["key"])
	}
}

func TestTTLClient_StaleReturnsError(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"key": "value"},
	})
	now := time.Now()
	client := NewTTLClient(mock, time.Minute)
	client.now = fixedClock(now)

	// Prime the fetch time.
	if _, err := client.ReadSecrets(context.Background(), "secret/foo"); err != nil {
		t.Fatalf("unexpected error on first read: %v", err)
	}

	// Advance clock beyond TTL.
	client.now = fixedClock(now.Add(2 * time.Minute))

	_, err := client.ReadSecrets(context.Background(), "secret/foo")
	if err == nil {
		t.Fatal("expected stale error, got nil")
	}
}

func TestTTLClient_ZeroMaxAgeDisablesTTL(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"k": "v"},
	})
	client := NewTTLClient(mock, 0)
	client.now = fixedClock(time.Now().Add(24 * time.Hour))

	// Even after a huge time gap, no stale error when maxAge == 0.
	if _, err := client.ReadSecrets(context.Background(), "secret/foo"); err != nil {
		t.Fatalf("unexpected error with zero maxAge: %v", err)
	}
}

func TestTTLClient_InvalidateAllowsRefresh(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/foo": {"k": "v"},
	})
	now := time.Now()
	client := NewTTLClient(mock, time.Minute)
	client.now = fixedClock(now)

	if _, err := client.ReadSecrets(context.Background(), "secret/foo"); err != nil {
		t.Fatal(err)
	}
	client.now = fixedClock(now.Add(2 * time.Minute))
	client.Invalidate("secret/foo")

	// After invalidation the stale check is skipped.
	if _, err := client.ReadSecrets(context.Background(), "secret/foo"); err != nil {
		t.Fatalf("expected success after invalidation, got: %v", err)
	}
}

func TestTTLClient_PropagatesInnerError(t *testing.T) {
	sentinel := errors.New("inner failure")
	mock := NewMockClient(nil)
	mock.SetError("secret/foo", sentinel)
	client := NewTTLClient(mock, time.Minute)

	_, err := client.ReadSecrets(context.Background(), "secret/foo")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got: %v", err)
	}
}

func TestTTLClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewTTLClient(NewMockClient(nil), time.Minute)
}
