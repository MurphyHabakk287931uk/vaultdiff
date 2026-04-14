package vault

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestCounterMetrics_InitialState(t *testing.T) {
	m := NewCounterMetrics()
	if m.TotalCalls() != 0 {
		t.Errorf("expected 0 calls, got %d", m.TotalCalls())
	}
	if m.TotalErrors() != 0 {
		t.Errorf("expected 0 errors, got %d", m.TotalErrors())
	}
}

func TestCounterMetrics_CountsSuccessfulReads(t *testing.T) {
	m := NewCounterMetrics()
	m.RecordRead("secret/a", 10, nil)
	m.RecordRead("secret/b", 20, nil)
	m.RecordRead("secret/a", 5, nil)

	if m.TotalCalls() != 3 {
		t.Errorf("expected 3 total calls, got %d", m.TotalCalls())
	}
	if m.TotalErrors() != 0 {
		t.Errorf("expected 0 errors, got %d", m.TotalErrors())
	}
}

func TestCounterMetrics_CountsErrors(t *testing.T) {
	m := NewCounterMetrics()
	m.RecordRead("secret/a", 10, nil)
	m.RecordRead("secret/a", 5, errors.New("boom"))
	m.RecordRead("secret/b", 8, errors.New("boom"))

	if m.TotalCalls() != 3 {
		t.Errorf("expected 3 calls, got %d", m.TotalCalls())
	}
	if m.TotalErrors() != 2 {
		t.Errorf("expected 2 errors, got %d", m.TotalErrors())
	}
}

func TestCounterMetrics_WriteSummary(t *testing.T) {
	m := NewCounterMetrics()
	m.RecordRead("secret/alpha", 12, nil)
	m.RecordRead("secret/beta", 8, errors.New("fail"))

	var buf bytes.Buffer
	m.WriteSummary(&buf)
	out := buf.String()

	if !strings.Contains(out, "secret/alpha") {
		t.Errorf("summary missing secret/alpha: %s", out)
	}
	if !strings.Contains(out, "secret/beta") {
		t.Errorf("summary missing secret/beta: %s", out)
	}
	if !strings.Contains(out, "total_ms=20") {
		t.Errorf("summary missing total_ms: %s", out)
	}
}

func TestCounterMetrics_WriteSummary_NilWriter(t *testing.T) {
	// Should not panic when w is nil (falls back to stderr).
	m := NewCounterMetrics()
	m.RecordRead("secret/x", 1, nil)
	m.WriteSummary(nil) // must not panic
}

func TestCounterMetrics_ImplementsMetrics(t *testing.T) {
	var _ Metrics = &CounterMetrics{}
}
