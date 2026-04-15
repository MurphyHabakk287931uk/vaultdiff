package vault

import (
	"context"
	"time"
)

// WatchConfig controls polling behaviour for the watch client.
type WatchConfig struct {
	// Interval is how often the secret path is re-read.
	Interval time.Duration
	// OnChange is called whenever the secrets at the path change.
	OnChange func(path string, secrets map[string]string)
	// OnError is called when a read fails; if nil errors are silently dropped.
	OnError func(path string, err error)
}

// DefaultWatchConfig returns a WatchConfig with sensible defaults.
func DefaultWatchConfig() WatchConfig {
	return WatchConfig{
		Interval: 30 * time.Second,
	}
}

// watchClient polls a SecretReader on a fixed interval and fires callbacks
// when the secret map changes.
type watchClient struct {
	inner  SecretReader
	cfg    WatchConfig
}

// NewWatchClient wraps inner and starts a background goroutine that polls
// path every cfg.Interval until ctx is cancelled.
// The returned CancelFunc stops the poller immediately.
func NewWatchClient(ctx context.Context, inner SecretReader, path string, cfg WatchConfig) context.CancelFunc {
	if inner == nil {
		panic("vault: NewWatchClient: inner must not be nil")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = DefaultWatchConfig().Interval
	}

	watchCtx, cancel := context.WithCancel(ctx)

	go func() {
		var last map[string]string
		ticker := time.NewTicker(cfg.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-watchCtx.Done():
				return
			case <-ticker.C:
				secrets, err := inner.ReadSecrets(watchCtx, path)
				if err != nil {
					if cfg.OnError != nil {
						cfg.OnError(path, err)
					}
					continue
				}
				if cfg.OnChange != nil && !mapsEqual(last, secrets) {
					cfg.OnChange(path, secrets)
				}
				last = secrets
			}
		}
	}()

	return cancel
}

func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
