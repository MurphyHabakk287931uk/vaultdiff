package vault

import (
	"regexp"
	"testing"
)

func TestSchemaClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewSchemaClient(nil, []SchemaRule{})
}

func TestSchemaClient_NoRules_ReturnsInner(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"key": "value"},
	})
	client := NewSchemaClient(mock, nil)
	if client != mock {
		t.Fatal("expected inner client returned directly")
	}
}

func TestSchemaClient_PassesValidValue(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"port": "8080"},
	})
	rules := []SchemaRule{
		{Key: "port", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	client := NewSchemaClient(mock, rules)
	secrets, err := client.ReadSecrets("secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["port"] != "8080" {
		t.Errorf("expected 8080, got %s", secrets["port"])
	}
}

func TestSchemaClient_FailsInvalidValue(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"port": "not-a-number"},
	})
	rules := []SchemaRule{
		{Key: "port", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	client := NewSchemaClient(mock, rules)
	_, err := client.ReadSecrets("secret/a")
	if err == nil {
		t.Fatal("expected schema violation error")
	}
}

func TestSchemaClient_SkipsMissingKey(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/a": {"other": "value"},
	})
	rules := []SchemaRule{
		{Key: "port", Pattern: regexp.MustCompile(`^\d+$`)},
	}
	client := NewSchemaClient(mock, rules)
	_, err := client.ReadSecrets("secret/a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSchemaClient_PropagatesInnerError(t *testing.T) {
	mock := NewMockClient(nil)
	rules := []SchemaRule{
		{Key: "k", Pattern: regexp.MustCompile(`.*`)},
	}
	client := NewSchemaClient(mock, rules)
	_, err := client.ReadSecrets("missing/path")
	if err == nil {
		t.Fatal("expected error from inner client")
	}
}

func TestSchemaClient_ImplementsSecretReader(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{})
	rules := []SchemaRule{{Key: "x", Pattern: regexp.MustCompile(`.*`)}}
	var _ SecretReader = NewSchemaClient(mock, rules)
}
