package cli

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/invoker"
	"github.com/spf13/cobra"
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Cria assinatura digital simulada",
	Run: func(cmd *cobra.Command, args []string) {
		if content == "" {
			fmt.Println("Erro: --content é obrigatório")
			cmd.Help()
			return
		}
		if token == "" {
			fmt.Println("Erro: --token é obrigatório")
			cmd.Help()
			return
		}

		resp, err := invoker.Sign(content, token)
		if err != nil {
			fmt.Println("Erro:", err)
			return
		}

		fmt.Println("✔ Assinatura:")
		fmt.Println(resp)
	},
}

func init() {
	signCmd.Flags().StringVar(&content, "content", "", "Conteúdo a ser assinado")
	signCmd.Flags().StringVar(&token, "token", "", "Token ou credencial")
	rootCmd.AddCommand(signCmd)
}
