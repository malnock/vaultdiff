package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information, set at build time
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd is the base command for the vaultdiff CLI
var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and audit changes between HashiCorp Vault secret versions",
	Long: `vaultdiff is a CLI tool for comparing Vault secret versions across
environments. It helps teams audit changes, detect drift, and track
secret modifications over time.

Examples:
  # Compare two versions of a secret
  vaultdiff diff secret/myapp/config --v1 3 --v2 4

  # Compare secrets across environments
  vaultdiff compare secret/myapp/config --env1 staging --env2 production

  # Show audit history for a secret path
  vaultdiff audit secret/myapp/config --limit 10`,
}

// versionCmd prints the current version of vaultdiff
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vaultdiff %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

func init() {
	// Global flags available to all subcommands
	rootCmd.PersistentFlags().String("vault-addr", "", "Vault server address (overrides VAULT_ADDR env var)")
	rootCmd.PersistentFlags().String("vault-token", "", "Vault token (overrides VAULT_TOKEN env var)")
	rootCmd.PersistentFlags().String("vault-namespace", "", "Vault namespace (overrides VAULT_NAMESPACE env var)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().String("output", "text", "Output format: text, json, yaml")

	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
