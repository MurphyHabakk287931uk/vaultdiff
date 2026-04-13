package vault

import "fmt"

// EngineType represents a supported KV secrets engine version.
type EngineType int

const (
	EngineKV1 EngineType = iota + 1
	EngineKV2
)

// ParseEngineType parses a string into an EngineType.
// Accepts "kv1", "1", "kv2", "2" (case-insensitive).
func ParseEngineType(s string) (EngineType, error) {
	switch s {
	case "kv1", "1":
		return EngineKV1, nil
	case "kv2", "2":
		return EngineKV2, nil
	default:
		return 0, fmt.Errorf("unsupported engine type %q: must be \"kv1\" or \"kv2\"", s)
	}
}

// String returns the canonical string representation of the EngineType.
func (e EngineType) String() string {
	switch e {
	case EngineKV1:
		return "kv1"
	case EngineKV2:
		return "kv2"
	default:
		return "unknown"
	}
}

// NewEngineClient constructs a SecretReader wrapping the given base client
// with the appropriate KV engine adapter and mount paths.
func NewEngineClient(base SecretReader, engine EngineType, mounts []string) (SecretReader, error) {
	switch engine {
	case EngineKV1:
		return NewKV1Client(base, mounts), nil
	case EngineKV2:
		return NewKV2Client(base, mounts), nil
	default:
		return nil, fmt.Errorf("cannot create client for unknown engine type %d", engine)
	}
}
