package vault

import (
	"fmt"
	"regexp"
)

// SchemaRule defines a validation rule for a secret key.
type SchemaRule struct {
	Key     string
	Pattern *regexp.Regexp
}

// schemaClient enforces value-format rules against secrets.
type schemaClient struct {
	inner SecretReader
	rules []SchemaRule
}

// NewSchemaClient wraps inner and validates secret values against rules.
// Rules are matched by exact key name. If a value does not match the
// associated pattern, an error is returned.
func NewSchemaClient(inner SecretReader, rules []SchemaRule) SecretReader {
	if inner == nil {
		panic("vault: NewSchemaClient: inner must not be nil")
	}
	if len(rules) == 0 {
		return inner
	}
	return &schemaClient{inner: inner, rules: rules}
}

func (c *schemaClient) ReadSecrets(path string) (map[string]string, error) {
	secrets, err := c.inner.ReadSecrets(path)
	if err != nil {
		return nil, err
	}
	for _, rule := range c.rules {
		val, ok := secrets[rule.Key]
		if !ok {
			continue
		}
		if !rule.Pattern.MatchString(val) {
			return nil, fmt.Errorf("vault: schema violation: key %q value does not match pattern %s", rule.Key, rule.Pattern)
		}
	}
	return secrets, nil
}
