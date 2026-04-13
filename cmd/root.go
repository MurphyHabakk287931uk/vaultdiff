package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/your-org/vaultdiff/internal/config"
)

var (
	cfgFile string
	Cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff secrets between two Vault paths or environments",
	Long: `vaultdiff compares secrets stored at two Vault paths and reports
added, removed, and modified keys with optional value redaction.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		Cfg, err = config.Load(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "",
		"path to config file (default: no file, built-in defaults used)",
	)
}
