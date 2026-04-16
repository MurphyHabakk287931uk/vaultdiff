package vault

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDedupeClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewDedupeClient(nil)
}

func TestDedupeClient_SingleRead(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	client := NewDedupeClient(mock)

	got, err := client.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("expected value, got %q", got["key"])
	}
}

func TestDedupeClient_PropagatesError(t *testing.T) {
	sentinel := errors.New("read failed")
	mock := NewMockClient(nil)
	mock.SetError("secret/bad", sentinel)
	client := NewDedupeClient(mock)

	_, err := client.ReadSecrets(context.Background(), "secret/bad")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestDedupeClient_CollapsesConcurrentReads(t *testing.T) {
	var callCount int64

	blocking := &blockingReader{
		ready: make(chan struct{}),
		onRead: func() { atomic.AddInt64(&callCount, 1) },
		result: map[string]string{"x": "1"},
	}
	client := NewDedupeClient(blocking)

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	// All goroutines fire before the blocking reader unblocks.
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			client.ReadSecrets(context.Background(), "secret/shared") //nolint:errcheck
		}()
	}

	// Give goroutines time to queue up, then unblock.
	time.Sleep(20 * time.Millisecond)
	close(blocking.ready)
	wg.Wait()

	if n := atomic.LoadInt64(&callCount); n != 1 {
		t.Errorf("expected 1 upstream call, got %d", n)
	}
}

func TestDedupeClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewDedupeClient(NewMockClient(nil))
}

// blockingReader is a SecretReader that blocks until ready is closed.
type blockingReader struct {
	ready  chan struct{}
	onRead func()
	result map[string]string
}

func (b *blockingReader) ReadSecrets(_ context.Context, _ string) (map[string]string, error) {
	b.onRead()
	<-b.ready
	return b.result, nil
}
