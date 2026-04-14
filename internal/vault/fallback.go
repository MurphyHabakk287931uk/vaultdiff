package vault

import (
	"context"
	"errors"
)

// FallbackClient tries a primary SecretReader and, if it returns an error
// matching the supplied predicate, transparently falls back to a secondary.
// This is useful for cross-environment promotion checks or graceful degradation
// when a path does not exist in one Vault cluster.
type FallbackClient struct {
	primary   SecretReader
	secondary SecretReader
	shouldFallback func(error) bool
}

// NewFallbackClient constructs a FallbackClient.
// If shouldFallback is nil, any non-nil error from primary triggers the fallback.
func NewFallbackClient(primary, secondary SecretReader, shouldFallback func(error) bool) *FallbackClient {
	if primary == nil {
		panic("vault: FallbackClient primary must not be nil")
	}
	if secondary == nil {
		panic("vault: FallbackClient secondary must not be nil")
	}
	if shouldFallback == nil {
		shouldFallback = func(err error) bool { return err != nil }
	}
	return &FallbackClient{
		primary:        primary,
		secondary:      secondary,
		shouldFallback: shouldFallback,
	}
}

// IsNotFound is a convenience predicate that triggers fallback only when the
// error wraps ErrNotFound, leaving other errors (network, permission) to
// propagate from the primary.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// ReadSecrets calls primary first. If shouldFallback returns true for the
// resulting error, secondary is tried instead and its result returned.
func (f *FallbackClient) ReadSecrets(ctx context.Context, path string) (map[string]string, error) {
	result, err := f.primary.ReadSecrets(ctx, path)
	if err != nil && f.shouldFallback(err) {
		return f.secondary.ReadSecrets(ctx, path)
	}
	return result, err
}
