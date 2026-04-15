package vault

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
)

// ErrQuotaExceeded is returned when the read quota for a client has been exhausted.
var ErrQuotaExceeded = errors.New("vault: read quota exceeded")

// QuotaConfig controls the behaviour of NewQuotaClient.
type QuotaConfig struct {
	// MaxReads is the maximum number of ReadSecrets calls allowed.
	// Zero means unlimited.
	MaxReads uint64
}

// DefaultQuotaConfig returns a QuotaConfig with no limit.
func DefaultQuotaConfig() QuotaConfig {
	return QuotaConfig{MaxReads: 0}
}

type quotaClient struct {
	inner    SecretReader
	max      uint64
	counter  atomic.Uint64
}

// NewQuotaClient wraps inner and returns ErrQuotaExceeded once more than
// cfg.MaxReads calls to ReadSecrets have been made. A MaxReads of zero
// disables the quota entirely.
func NewQuotaClient(inner SecretReader, cfg QuotaConfig) SecretReader {
	if inner == nil {
		panic("vault: NewQuotaClient: inner must not be nil")
	}
	return &quotaClient{inner: inner, max: cfg.MaxReads}
}

func (q *quotaClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	if q.max > 0 {
		n := q.counter.Add(1)
		if n > q.max {
			return nil, fmt.Errorf("%w: limit is %d", ErrQuotaExceeded, q.max)
		}
	}
	return q.inner.ReadSecrets(ctx, path)
}

// Remaining returns the number of reads still permitted.
// Returns 0 when no quota is configured.
func (q *quotaClient) Remaining() uint64 {
	if q.max == 0 {
		return 0
	}
	used := q.counter.Load()
	if used >= q.max {
		return 0
	}
	return q.max - used
}
