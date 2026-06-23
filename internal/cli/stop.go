package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Encerra o processo do assinador HTTP em background",
	Run: func(cmd *cobra.Command, args []string) {
		home, _ := os.UserHomeDir()
		pidFile := filepath.Join(home, ".hubsaude", "assinador.pid")

		data, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Println("ℹ️ Nenhum servidor em execução detectado (arquivo PID ausente).")
			return
		}

		pid, err := strconv.Atoi(string(data))
		if err != nil {
			fmt.Println("❌ Erro ao ler o PID. O arquivo pode estar corrompido.")
			return
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			fmt.Printf("ℹ️ Processo %d não encontrado. Limpando arquivo de registro...\n", pid)
			os.Remove(pidFile)
			return
		}
		err = process.Kill()
		if err != nil {
	
			fmt.Printf("ℹ️ O processo %d já foi encerrado ou o acesso foi negado.\n", pid)
		} else {
			fmt.Printf("✔ Servidor (PID: %d) encerrado com sucesso.\n", pid)
		}


		os.Remove(pidFile)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}