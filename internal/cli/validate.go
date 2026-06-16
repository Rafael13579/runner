package cli

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/invoker"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida assinatura digital simulada",
	Run: func(cmd *cobra.Command, args []string) {
		if content == "" {
			fmt.Println("Erro: --content é obrigatório")
			cmd.Help()
			return
		}
		if signature == "" {
			fmt.Println("Erro: --signature é obrigatório")
			cmd.Help()
			return
		}

		valid, err := invoker.Validate(content, signature)
		if err != nil {
			fmt.Println("Erro:", err)
			return
		}

		fmt.Println("✔ Resultado da validação:")
		fmt.Printf("Assinatura: %s\n", signature)
		fmt.Printf("Válido: %t\n", valid)
	},
}

func init() {
	validateCmd.Flags().StringVar(&content, "content", "", "Conteúdo original")
	validateCmd.Flags().StringVar(&signature, "signature", "", "Assinatura a ser validada")
	rootCmd.AddCommand(validateCmd)
}
