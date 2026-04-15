package vault

import (
	"errors"
	"sort"
	"testing"
)

func TestPool_RegisterAndGet(t *testing.T) {
	p := NewPool()
	mc := NewMockClient(map[string]map[string]string{})
	if err := p.Register("prod", mc); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := p.Get("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != mc {
		t.Error("Get did not return the registered client")
	}
}

func TestPool_Register_EmptyName(t *testing.T) {
	p := NewPool()
	err := p.Register("", NewMockClient(nil))
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestPool_Register_NilClient(t *testing.T) {
	p := NewPool()
	err := p.Register("prod", nil)
	if err == nil {
		t.Fatal("expected error for nil client")
	}
}

func TestPool_Get_Unknown(t *testing.T) {
	p := NewPool()
	_, err := p.Get("missing")
	if err == nil {
		t.Fatal("expected error for unknown name")
	}
}

func TestPool_Remove(t *testing.T) {
	p := NewPool()
	_ = p.Register("prod", NewMockClient(nil))
	p.Remove("prod")
	_, err := p.Get("prod")
	if err == nil {
		t.Fatal("expected error after removal")
	}
}

func TestPool_Names(t *testing.T) {
	p := NewPool()
	_ = p.Register("b", NewMockClient(nil))
	_ = p.Register("a", NewMockClient(nil))
	names := p.Names()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestPool_ReadSecrets_Dispatches(t *testing.T) {
	p := NewPool()
	mc := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "value"},
	})
	_ = p.Register("prod", mc)

	got, err := p.ReadSecrets("prod/secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("expected value, got %q", got["key"])
	}
}

func TestPool_ReadSecrets_MissingSlash(t *testing.T) {
	p := NewPool()
	_, err := p.ReadSecrets("noslash")
	if err == nil {
		t.Fatal("expected error for path without slash")
	}
}

func TestPool_ReadSecrets_UnknownClient(t *testing.T) {
	p := NewPool()
	_, err := p.ReadSecrets("ghost/secret/app")
	if err == nil {
		t.Fatal("expected error for unregistered client")
	}
}

func TestPool_ReadSecrets_PropagatesError(t *testing.T) {
	p := NewPool()
	sentinel := errors.New("vault down")
	mc := NewMockClient(nil)
	mc.SetError("secret/app", sentinel)
	_ = p.Register("prod", mc)

	_, err := p.ReadSecrets("prod/secret/app")
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestPool_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewPool()
}
