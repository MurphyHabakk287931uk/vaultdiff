package vault

import (
	"time"
)

// Metrics is the interface for recording vault client observations.
type Metrics interface {
	// RecordRead records a completed ReadSecrets call.
	RecordRead(path string, durationMs int64, err error)
}

// ObserveClient wraps a SecretReader and records metrics for every read.
type ObserveClient struct {
	inner   SecretReader
	metrics Metrics
	now     func() time.Time
}

// NewObserveClient returns an ObserveClient that delegates to inner and
// records observations via metrics. Panics if either argument is nil.
func NewObserveClient(inner SecretReader, metrics Metrics) *ObserveClient {
	if inner == nil {
		panic("observe: inner SecretReader must not be nil")
	}
	if metrics == nil {
		panic("observe: Metrics must not be nil")
	}
	return &ObserveClient{
		inner:   inner,
		metrics: metrics,
		now:     time.Now,
	}
}

// ReadSecrets delegates to the inner client and records the call duration
// and any error via the Metrics implementation.
func (o *ObserveClient) ReadSecrets(path string) (map[string]string, error) {
	start := o.now()
	secrets, err := o.inner.ReadSecrets(path)
	duration := o.now().Sub(start).Milliseconds()
	o.metrics.RecordRead(path, duration, err)
	return secrets, err
}
