package vault

import (
	"errors"
	"testing"
)

func TestChainClient_PanicsOnNoClients(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for empty client list")
		}
	}()
	NewChainClient()
}

func TestChainClient_FirstSucceeds(t *testing.T) {
	want := map[string]string{"key": "value"}
	a := NewMockClient(map[string]map[string]string{
		"secret/app": want,
	})
	b := NewMockClient(nil)

	chain := NewChainClient(a, b)
	got, err := chain.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != "value" {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestChainClient_FallsToSecond(t *testing.T) {
	want := map[string]string{"k": "v"}
	a := &errorReader{err: errors.New("not found")}
	b := NewMockClient(map[string]map[string]string{
		"secret/app": want,
	})

	chain := NewChainClient(a, b)
	got, err := chain.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["k"] != "v" {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestChainClient_AllFail(t *testing.T) {
	sentinel := errors.New("last error")
	a := &errorReader{err: errors.New("first error")}
	b := &errorReader{err: sentinel}

	chain := NewChainClient(a, b)
	_, err := chain.ReadSecrets("secret/app")
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestChainClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewChainClient(NewMockClient(nil))
}

func TestChainClient_SingleClient(t *testing.T) {
	want := map[string]string{"a": "1"}
	c := NewMockClient(map[string]map[string]string{"p": want})
	chain := NewChainClient(c)
	got, err := chain.ReadSecrets("p")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["a"] != "1" {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// errorReader is a minimal SecretReader that always returns an error.
type errorReader struct{ err error }

func (e *errorReader) ReadSecrets(_ string) (map[string]string, error) {
	return nil, e.err
}
