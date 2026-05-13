package cli

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/invoker"
	"github.com/spf13/cobra"
)

var signature string

var validateCmd = &cobra.Command{
	Use: "validate",
	Run: func(cmd *cobra.Command, args []string) {

		valid, err := invoker.Validate(content, signature)
		if err != nil {
			fmt.Println("Erro:", err)
			return
		}

		fmt.Println("✔ Válido:", valid)
	},
}

func init() {
	validateCmd.Flags().StringVar(&content, "content", "", "Conteúdo")
	validateCmd.Flags().StringVar(&signature, "signature", "", "Assinatura")
	rootCmd.AddCommand(validateCmd)
}
