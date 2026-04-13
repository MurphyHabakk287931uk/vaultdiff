package vault

import (
	"context"
	"fmt"
	"time"
)

// TimeoutClient wraps a SecretReader and enforces a per-call deadline.
type TimeoutClient struct {
	inner   SecretReader
	timeout time.Duration
}

// NewTimeoutClient returns a SecretReader that cancels any ReadSecrets call
// that exceeds the given timeout duration. A zero or negative timeout disables
// enforcement and delegates directly to the inner client.
func NewTimeoutClient(inner SecretReader, timeout time.Duration) *TimeoutClient {
	return &TimeoutClient{inner: inner, timeout: timeout}
}

// ReadSecrets reads secrets from the inner client, cancelling the call if the
// configured timeout elapses before a response is received.
func (c *TimeoutClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	if c.timeout <= 0 {
		return c.inner.ReadSecrets(ctx, path)
	}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	type result struct {
		data map[string]string
		err  error
	}

	ch := make(chan result, 1)
	go func() {
		data, err := c.inner.ReadSecrets(ctx, path)
		ch <- result{data, err}
	}()

	select {
	case res := <-ch:
		return res.data, res.err
	case <-ctx.Done():
		return nil, fmt.Errorf("vault: ReadSecrets timed out after %s for path %q", c.timeout, path)
	}
}
