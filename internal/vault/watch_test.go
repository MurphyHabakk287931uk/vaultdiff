package vault

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestWatchClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewWatchClient(context.Background(), nil, "secret/a", DefaultWatchConfig())
}

func TestWatchClient_FiresOnChangeWhenSecretsChange(t *testing.T) {
	call := 0
	mock := NewMockClient(map[string]map[string]string{})
	mock.SetSecrets("secret/a", map[string]string{"k": "v1"})

	var mu sync.Mutex
	changed := make(chan map[string]string, 2)

	cfg := WatchConfig{
		Interval: 20 * time.Millisecond,
		OnChange: func(_ string, s map[string]string) {
			mu.Lock()
			defer mu.Unlock()
			call++
			changed <- s
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stop := NewWatchClient(ctx, mock, "secret/a", cfg)
	defer stop()

	// First tick fires OnChange (nil -> v1).
	select {
	case s := <-changed:
		if s["k"] != "v1" {
			t.Fatalf("expected v1, got %s", s["k"])
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for first change")
	}

	// Update the secret and wait for second callback.
	mock.SetSecrets("secret/a", map[string]string{"k": "v2"})
	select {
	case s := <-changed:
		if s["k"] != "v2" {
			t.Fatalf("expected v2, got %s", s["k"])
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for second change")
	}
	_ = call
}

func TestWatchClient_CallsOnError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := NewMockClient(map[string]map[string]string{})
	mock.SetError("secret/a", sentinel)

	errCh := make(chan error, 1)
	cfg := WatchConfig{
		Interval: 20 * time.Millisecond,
		OnError: func(_ string, err error) { errCh <- err },
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	stop := NewWatchClient(ctx, mock, "secret/a", cfg)
	defer stop()

	select {
	case err := <-errCh:
		if !errors.Is(err, sentinel) {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for error callback")
	}
}

func TestWatchClient_StopsOnCancel(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	mock.SetSecrets("secret/a", map[string]string{"x": "1"})

	count := 0
	cfg := WatchConfig{
		Interval: 20 * time.Millisecond,
		OnChange: func(_ string, _ map[string]string) { count++ },
	}

	ctx := context.Background()
	stop := NewWatchClient(ctx, mock, "secret/a", cfg)

	time.Sleep(80 * time.Millisecond)
	stop()
	snapshot := count
	time.Sleep(60 * time.Millisecond)

	if count != snapshot {
		t.Fatalf("poller still running after cancel: count went from %d to %d", snapshot, count)
	}
}

func TestMapsEqual(t *testing.T) {
	if !mapsEqual(map[string]string{"a": "1"}, map[string]string{"a": "1"}) {
		t.Fatal("expected equal")
	}
	if mapsEqual(map[string]string{"a": "1"}, map[string]string{"a": "2"}) {
		t.Fatal("expected not equal")
	}
	if mapsEqual(map[string]string{"a": "1"}, map[string]string{}) {
		t.Fatal("expected not equal for different lengths")
	}
}
