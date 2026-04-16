// Package vault provides a rewriteClient that rewrites secret paths before
// forwarding reads to an inner SecretReader.
//
// # Usage
//
//	rules := []vault.RewriteRule{
//		{From: "staging/", To: "prod/"},
//	}
//	client := vault.NewRewriteClient(inner, rules)
//	// ReadSecrets("staging/db") → inner.ReadSecrets("prod/db")
//
// Rules are evaluated in declaration order; the first matching prefix wins.
// If no rules are provided, the inner client is returned directly with no
// overhead.
package vault
