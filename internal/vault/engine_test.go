package vault_test

import (
	"testing"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func TestParseEngineType_Valid(t *testing.T) {
	cases := []struct {
		input    string
		want     vault.EngineType
		wantStr  string
	}{
		{"kv1", vault.EngineKV1, "kv1"},
		{"1", vault.EngineKV1, "kv1"},
		{"kv2", vault.EngineKV2, "kv2"},
		{"2", vault.EngineKV2, "kv2"},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := vault.ParseEngineType(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
			if got.String() != tc.wantStr {
				t.Errorf("String() = %q, want %q", got.String(), tc.wantStr)
			}
		})
	}
}

func TestParseEngineType_Invalid(t *testing.T) {
	_, err := vault.ParseEngineType("kv3")
	if err == nil {
		t.Fatal("expected error for invalid engine type, got nil")
	}
}

func TestEngineType_String_Unknown(t *testing.T) {
	var e vault.EngineType // zero value
	if e.String() != "unknown" {
		t.Errorf("expected \"unknown\", got %q", e.String())
	}
}

func TestNewEngineClient_KV1(t *testing.T) {
	base := vault.NewMockClient(nil, nil)
	client, err := vault.NewEngineClient(base, vault.EngineKV1, []string{"secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewEngineClient_KV2(t *testing.T) {
	base := vault.NewMockClient(nil, nil)
	client, err := vault.NewEngineClient(base, vault.EngineKV2, []string{"secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewEngineClient_Unknown(t *testing.T) {
	base := vault.NewMockClient(nil, nil)
	_, err := vault.NewEngineClient(base, vault.EngineType(99), nil)
	if err == nil {
		t.Fatal("expected error for unknown engine type, got nil")
	}
}
