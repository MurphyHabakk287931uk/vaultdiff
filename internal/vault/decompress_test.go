package vault

import (
	"errors"
	"testing"
)

func TestDecompressClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewDecompressClient(nil)
}

func TestDecompressClient_DecompressesValues(t *testing.T) {
	original := "secret payload"
	compressed, err := compressValue(original)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": compressed},
	})
	client := NewDecompressClient(mock)
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] != original {
		t.Fatalf("got %q want %q", got["key"], original)
	}
}

func TestDecompressClient_InvalidBase64ReturnsError(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "not-valid-base64!!!"},
	})
	client := NewDecompressClient(mock)
	_, err := client.ReadSecrets("secret/app")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecompressClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/missing", errors.New("not found"))
	client := NewDecompressClient(mock)
	_, err := client.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDecompressClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewDecompressClient(NewMockClient(nil))
}

func TestCompressDecompressRoundTrip_ViaClients(t *testing.T) {
	original := map[string]string{"alpha": "foo bar baz", "beta": "another value"}
	mock := NewMockClient(map[string]map[string]string{"secret/data": original})
	compressed := NewCompressClient(mock)
	// store compressed result in a new mock
	result, err := compressed.ReadSecrets("secret/data")
	if err != nil {
		t.Fatalf("compress read: %v", err)
	}
	mock2 := NewMockClient(map[string]map[string]string{"secret/data": result})
	decompressed := NewDecompressClient(mock2)
	final, err := decompressed.ReadSecrets("secret/data")
	if err != nil {
		t.Fatalf("decompress read: %v", err)
	}
	for k, v := range original {
		if final[k] != v {
			t.Fatalf("key %q: got %q want %q", k, final[k], v)
		}
	}
}
