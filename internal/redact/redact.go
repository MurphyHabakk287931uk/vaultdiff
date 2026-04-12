package redact

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

const (
	// RedactedPlaceholder is shown when a value is fully redacted
	RedactedPlaceholder = "[REDACTED]"
	// MaskedPrefix is the prefix shown for hashed values
	MaskedPrefix = "sha256:"
)

// Mode controls how secret values are displayed in diff output
type Mode int

const (
	// ModeNone shows values in plaintext (default)
	ModeNone Mode = iota
	// ModeRedact replaces values with [REDACTED]
	ModeRedact
	// ModeMask replaces values with a short SHA256 hash prefix
	ModeMask
)

// ParseMode converts a string flag value into a Mode.
func ParseMode(s string) (Mode, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "none":
		return ModeNone, nil
	case "redact":
		return ModeRedact, nil
	case "mask":
		return ModeMask, nil
	default:
		return ModeNone, fmt.Errorf("unknown redact mode %q: must be none, redact, or mask", s)
	}
}

// Apply transforms a secret value according to the given Mode.
func Apply(value string, mode Mode) string {
	switch mode {
	case ModeRedact:
		return RedactedPlaceholder
	case ModeMask:
		sum := sha256.Sum256([]byte(value))
		// Show first 8 hex chars so users can confirm whether two values match
		return MaskedPrefix + fmt.Sprintf("%x", sum[:4])
	default:
		return value
	}
}
