package cmd

import (
	"fmt"
	"strings"

	"github.com/your-org/vaultdiff/internal/vault"
)

// parseRewriteRules parses a slice of "from=to" strings into RewriteRule values.
// Each entry must contain exactly one "=" separator.
func parseRewriteRules(raw []string) ([]vault.RewriteRule, error) {
	rules := make([]vault.RewriteRule, 0, len(raw))
	for _, r := range raw {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return nil, fmt.Errorf("invalid rewrite rule %q: expected format from=to", r)
		}
		rules = append(rules, vault.RewriteRule{From: parts[0], To: parts[1]})
	}
	return rules, nil
}
