package vault

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AuditEntry records a single secret-read operation.
type AuditEntry struct {
	Timestamp time.Time
	Path      string
	Keys      int
	Err       error
}

// AuditLogger writes structured audit entries to a writer.
type AuditLogger struct {
	out io.Writer
}

// NewAuditLogger returns an AuditLogger that writes to w.
// If w is nil, os.Stderr is used.
func NewAuditLogger(w io.Writer) *AuditLogger {
	if w == nil {
		w = os.Stderr
	}
	return &AuditLogger{out: w}
}

// Log writes a single audit entry in a human-readable format.
func (l *AuditLogger) Log(e AuditEntry) {
	status := "ok"
	if e.Err != nil {
		status = fmt.Sprintf("error: %v", e.Err)
	}
	fmt.Fprintf(l.out, "[audit] %s path=%s keys=%d status=%s\n",
		e.Timestamp.UTC().Format(time.RFC3339),
		e.Path,
		e.Keys,
		status,
	)
}

// AuditClient wraps a SecretReader and emits an audit log entry per read.
type AuditClient struct {
	inner  SecretReader
	logger *AuditLogger
}

// NewAuditClient returns an AuditClient that logs every ReadSecrets call.
func NewAuditClient(inner SecretReader, logger *AuditLogger) *AuditClient {
	if logger == nil {
		logger = NewAuditLogger(nil)
	}
	return &AuditClient{inner: inner, logger: logger}
}

// ReadSecrets delegates to the inner client and logs the outcome.
func (a *AuditClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := a.inner.ReadSecrets(path)
	a.logger.Log(AuditEntry{
		Timestamp: time.Now(),
		Path:      path,
		Keys:      len(secrets),
		Err:       err,
	})
	return secrets, err
}
