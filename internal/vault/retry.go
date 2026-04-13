package vault

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// RetryConfig holds configuration for retry behaviour.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// DefaultRetryConfig returns sensible retry defaults.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    2 * time.Second,
	}
}

// retryClient wraps a SecretReader and retries transient errors.
type retryClient struct {
	inner  SecretReader
	config RetryConfig
}

// NewRetryClient returns a SecretReader that retries on transient failures.
func NewRetryClient(inner SecretReader, cfg RetryConfig) SecretReader {
	return &retryClient{inner: inner, config: cfg}
}

// ReadSecrets attempts to read secrets, retrying on transient errors.
func (r *retryClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	var lastErr error
	delay := r.config.BaseDelay

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		result, err := r.inner.ReadSecrets(ctx, path)
		if err == nil {
			return result, nil
		}

		if !isTransient(err) {
			return nil, err
		}

		lastErr = err
		if attempt < r.config.MaxAttempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
			delay = min(delay*2, r.config.MaxDelay)
		}
	}

	return nil, fmt.Errorf("all %d attempts failed: %w", r.config.MaxAttempts, lastErr)
}

// TransientError marks an error as retryable.
type TransientError struct {
	Cause error
}

func (e *TransientError) Error() string { return "transient: " + e.Cause.Error() }
func (e *TransientError) Unwrap() error { return e.Cause }

func isTransient(err error) bool {
	var t *TransientError
	return errors.As(err, &t)
}

func min(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
