package vault

import (
	"errors"
	"testing"
)

func TestCompressClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewCompressClient(nil)
}

func TestCompressClient_CompressesValues(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"key": "hello world"},
	})
	client := NewCompressClient(mock)
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["key"] == "hello world" {
		t.Fatal("expected compressed value, got original")
	}
	if got["key"] == "" {
		t.Fatal("expected non-empty compressed value")
	}
}

func TestCompressClient_RoundTrip(t *testing.T) {
	original := "the quick brown fox jumps over the lazy dog"
	compressed, err := compressValue(original)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	got, err := decompressValue(compressed)
	if err != nil {
		t.Fatalf("decompress: %v", err)
	}
	if got != original {
		t.Fatalf("round-trip mismatch: got %q want %q", got, original)
	}
}

func TestCompressClient_PropagatesError(t *testing.T) {
	mock := NewMockClient(nil)
	mock.SetError("secret/missing", errors.New("not found"))
	client := NewCompressClient(mock)
	_, err := client.ReadSecrets("secret/missing")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCompressClient_EmptyValue(t *testing.T) {
	mock := NewMockClient(map[string]map[string]string{
		"secret/app": {"empty": ""},
	})
	client := NewCompressClient(mock)
	got, err := client.ReadSecrets("secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	decoded, err := decompressValue(got["empty"])
	if err != nil {
		t.Fatalf("decompress: %v", err)
	}
	if decoded != "" {
		t.Fatalf("expected empty string, got %q", decoded)
	}
}

func TestCompressClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewCompressClient(NewMockClient(nil))
}
