package vault

import (
	"errors"
	"testing"
	"time"
)

// fakeMetrics records every RecordRead call for assertions.
type fakeMetrics struct {
	calls []metricCall
}

type metricCall struct {
	path        string
	durationMs  int64
	err         error
}

func (f *fakeMetrics) RecordRead(path string, durationMs int64, err error) {
	f.calls = append(f.calls, metricCall{path: path, durationMs: durationMs, err: err})
}

func TestObserveClient_RecordsSuccessfulRead(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	fm := &fakeMetrics{}
	client := NewObserveClient(mock, fm)

	secrets, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["key"] != "value" {
		t.Errorf("expected value, got %q", secrets["key"])
	}
	if len(fm.calls) != 1 {
		t.Fatalf("expected 1 metric call, got %d", len(fm.calls))
	}
	if fm.calls[0].path != "secret/app" {
		t.Errorf("expected path secret/app, got %q", fm.calls[0].path)
	}
	if fm.calls[0].err != nil {
		t.Errorf("expected nil error in metric, got %v", fm.calls[0].err)
	}
}

func TestObserveClient_RecordsError(t *testing.T) {
	sentinel := errors.New("vault unavailable")
	mock := NewMockClient(nil)
	mock.SetError("secret/fail", sentinel)
	fm := &fakeMetrics{}
	client := NewObserveClient(mock, fm)

	_, err := client.ReadSecrets("secret/fail")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
	if len(fm.calls) != 1 {
		t.Fatalf("expected 1 metric call, got %d", len(fm.calls))
	}
	if !errors.Is(fm.calls[0].err, sentinel) {
		t.Errorf("expected sentinel in metric, got %v", fm.calls[0].err)
	}
}

func TestObserveClient_RecordsDuration(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"k": "v"},
	})
	fm := &fakeMetrics{}
	client := NewObserveClient(mock, fm)

	// Advance the clock by 50 ms via a custom now function.
	calls := 0
	base := time.Now()
	client.now = func() time.Time {
		calls++
		if calls == 1 {
			return base
		}
		return base.Add(50 * time.Millisecond)
	}

	_, _ = client.ReadSecrets("secret/app")
	if fm.calls[0].durationMs != 50 {
		t.Errorf("expected 50 ms, got %d", fm.calls[0].durationMs)
	}
}

func TestObserveClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil inner")
		}
	}()
	NewObserveClient(nil, &fakeMetrics{})
}

func TestObserveClient_PanicsOnNilMetrics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil metrics")
		}
	}()
	mock := NewMockClient(nil)
	NewObserveClient(mock, nil)
}

func TestObserveClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = &ObserveClient{}
}
