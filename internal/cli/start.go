package cli

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kyriosdata/assinatura/internal/jdk"
	"github.com/spf13/cobra"
)

var (
	port    string
	timeout int
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o assinador.jar em background no modo servidor",
	Long:  `Inicia o servidor HTTP do assinador em segundo plano, monitoriza a sua inicialização através de health checks e configura janelas de inatividade.`,
	Run: func(cmd *cobra.Command, args []string) {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("❌ Erro ao obter o diretório home do utilizador:", err)
			return
		}

		jarPath := filepath.Join(home, ".hubsaude", "assinador.jar")

		javaBin, err := jdk.EnsureJava()
		if err != nil {
			fmt.Println("❌ Erro ao localizar ou descarregar o JRE:", err)
			return
		}

		if _, err := os.Stat(jarPath); os.IsNotExist(err) {
			if _, errLocal := os.Stat("assinador.jar"); errLocal == nil {
				jarPath = "assinador.jar"
			} else {
				fmt.Println("❌ Erro: 'assinador.jar' não encontrado. Certifique-se de que o artefato está presente.")
				return
			}
		}

		javaArgs := []string{
			fmt.Sprintf("-Dserver.port=%s", port),
			"-jar",
			jarPath,
			fmt.Sprintf("--assinador.timeout.minutes=%d", timeout),
		}

		c := exec.Command(javaBin, javaArgs...)

		err = c.Start()
		if err != nil {
			fmt.Println("❌ Erro ao iniciar o processo do servidor:", err)
			return
		}

		if err := savePID(c.Process.Pid); err != nil {
			fmt.Println("⚠️ Servidor iniciado, mas falhou ao gravar o ficheiro de PID:", err)
		} else {
			fmt.Printf("✔ Processo do servidor disparado em segundo plano (PID: %d). Aguardando inicialização...\n", c.Process.Pid)
		}

		success := false
		for i := 0; i < 15; i++ {
			time.Sleep(1 * time.Second)
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/health", port))
			if err == nil && resp.StatusCode == http.StatusOK {
				success = true
				resp.Body.Close()
				break
			}
		}

		if success {
			fmt.Printf("✔ Servidor ativo e pronto para receber requisições na porta %s!\n", port)
		} else {
			fmt.Println("⚠️ O processo foi iniciado, mas o endpoint de health check não respondeu dentro do tempo limite.")
		}
	},
}

func savePID(pid int) error {
	home, _ := os.UserHomeDir()
	folder := filepath.Join(home, ".hubsaude")
	
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}

	path := filepath.Join(folder, "assinador.pid")
	return os.WriteFile(path, []byte(fmt.Sprintf("%d", pid)), 0644)
}

func init() {
	startCmd.Flags().StringVar(&port, "port", "8080", "Porta para execução do servidor HTTP")
	startCmd.Flags().IntVar(&timeout, "timeout", 0, "Encerra o servidor após X minutos de inatividade (0 para desativar)")
	rootCmd.AddCommand(startCmd)
}
