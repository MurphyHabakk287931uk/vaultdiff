package vault

import "fmt"

// SplitClient reads from two named inner clients and returns their secrets
// under separate top-level key namespaces, allowing side-by-side comparison
// of secrets from two distinct sources in a single map.
type SplitClient struct {
	left  SecretReader
	right SecretReader
	leftNS  string
	rightNS string
}

// NewSplitClient wraps two SecretReaders, prefixing their keys with the
// provided namespace strings so results can coexist in one map.
func NewSplitClient(left SecretReader, leftNS string, right SecretReader, rightNS string) *SplitClient {
	if left == nil {
		panic("vault: NewSplitClient: left client must not be nil")
	}
	if right == nil {
		panic("vault: NewSplitClient: right client must not be nil")
	}
	if leftNS == "" || rightNS == "" {
		panic("vault: NewSplitClient: namespace strings must not be empty")
	}
	return &SplitClient{left: left, leftNS: leftNS, right: right, rightNS: rightNS}
}

// ReadSecrets reads the same path from both inner clients and returns a merged
// map whose keys are prefixed with each namespace, e.g. "left/KEY".
func (s *SplitClient) ReadSecrets(path string) (map[string]string, error) {
	lSecrets, err := s.left.ReadSecrets(path)
	if err != nil {
		return nil, fmt.Errorf("split(%s): %w", s.leftNS, err)
	}
	rSecrets, err := s.right.ReadSecrets(path)
	if err != nil {
		return nil, fmt.Errorf("split(%s): %w", s.rightNS, err)
	}
	out := make(map[string]string, len(lSecrets)+len(rSecrets))
	for k, v := range lSecrets {
		out[s.leftNS+"/"+k] = v
	}
	for k, v := range rSecrets {
		out[s.rightNS+"/"+k] = v
	}
	return out, nil
}
