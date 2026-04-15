package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultdiff/internal/config"
	"github.com/your-org/vaultdiff/internal/diff"
	"github.com/your-org/vaultdiff/internal/output"
	"github.com/your-org/vaultdiff/internal/redact"
	"github.com/your-org/vaultdiff/internal/vault"
)

var (
	redactFlag string
	formatFlag string
	showAllFlag bool
	scopeFlag string
)

var diffCmd = &cobra.Command{
	Use:   "diff <src-path> <dst-path>",
	Short: "Diff secrets between two Vault paths",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

func init() {
	diffCmd.Flags().StringVar(&redactFlag, "redact", "none", "Redact mode: none, redact, mask")
	diffCmd.Flags().StringVar(&formatFlag, "format", "text", "Output format: text, json")
	diffCmd.Flags().BoolVar(&showAllFlag, "show-all", false, "Include unchanged keys in output")
	diffCmd.Flags().StringVar(&scopeFlag, "scope", "", "Restrict reads to a path scope prefix")
	rootCmd.AddCommand(diffCmd)
}

func runDiff(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load("")
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	mode, err := redact.ParseMode(redactFlag)
	if err != nil {
		return err
	}

	fmt, err := output.ParseFormat(formatFlag)
	if err != nil {
		return err
	}

	baseClient, err := vault.NewClient(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	var src, dst vault.SecretReader
	src = baseClient
	dst = baseClient

	if scopeFlag != "" {
		src = vault.NewScopeClient(src, scopeFlag)
		dst = vault.NewScopeClient(dst, scopeFlag)
	}

	runner := diff.NewRunner(diff.RunnerConfig{
		Src:        src,
		Dst:        dst,
		SrcPath:    args[0],
		DstPath:    args[1],
		RedactMode: mode,
		ShowAll:    showAllFlag,
	})

	results, err := runner.Run(cmd.Context())
	if err != nil {
		return err
	}

	printer := output.NewPrinter(fmt, os.Stdout)
	for _, r := range results {
		printer.PrintLine(r)
	}
	printer.PrintSummary(results)
	return nil
}
