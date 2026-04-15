package vault

import (
	"context"
	"time"
)

// HedgeConfig controls hedged-request behaviour.
type HedgeConfig struct {
	// Delay is how long to wait before issuing the second request.
	// Zero disables hedging and delegates directly to the inner client.
	Delay time.Duration
}

// DefaultHedgeConfig returns a HedgeConfig with a 200 ms hedge delay.
func DefaultHedgeConfig() HedgeConfig {
	return HedgeConfig{Delay: 200 * time.Millisecond}
}

type hedgeClient struct {
	inner  SecretReader
	config HedgeConfig
}

// NewHedgeClient wraps inner so that, if a ReadSecrets call does not return
// within config.Delay, a second identical request is issued concurrently.
// Whichever request completes first (successfully or not) wins; the slower
// goroutine is abandoned when the returned context is cancelled by the caller.
//
// Panics if inner is nil.
func NewHedgeClient(inner SecretReader, config HedgeConfig) SecretReader {
	if inner == nil {
		panic("hedge: inner client must not be nil")
	}
	if config.Delay <= 0 {
		return inner
	}
	return &hedgeClient{inner: inner, config: config}
}

func (h *hedgeClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	type result struct {
		secrets map[string]string
		err     error
	}

	ch := make(chan result, 2)

	issue := func() {
		s, err := h.inner.ReadSecrets(ctx, path)
		ch <- result{s, err}
	}

	go issue()

	select {
	case res := <-ch:
		return res.secrets, res.err
	case <-time.After(h.config.Delay):
		// Hedge: fire a second request.
		go issue()
		res := <-ch
		return res.secrets, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
