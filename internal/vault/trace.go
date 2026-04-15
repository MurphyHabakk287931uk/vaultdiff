package vault

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

// TraceClient wraps a SecretReader and writes a human-readable trace of every
// ReadSecrets call — path, duration, and whether it succeeded — to a Writer.
type TraceClient struct {
	inner  SecretReader
	out    io.Writer
	prefix string
}

// NewTraceClient returns a TraceClient that writes traces to out.
// If out is nil, os.Stderr is used. prefix is prepended to every line.
func NewTraceClient(inner SecretReader, out io.Writer, prefix string) *TraceClient {
	if inner == nil {
		panic("vault: NewTraceClient requires a non-nil inner client")
	}
	if out == nil {
		out = os.Stderr
	}
	return &TraceClient{inner: inner, out: out, prefix: prefix}
}

// ReadSecrets delegates to the inner client and traces the call.
func (t *TraceClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	start := time.Now()
	secrets, err := t.inner.ReadSecrets(ctx, path)
	dur := time.Since(start)

	status := "ok"
	if err != nil {
		status = fmt.Sprintf("error: %v", err)
	}

	line := fmt.Sprintf("%s[trace] path=%q duration=%s status=%s keys=%d\n",
		t.prefix, path, dur.Round(time.Microsecond), status, len(secrets))
	_, _ = fmt.Fprint(t.out, line)

	return secrets, err
}
