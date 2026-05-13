package cli

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em background",
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("java", "-jar", "assinador.jar")

		err := c.Start() // Inicia sem travar o terminal
		if err != nil {
			fmt.Println("❌ Erro ao iniciar:", err)
			return
		}

		// Tenta salvar o PID para uso futuro do 'stop'
		if err := savePID(c.Process.Pid); err != nil {
			fmt.Println("⚠️ Servidor iniciado, mas falhou ao salvar PID:", err)
		} else {
			fmt.Printf("✔ Servidor iniciado (PID: %d)\n", c.Process.Pid)
		}
	},
}

func savePID(pid int) error {
	home, _ := os.UserHomeDir()
	folder := filepath.Join(home, ".hubsaude")
	os.MkdirAll(folder, 0755)

	path := filepath.Join(folder, "assinador.pid")
	return os.WriteFile(path, []byte(fmt.Sprintf("%d", pid)), 0644)
}

func init() {
	rootCmd.AddCommand(startCmd)
}
