package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	vaultAddr  string
	vaultToken string
	redactMode string
	outputFmt  string
	showAll    bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff secrets between two Vault paths or environments",
	Long: `vaultdiff compares secrets stored at two Vault paths and outputs
the differences with optional redaction of sensitive values.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&vaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.PersistentFlags().StringVar(&redactMode, "redact", "none", "Redact mode: none, redact, mask")
	rootCmd.PersistentFlags().StringVar(&outputFmt, "output", "text", "Output format: text, json")
	rootCmd.PersistentFlags().BoolVar(&showAll, "all", false, "Show unchanged keys as well")
}
