package vault

import (
	"errors"
	"testing"
	"time"
)

var fastCircuitConfig = CircuitBreakerConfig{
	FailureThreshold: 3,
	SuccessThreshold: 2,
	OpenDuration:     50 * time.Millisecond,
}

func TestCircuitBreaker_SucceedsNormally(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"k": "v"},
	})
	cb := NewCircuitBreakerClient(mock, fastCircuitConfig)

	got, err := cb.ReadSecrets("secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["k"] != "v" {
		t.Errorf("expected v, got %q", got["k"])
	}
}

func TestCircuitBreaker_OpensAfterThreshold(t *testing.T) {
	errBoom := errors.New("vault unavailable")
	mock := NewMockClient(nil)
	mock.SetError("secret/a", errBoom)
	cb := NewCircuitBreakerClient(mock, fastCircuitConfig)

	for i := 0; i < fastCircuitConfig.FailureThreshold; i++ {
		_, _ = cb.ReadSecrets("secret/a")
	}

	_, err := cb.ReadSecrets("secret/a")
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenAfterDuration(t *testing.T) {
	errBoom := errors.New("vault unavailable")
	mock := NewMockClient(nil)
	mock.SetError("secret/a", errBoom)
	cb := NewCircuitBreakerClient(mock, fastCircuitConfig)

	for i := 0; i < fastCircuitConfig.FailureThreshold; i++ {
		_, _ = cb.ReadSecrets("secret/a")
	}

	time.Sleep(fastCircuitConfig.OpenDuration + 10*time.Millisecond)

	// Should probe (half-open), fail, and re-open
	_, err := cb.ReadSecrets("secret/a")
	if errors.Is(err, ErrCircuitOpen) {
		t.Fatal("expected probe attempt, not circuit open error")
	}
	if !errors.Is(err, errBoom) {
		t.Fatalf("expected underlying error, got %v", err)
	}
}

func TestCircuitBreaker_ClosesAfterSuccessThreshold(t *testing.T) {
	errBoom := errors.New("vault unavailable")
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"x": "1"},
	})
	mock.SetError("secret/a", errBoom)
	cb := NewCircuitBreakerClient(mock, fastCircuitConfig)

	for i := 0; i < fastCircuitConfig.FailureThreshold; i++ {
		_, _ = cb.ReadSecrets("secret/a")
	}

	time.Sleep(fastCircuitConfig.OpenDuration + 10*time.Millisecond)
	mock.ClearError("secret/a")

	for i := 0; i < fastCircuitConfig.SuccessThreshold; i++ {
		_, err := cb.ReadSecrets("secret/a")
		if err != nil {
			t.Fatalf("probe %d failed: %v", i, err)
		}
	}

	// Circuit should now be closed — normal call should succeed
	_, err := cb.ReadSecrets("secret/a")
	if err != nil {
		t.Fatalf("expected closed circuit, got: %v", err)
	}
}

func TestCircuitBreaker_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewCircuitBreakerClient(NewMockClient(nil), DefaultCircuitBreakerConfig())
}
