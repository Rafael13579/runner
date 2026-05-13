package cli

import (
	"fmt"

	"github.com/kyriosdata/assinatura/internal/invoker"
	"github.com/spf13/cobra"
)

var content string

var signCmd = &cobra.Command{
	Use: "sign",
	Run: func(cmd *cobra.Command, args []string) {

		resp, err := invoker.Sign(content)
		if err != nil {
			fmt.Println("Erro:", err)
			return
		}

		fmt.Println("✔ Assinatura:", resp)
	},
}

func init() {
	signCmd.Flags().StringVar(&content, "content", "", "Conteúdo")
	rootCmd.AddCommand(signCmd)
}
