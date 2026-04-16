package vault

import (
	"fmt"
	"math/rand"
)

// SampleClient wraps a SecretReader and returns only a random subset of the
// secret keys on each read. It is useful for spot-checking large secret stores
// without pulling every key across the wire.
type SampleClient struct {
	inner SecretReader
	rate  float64 // 0.0–1.0: fraction of keys to keep
	rng   *rand.Rand
}

// NewSampleClient returns a SampleClient that forwards reads to inner and
// retains each key with probability rate (0.0 = drop all, 1.0 = keep all).
// A nil inner or a rate outside [0, 1] causes a panic.
func NewSampleClient(inner SecretReader, rate float64, rng *rand.Rand) *SampleClient {
	if inner == nil {
		panic("vault: NewSampleClient: inner must not be nil")
	}
	if rate < 0 || rate > 1 {
		panic(fmt.Sprintf("vault: NewSampleClient: rate %v out of range [0, 1]", rate))
	}
	if rng == nil {
		rng = rand.New(rand.NewSource(42)) //nolint:gosec
	}
	return &SampleClient{inner: inner, rate: rate, rng: rng}
}

// ReadSecrets delegates to the inner client and then randomly drops keys
// according to the configured sample rate.
func (s *SampleClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := s.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}

	// rate == 1 → keep everything, no allocation needed
	if s.rate == 1.0 {
		return secrets, nil
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if s.rng.Float64() < s.rate {
			out[k] = v
		}
	}
	return out, nil
}
