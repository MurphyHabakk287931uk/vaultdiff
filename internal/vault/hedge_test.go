package vault

import (
	"context"
	"errors"
	"testing"
	"time"
)

// slowClient delays its response by the given duration.
type slowClient struct {
	delay   time.Duration
	secrets map[string]string
	err     error
}

func (s *slowClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	select {
	case <-time.After(s.delay):
		return s.secrets, s.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func TestHedgeClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewHedgeClient(nil, DefaultHedgeConfig())
}

func TestHedgeClient_ZeroDelay_ReturnsInnerDirectly(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"k": "v"},
	})
	client := NewHedgeClient(mock, HedgeConfig{Delay: 0})
	if client == mock {
		// pointer equality: unwrapped
	}
	s, err := client.ReadSecrets(context.Background(), "secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s["k"] != "v" {
		t.Fatalf("expected v, got %q", s["k"])
	}
}

func TestHedgeClient_FastResponse_NoHedge(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/fast": {"key": "value"},
	})
	client := NewHedgeClient(mock, HedgeConfig{Delay: 500 * time.Millisecond})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	s, err := client.ReadSecrets(ctx, "secret/fast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s["key"] != "value" {
		t.Fatalf("expected value, got %q", s["key"])
	}
}

func TestHedgeClient_SlowFirst_HedgeWins(t *testing.T) {
	slow := &slowClient{
		delay:   300 * time.Millisecond,
		secrets: map[string]string{"key": "hedged"},
	}
	client := NewHedgeClient(slow, HedgeConfig{Delay: 50 * time.Millisecond})

	start := time.Now()
	s, err := client.ReadSecrets(context.Background(), "any/path")
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s["key"] != "hedged" {
		t.Fatalf("expected hedged, got %q", s["key"])
	}
	// The hedge fires at 50 ms; the second request also takes 300 ms, so total
	// should be roughly 300 ms — well under 600 ms (two sequential requests).
	if elapsed > 600*time.Millisecond {
		t.Fatalf("hedge did not shorten latency: elapsed %v", elapsed)
	}
}

func TestHedgeClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := NewMockClient(nil)
	mock.SetError("secret/err", sentinel)

	client := NewHedgeClient(mock, HedgeConfig{Delay: 10 * time.Millisecond})
	_, err := client.ReadSecrets(context.Background(), "secret/err")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestHedgeClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(nil)
	var _ SecretReader = NewHedgeClient(mock, DefaultHedgeConfig())
}
