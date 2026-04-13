package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultdiff/internal/diff"
	"github.com/example/vaultdiff/internal/output"
	"github.com/example/vaultdiff/internal/redact"
	"github.com/example/vaultdiff/internal/vault"
)

var diffCmd = &cobra.Command{
	Use:   "diff <src-path> <dst-path>",
	Short: "Compare secrets between two Vault paths",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	srcPath := args[0]
	dstPath := args[1]

	redactModeVal, err := redact.ParseMode(redactMode)
	if err != nil {
		return fmt.Errorf("invalid redact mode: %w", err)
	}

	fmt, err := output.ParseFormat(outputFmt)
	if err != nil {
		return fmt.Errorf("invalid output format: %w", err)
	}

	clientOpts := []vault.Option{}
	if vaultAddr != "" {
		clientOpts = append(clientOpts, vault.WithAddress(vaultAddr))
	}
	if vaultToken != "" {
		clientOpts = append(clientOpts, vault.WithToken(vaultToken))
	}

	client, err := vault.NewClient(clientOpts...)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	printer := output.NewPrinter(os.Stdout, fmt)

	runner := diff.NewRunner(client, client, redactModeVal, printer)
	result, err := runner.Run(cmd.Context(), srcPath, dstPath)
	if err != nil {
		return fmt.Errorf("diff failed: %w", err)
	}

	if fmt == output.FormatText {
		printer.PrintSummary(result.Added, result.Removed, result.Modified, result.Unchanged)
	} else {
		if err := output.WriteJSON(os.Stdout, result.Entries, result.Added, result.Removed, result.Modified, result.Unchanged); err != nil {
			return fmt.Errorf("failed to write JSON output: %w", err)
		}
	}

	return nil
}
