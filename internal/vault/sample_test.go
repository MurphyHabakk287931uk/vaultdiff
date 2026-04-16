package vault

import (
	"errors"
	"math/rand"
	"testing"
)

func TestSampleClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewSampleClient(nil, 0.5, nil)
}

func TestSampleClient_PanicsOnNegativeRate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative rate")
		}
	}()
	NewSampleClient(NewMockClient(nil, nil), -0.1, nil)
}

func TestSampleClient_PanicsOnRateAboveOne(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for rate > 1")
		}
	}()
	NewSampleClient(NewMockClient(nil, nil), 1.1, nil)
}

func TestSampleClient_RateOne_KeepsAll(t *testing.T) {
	secrets := map[string]string{"a": "1", "b": "2", "c": "3"}
	mock := NewMockClient(secrets, nil)
	client := NewSampleClient(mock, 1.0, rand.New(rand.NewSource(0)))

	got, err := client.ReadSecrets("secret/data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(secrets) {
		t.Fatalf("expected %d keys, got %d", len(secrets), len(got))
	}
}

func TestSampleClient_RateZero_DropsAll(t *testing.T) {
	secrets := map[string]string{"a": "1", "b": "2", "c": "3"}
	mock := NewMockClient(secrets, nil)
	// Use a fixed RNG; with rate 0 no key should pass the threshold.
	client := NewSampleClient(mock, 0.0, rand.New(rand.NewSource(0)))

	got, err := client.ReadSecrets("secret/data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(got))
	}
}

func TestSampleClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := NewMockClient(nil, sentinel)
	client := NewSampleClient(mock, 0.5, nil)

	_, err := client.ReadSecrets("secret/data")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestSampleClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewSampleClient(NewMockClient(nil, nil), 0.5, nil)
}

func TestSampleClient_NilRng_UsesDefault(t *testing.T) {
	secrets := map[string]string{"key": "value"}
	mock := NewMockClient(secrets, nil)
	// Should not panic when rng is nil.
	client := NewSampleClient(mock, 1.0, nil)

	got, err := client.ReadSecrets("any/path")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Fatalf("expected key=value, got %v", got)
	}
}
