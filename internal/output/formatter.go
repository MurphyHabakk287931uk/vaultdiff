package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

// Format represents the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// ParseFormat parses a string into a Format, returning an error if invalid.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case string(FormatText), "":
		return FormatText, nil
	case string(FormatJSON):
		return FormatJSON, nil
	default:
		return "", fmt.Errorf("unknown output format %q: must be one of [text, json]", s)
	}
}

// Printer writes formatted diff output to a writer.
type Printer struct {
	w      io.Writer
	format Format
	noColor bool
}

// NewPrinter creates a new Printer for the given writer and format.
func NewPrinter(w io.Writer, format Format, noColor bool) *Printer {
	return &Printer{w: w, format: format, noColor: noColor}
}

// PrintLine writes a single diff line with optional color coding.
func (p *Printer) PrintLine(prefix, key, value string) {
	if p.format == FormatJSON {
		return
	}
	line := fmt.Sprintf("%s %s = %s", prefix, key, value)
	if p.noColor {
		fmt.Fprintln(p.w, line)
		return
	}
	switch prefix {
	case "+":
		color.New(color.FgGreen).Fprintln(p.w, line)
	case "-":
		color.New(color.FgRed).Fprintln(p.w, line)
	case "~":
		color.New(color.FgYellow).Fprintln(p.w, line)
	default:
		fmt.Fprintln(p.w, line)
	}
}

// PrintSummary writes a summary line to the writer.
func (p *Printer) PrintSummary(added, removed, modified int) {
	if p.format == FormatJSON {
		return
	}
	fmt.Fprintf(p.w, "\nSummary: +%d added, -%d removed, ~%d modified\n", added, removed, modified)
}
