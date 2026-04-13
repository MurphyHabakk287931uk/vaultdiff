package vault

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitClient wraps a SecretReader and enforces a token-bucket rate limit
// on calls to ReadSecrets. This prevents overwhelming a Vault cluster when
// diffing many paths in rapid succession.
type RateLimitClient struct {
	inner   SecretReader
	limiter *rate.Limiter
}

// RateLimitConfig controls the token-bucket parameters.
type RateLimitConfig struct {
	// RequestsPerSecond is the sustained rate allowed.
	RequestsPerSecond float64
	// Burst is the maximum number of requests allowed in a single instant.
	Burst int
}

// DefaultRateLimitConfig returns a conservative default: 10 RPS, burst of 5.
func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		RequestsPerSecond: 10,
		Burst:             5,
	}
}

// NewRateLimitClient wraps inner with a rate limiter using the given config.
// If cfg.RequestsPerSecond <= 0 the client is returned unwrapped.
func NewRateLimitClient(inner SecretReader, cfg RateLimitConfig) SecretReader {
	if inner == nil {
		panic("vault: NewRateLimitClient: inner must not be nil")
	}
	if cfg.RequestsPerSecond <= 0 {
		return inner
	}
	if cfg.Burst <= 0 {
		cfg.Burst = 1
	}
	return &RateLimitClient{
		inner:   inner,
		limiter: rate.NewLimiter(rate.Limit(cfg.RequestsPerSecond), cfg.Burst),
	}
}

// ReadSecrets waits until the rate limiter grants a token, then delegates to
// the inner SecretReader. If the context is cancelled while waiting the error
// is propagated immediately.
func (r *RateLimitClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	if err := r.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}
	return r.inner.ReadSecrets(ctx, path)
}

// WaitDuration returns the estimated wait time before the next token is
// available. Useful for diagnostics and tests.
func (r *RateLimitClient) WaitDuration() time.Duration {
	return r.limiter.Reserve().Delay()
}
