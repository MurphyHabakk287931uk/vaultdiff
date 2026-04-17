package cmd

import (
	"fmt"
	"strings"

	"github.com/yourusername/vaultdiff/internal/vault"
)

// parseEnrichPairs converts a slice of "key=value" strings into a map
// suitable for vault.NewEnrichClient. Returns an error if any entry is
// malformed.
func parseEnrichPairs(pairs []string) (map[string]string, error) {
	result := make(map[string]string, len(pairs))
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx <= 0 {
			return nil, fmt.Errorf("enrich: invalid pair %q (expected key=value)", p)
		}
		result[p[:idx]] = p[idx+1:]
	}
	return result, nil
}

// buildEnrichClient wraps inner with an EnrichClient when pairs is non-empty.
// If pairs is empty the original client is returned unchanged.
func buildEnrichClient(inner vault.SecretReader, pairs []string) (vault.SecretReader, error) {
	if len(pairs) == 0 {
		return inner, nil
	}
	extra, err := parseEnrichPairs(pairs)
	if err != nil {
		return nil, err
	}
	return vault.NewEnrichClient(inner, extra), nil
}
