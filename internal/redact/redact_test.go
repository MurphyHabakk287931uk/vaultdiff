package redact_test

import (
	"strings"
	"testing"

	"github.com/yourorg/vaultdiff/internal/redact"
)

func TestParseMode(t *testing.T) {
	tests := []struct {
		input   string
		want    redact.Mode
		wantErr bool
	}{
		{"none", redact.ModeNone, false},
		{"", redact.ModeNone, false},
		{"redact", redact.ModeRedact, false},
		{"REDACT", redact.ModeRedact, false},
		{"mask", redact.ModeMask, false},
		{"Mask", redact.ModeMask, false},
		{"unknown", redact.ModeNone, true},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			got, err := redact.ParseMode(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("ParseMode(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestApply_None(t *testing.T) {
	got := redact.Apply("supersecret", redact.ModeNone)
	if got != "supersecret" {
		t.Errorf("ModeNone should return plaintext, got %q", got)
	}
}

func TestApply_Redact(t *testing.T) {
	got := redact.Apply("supersecret", redact.ModeRedact)
	if got != redact.RedactedPlaceholder {
		t.Errorf("ModeRedact should return %q, got %q", redact.RedactedPlaceholder, got)
	}
}

func TestApply_Mask(t *testing.T) {
	value := "supersecret"
	got := redact.Apply(value, redact.ModeMask)

	if !strings.HasPrefix(got, redact.MaskedPrefix) {
		t.Errorf("ModeMask output should start with %q, got %q", redact.MaskedPrefix, got)
	}

	// Same value should produce same mask (deterministic)
	got2 := redact.Apply(value, redact.ModeMask)
	if got != got2 {
		t.Errorf("ModeMask should be deterministic: %q != %q", got, got2)
	}

	// Different values should produce different masks
	diff := redact.Apply("different", redact.ModeMask)
	if got == diff {
		t.Errorf("ModeMask should differ for different inputs")
	}
}
