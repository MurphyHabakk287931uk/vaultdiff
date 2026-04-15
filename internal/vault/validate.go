package vault

import (
	"errors"
	"fmt"
	"strings"
)

// ValidationError is returned when a secret map fails validation.
type ValidationError struct {
	Path   string
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for %q: %s", e.Path, strings.Join(e.Issues, "; "))
}

// ValidateFunc is a function that inspects a secrets map and returns a list of
// human-readable issues. An empty slice means the secrets are valid.
type ValidateFunc func(secrets map[string]string) []string

// validateClient wraps a SecretReader and applies a ValidateFunc to every
// successful read, returning a ValidationError when issues are found.
type validateClient struct {
	inner    SecretReader
	validate ValidateFunc
}

// NewValidateClient returns a SecretReader that applies fn to every result
// returned by inner. If fn reports issues the read fails with a
// *ValidationError. Panics if inner is nil or fn is nil.
func NewValidateClient(inner SecretReader, fn ValidateFunc) SecretReader {
	if inner == nil {
		panic("vault: NewValidateClient: inner must not be nil")
	}
	if fn == nil {
		panic("vault: NewValidateClient: fn must not be nil")
	}
	return &validateClient{inner: inner, validate: fn}
}

func (v *validateClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := v.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	if issues := v.validate(secrets); len(issues) > 0 {
		return nil, &ValidationError{Path: path, Issues: issues}
	}
	return secrets, nil
}

// RequireKeys returns a ValidateFunc that fails when any of the given keys are
// absent from the secrets map.
func RequireKeys(keys ...string) ValidateFunc {
	return func(secrets map[string]string) []string {
		var issues []string
		for _, k := range keys {
			if _, ok := secrets[k]; !ok {
				issues = append(issues, fmt.Sprintf("missing required key %q", k))
			}
		}
		return issues
	}
}

// RejectEmptyValues returns a ValidateFunc that fails when any secret value is
// an empty string.
func RejectEmptyValues() ValidateFunc {
	return func(secrets map[string]string) []string {
		var issues []string
		for k, v := range secrets {
			if strings.TrimSpace(v) == "" {
				issues = append(issues, fmt.Sprintf("key %q has empty value", k))
			}
		}
		return issues
	}
}

// IsValidationError reports whether err (or its chain) is a *ValidationError.
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}
