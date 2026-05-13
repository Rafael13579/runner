package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra o processo do assinador",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		pidFile := filepath.Join(home, ".hubsaude", "assinador.pid")

		// 1. Lê o arquivo
		data, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Println("❌ Nenhum processo em execução encontrado (arquivo PID ausente).")
			return
		}

		// 2. Converte para inteiro
		pid, _ := strconv.Atoi(string(data))

		// 3. Busca o processo e envia sinal de interrupção
		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("❌ Não foi possível encontrar o processo %d\n", pid)
			return
		}

		err = process.Signal(os.Interrupt) // Envia o sinal de "parar"
		if err != nil {
			fmt.Println("❌ Erro ao parar o processo:", err)
			return
		}

		// 4. Limpa o arquivo de PID
		os.Remove(pidFile)
		fmt.Printf("✔ Processo %d encerrado com sucesso.\n", pid)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
