package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var Version = "dev" // O GoReleaser substituirá isso no build

var rootCmd = &cobra.Command{
	Use:     "assinatura",
	Short:   "CLI para assinatura digital",
	Version: Version, // Habilita automaticamente o flag --version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
