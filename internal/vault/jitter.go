package vault

import (
	"context"
	"math/rand"
	"time"
)

// JitterClient wraps a SecretReader and introduces a random delay before each
// read. This is useful for spreading load when many clients start simultaneously.
type JitterClient struct {
	inner  SecretReader
	minDur time.Duration
	maxDur time.Duration
	rng    *rand.Rand
}

// JitterConfig controls the bounds of the random delay.
type JitterConfig struct {
	Min time.Duration
	Max time.Duration
}

// DefaultJitterConfig returns a sensible default (0–100 ms).
func DefaultJitterConfig() JitterConfig {
	return JitterConfig{
		Min: 0,
		Max: 100 * time.Millisecond,
	}
}

// NewJitterClient returns a SecretReader that sleeps for a random duration in
// [cfg.Min, cfg.Max) before delegating to inner.
// Panics if inner is nil or cfg.Max < cfg.Min.
func NewJitterClient(inner SecretReader, cfg JitterConfig) SecretReader {
	if inner == nil {
		panic("jitter: inner client must not be nil")
	}
	if cfg.Max < cfg.Min {
		panic("jitter: max must be >= min")
	}
	return &JitterClient{
		inner:  inner,
		minDur: cfg.Min,
		maxDur: cfg.Max,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (j *JitterClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	spread := j.maxDur - j.minDur
	delay := j.minDur
	if spread > 0 {
		delay += time.Duration(j.rng.Int63n(int64(spread)))
	}
	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return j.inner.ReadSecrets(ctx, path)
}
