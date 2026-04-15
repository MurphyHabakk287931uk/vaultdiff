package vault

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

func TestTraceClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil inner")
		}
	}()
	NewTraceClient(nil, nil, "")
}

func TestTraceClient_NilWriterUsesStderr(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"k": "v"},
	})
	// Should not panic even with nil writer.
	client := NewTraceClient(mock, nil, "")
	if client.out == nil {
		t.Fatal("expected out to default to stderr, got nil")
	}
}

func TestTraceClient_WritesTraceOnSuccess(t *testing.T) {
	var buf bytes.Buffer
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"user": "admin", "pass": "secret"},
	})
	client := NewTraceClient(mock, &buf, "pfx:")

	_, err := client.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	if !strings.HasPrefix(line, "pfx:") {
		t.Errorf("expected prefix 'pfx:' in output, got: %s", line)
	}
	if !strings.Contains(line, `path="secret/app"`) {
		t.Errorf("expected path in trace output, got: %s", line)
	}
	if !strings.Contains(line, "status=ok") {
		t.Errorf("expected status=ok in trace output, got: %s", line)
	}
	if !strings.Contains(line, "keys=2") {
		t.Errorf("expected keys=2 in trace output, got: %s", line)
	}
}

func TestTraceClient_WritesTraceOnError(t *testing.T) {
	var buf bytes.Buffer
	mock := NewMockClient(nil)
	mock.SetError("secret/missing", errors.New("not found"))
	client := NewTraceClient(mock, &buf, "")

	_, err := client.ReadSecrets(context.Background(), "secret/missing")
	if err == nil {
		t.Fatal("expected error")
	}

	line := buf.String()
	if !strings.Contains(line, "error:") {
		t.Errorf("expected error in trace output, got: %s", line)
	}
	if !strings.Contains(line, "keys=0") {
		t.Errorf("expected keys=0 in trace output, got: %s", line)
	}
}

func TestTraceClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewTraceClient(NewMockClient(nil), nil, "")
}

func TestTraceClient_EmptyPrefix(t *testing.T) {
	var buf bytes.Buffer
	mock := NewMockClient(map[string]map[string]string{
		"a": {"x": "1"},
	})
	client := NewTraceClient(mock, &buf, "")
	_, _ = client.ReadSecrets(context.Background(), "a")
	if !strings.HasPrefix(buf.String(), "[trace]") {
		t.Errorf("expected line to start with '[trace]', got: %s", buf.String())
	}
}
