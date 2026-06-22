package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultdiff",
	Short: "Diff and audit changes between HashiCorp Vault secret versions",
	Long: `vaultdiff is a CLI tool for comparing Vault secret versions across
environments, paths, and time ranges. It helps operators audit changes,
detect drift, and review secret history with configurable output formats.`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
