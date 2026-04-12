package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected Format
	}{
		{"text", FormatText},
		{"TEXT", FormatText},
		{"", FormatText},
		{"json", FormatJSON},
		{"JSON", FormatJSON},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := ParseFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
	if !strings.Contains(err.Error(), "yaml") {
		t.Errorf("error should mention the bad value, got: %v", err)
	}
}

func TestPrinter_PrintLine_NoColor(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, FormatText, true)
	p.PrintLine("+", "DB_PASS", "***")
	out := buf.String()
	if !strings.Contains(out, "+ DB_PASS = ***") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestPrinter_PrintLine_JSONSkips(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, FormatJSON, true)
	p.PrintLine("+", "KEY", "val")
	if buf.Len() != 0 {
		t.Errorf("expected no output for JSON format, got: %q", buf.String())
	}
}

func TestPrinter_PrintSummary(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, FormatText, true)
	p.PrintSummary(2, 1, 3)
	out := buf.String()
	if !strings.Contains(out, "+2 added") || !strings.Contains(out, "-1 removed") || !strings.Contains(out, "~3 modified") {
		t.Errorf("unexpected summary output: %q", out)
	}
}

func TestPrinter_PrintSummary_JSONSkips(t *testing.T) {
	var buf bytes.Buffer
	p := NewPrinter(&buf, FormatJSON, true)
	p.PrintSummary(1, 2, 3)
	if buf.Len() != 0 {
		t.Errorf("expected no output for JSON format")
	}
}
