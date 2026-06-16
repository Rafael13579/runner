package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	Version   = "dev" // O GoReleaser substituirá isso no build
	content   string
	token     string
	signature string
	jarPath   string
)

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
