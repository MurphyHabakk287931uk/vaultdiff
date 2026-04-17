package vault

import (
	"errors"
	"testing"
)

func TestSplitClient_PanicsOnNilLeft(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewSplitClient(nil, "left", NewMockClient(nil, nil), "right")
}

func TestSplitClient_PanicsOnNilRight(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewSplitClient(NewMockClient(nil, nil), "left", nil, "right")
}

func TestSplitClient_PanicsOnEmptyNamespace(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewSplitClient(NewMockClient(nil, nil), "", NewMockClient(nil, nil), "right")
}

func TestSplitClient_MergesWithNamespacedKeys(t *testing.T) {
	left := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "left-pass"},
	}, nil)
	right := NewMockClient(map[string]map[string]string{
		"secret/app": {"DB_PASS": "right-pass"},
	}, nil)
	c := NewSplitClient(left, "prod", right, "staging")
	out, err := c.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["prod/DB_PASS"] != "left-pass" {
		t.Errorf("expected prod/DB_PASS=left-pass, got %q", out["prod/DB_PASS"])
	}
	if out["staging/DB_PASS"] != "right-pass" {
		t.Errorf("expected staging/DB_PASS=right-pass, got %q", out["staging/DB_PASS"])
	}
}

func TestSplitClient_PropagatesLeftError(t *testing.T) {
	left := NewMockClient(nil, errors.New("left boom"))
	right := NewMockClient(map[string]map[string]string{"p": {"k": "v"}}, nil)
	c := NewSplitClient(left, "a", right, "b")
	_, err := c.ReadSecrets("p")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSplitClient_PropagatesRightError(t *testing.T) {
	left := NewMockClient(map[string]map[string]string{"p": {"k": "v"}}, nil)
	right := NewMockClient(nil, errors.New("right boom"))
	c := NewSplitClient(left, "a", right, "b")
	_, err := c.ReadSecrets("p")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSplitClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewSplitClient(
		NewMockClient(nil, nil), "a",
		NewMockClient(nil, nil), "b",
	)
}
