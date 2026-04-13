package vault

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestAuditLogger_LogSuccess(t *testing.T) {
	var buf bytes.Buffer
	l := NewAuditLogger(&buf)
	l.Log(AuditEntry{Path: "secret/app", Keys: 3})

	got := buf.String()
	if !strings.Contains(got, "path=secret/app") {
		t.Errorf("expected path in log, got: %s", got)
	}
	if !strings.Contains(got, "keys=3") {
		t.Errorf("expected keys=3 in log, got: %s", got)
	}
	if !strings.Contains(got, "status=ok") {
		t.Errorf("expected status=ok in log, got: %s", got)
	}
}

func TestAuditLogger_LogError(t *testing.T) {
	var buf bytes.Buffer
	l := NewAuditLogger(&buf)
	l.Log(AuditEntry{Path: "secret/missing", Err: errors.New("not found")})

	got := buf.String()
	if !strings.Contains(got, "error: not found") {
		t.Errorf("expected error in log, got: %s", got)
	}
}

func TestNewAuditLogger_NilUsesStderr(t *testing.T) {
	l := NewAuditLogger(nil)
	if l.out == nil {
		t.Fatal("expected non-nil writer")
	}
}

func TestAuditClient_LogsOnSuccess(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	var buf bytes.Buffer
	client := NewAuditClient(mock, NewAuditLogger(&buf))

	secrets, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(secrets))
	}
	if !strings.Contains(buf.String(), "keys=1") {
		t.Errorf("expected keys=1 in audit log, got: %s", buf.String())
	}
}

func TestAuditClient_LogsOnError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/fail", errors.New("permission denied"))
	var buf bytes.Buffer
	client := NewAuditClient(mock, NewAuditLogger(&buf))

	_, err := client.ReadSecrets("secret/fail")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(buf.String(), "error: permission denied") {
		t.Errorf("expected error in audit log, got: %s", buf.String())
	}
}

func TestAuditClient_NilLogger_UsesDefault(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/x": {"a": "b"},
	})
	client := NewAuditClient(mock, nil)
	if client.logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestAuditClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = (*AuditClient)(nil)
}
