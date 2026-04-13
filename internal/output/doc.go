// Package output provides formatting and writing utilities for vaultdiff results.
//
// It supports two output formats:
//
//   - text: human-readable, color-coded diff lines with a summary footer.
//   - json: machine-readable JSON containing all diff entries and a summary.
//
// # Text Format
//
// Text output prefixes each changed line with a symbol indicating the type of
// change: "+" for added keys, "-" for removed keys, and "~" for modified keys.
// A summary footer is printed at the end showing the total counts.
//
// # JSON Format
//
// JSON output emits a single object with a "summary" field and an "entries"
// array. Each entry contains the key, change type, and old/new values where
// applicable. This format is suitable for programmatic consumption or CI
// pipelines that parse structured output.
//
// # Usage
//
//	format, err := output.ParseFormat(flagValue)
//	printer := output.NewPrinter(os.Stdout, format, noColor)
//	printer.PrintLine("+", "MY_KEY", "value")
//	printer.PrintSummary(added, removed, modified)
//
// For JSON output use WriteJSON directly:
//
//	err := output.WriteJSON(os.Stdout, srcPath, dstPath, entries)
package output
