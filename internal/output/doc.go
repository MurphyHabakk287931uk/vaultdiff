// Package output provides formatting and writing utilities for vaultdiff results.
//
// It supports two output formats:
//
//   - text: human-readable, color-coded diff lines with a summary footer.
//   - json: machine-readable JSON containing all diff entries and a summary.
//
// Usage:
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
