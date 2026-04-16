package vault

import (
	"context"
	"errors"
	"testing"
)

var testKey = []byte("0123456789abcdef") // 16-byte AES-128 key

func TestEncryptClient_PanicsOnNilInner(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewEncryptClient(nil, testKey)
}

func TestEncryptClient_PanicsOnBadKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	NewEncryptClient(NewMockClient(nil, nil), []byte("short"))
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	orig := map[string]string{"DB_PASS": "s3cr3t", "API_KEY": "abc123"}
	mock := NewMockClient(orig, nil)

	enc := NewEncryptClient(mock, testKey)
	dec := NewDecryptClient(enc, testKey)

	ctx := context.Background()
	got, err := dec.ReadSecrets(ctx, "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, want := range orig {
		if got[k] != want {
			t.Errorf("key %q: got %q, want %q", k, got[k], want)
		}
	}
}

func TestEncryptClient_ValuesAreNotPlaintext(t *testing.T) {
	orig := map[string]string{"SECRET": "plainvalue"}
	mock := NewMockClient(orig, nil)
	enc := NewEncryptClient(mock, testKey)

	got, err := enc.ReadSecrets(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["SECRET"] == "plainvalue" {
		t.Error("expected encrypted value, got plaintext")
	}
}

func TestEncryptClient_PropagatesError(t *testing.T) {
	expected := errors.New("read failure")
	mock := NewMockClient(nil, expected)
	enc := NewEncryptClient(mock, testKey)

	_, err := enc.ReadSecrets(context.Background(), "secret/app")
	if !errors.Is(err, expected) {
		t.Errorf("got %v, want %v", err, expected)
	}
}

func TestDecryptClient_InvalidCiphertextReturnsError(t *testing.T) {
	bad := map[string]string{"KEY": "not-valid-base64!!"}
	mock := NewMockClient(bad, nil)
	dec := NewDecryptClient(mock, testKey)

	_, err := dec.ReadSecrets(context.Background(), "secret/app")
	if err == nil {
		t.Fatal("expected error for invalid ciphertext")
	}
}

func TestEncryptClient_ImplementsSecretReader(t *testing.T) {
	var _ SecretReader = NewEncryptClient(NewMockClient(nil, nil), testKey)
	var _ SecretReader = NewDecryptClient(NewMockClient(nil, nil), testKey)
}
