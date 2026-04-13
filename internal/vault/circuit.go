package vault

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrCircuitOpen is returned when the circuit breaker is in the open state.
var ErrCircuitOpen = errors.New("circuit breaker open: too many recent failures")

// circuitState represents the current state of the circuit breaker.
type circuitState int

const (
	stateClosed circuitState = iota
	stateOpen
	stateHalfOpen
)

// CircuitBreakerConfig holds configuration for the circuit breaker.
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of consecutive failures before opening.
	FailureThreshold int
	// SuccessThreshold is the number of successes in half-open to close again.
	SuccessThreshold int
	// OpenDuration is how long to stay open before moving to half-open.
	OpenDuration time.Duration
}

// DefaultCircuitBreakerConfig returns a sensible default configuration.
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		OpenDuration:     30 * time.Second,
	}
}

// circuitBreakerClient wraps a SecretReader with circuit breaker logic.
type circuitBreakerClient struct {
	inner   SecretReader
	cfg     CircuitBreakerConfig
	mu      sync.Mutex
	state   circuitState
	failures int
	successes int
	openedAt time.Time
}

// NewCircuitBreakerClient wraps inner with a circuit breaker using cfg.
func NewCircuitBreakerClient(inner SecretReader, cfg CircuitBreakerConfig) SecretReader {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = DefaultCircuitBreakerConfig().FailureThreshold
	}
	if cfg.OpenDuration <= 0 {
		cfg.OpenDuration = DefaultCircuitBreakerConfig().OpenDuration
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = DefaultCircuitBreakerConfig().SuccessThreshold
	}
	return &circuitBreakerClient{inner: inner, cfg: cfg}
}

func (c *circuitBreakerClient) ReadSecrets(path string) (map[string]string, error) {
	c.mu.Lock()
	switch c.state {
	case stateOpen:
		if time.Since(c.openedAt) >= c.cfg.OpenDuration {
			c.state = stateHalfOpen
			c.successes = 0
		} else {
			c.mu.Unlock()
			return nil, fmt.Errorf("%w (retry after %s)", ErrCircuitOpen,
				c.cfg.OpenDuration-time.Since(c.openedAt).Round(time.Second))
		}
	case stateHalfOpen:
		// allow one probe through
	}
	c.mu.Unlock()

	secrets, err := c.inner.ReadSecrets(path)

	c.mu.Lock()
	defer c.mu.Unlock()

	if err != nil {
		c.failures++
		c.successes = 0
		if c.state == stateHalfOpen || c.failures >= c.cfg.FailureThreshold {
			c.state = stateOpen
			c.openedAt = time.Now()
		}
		return nil, err
	}

	if c.state == stateHalfOpen {
		c.successes++
		if c.successes >= c.cfg.SuccessThreshold {
			c.state = stateClosed
			c.failures = 0
		}
	} else {
		c.failures = 0
	}
	return secrets, nil
}
