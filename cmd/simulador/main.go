package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kyriosdata/assinatura/internal/simulador"
	"github.com/spf13/cobra"
)

var (
	version     = "dev"
	port        int
	sourceURL   string
	checksum    string
	healthPath  string
	waitSeconds int
)

func main() {
	root := newRootCommand()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:     "simulador",
		Short:   "Gerencia o ciclo de vida do Simulador do HubSaude",
		Version: version,
	}

	root.PersistentFlags().IntVar(&port, "port", 8090, "Porta HTTP do simulador")
	root.PersistentFlags().StringVar(&healthPath, "health-path", "/actuator/health", "Caminho do health check do simulador")

	root.AddCommand(newStartCommand())
	root.AddCommand(newStopCommand())
	root.AddCommand(newStatusCommand())
	return root
}

func newStartCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Inicia o simulador.jar",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := options()
			status, err := simulador.Start(cmd.Context(), opts)
			if err != nil {
				return err
			}
			printStatus(status)
			return nil
		},
	}

	cmd.Flags().StringVar(&sourceURL, "source", "", "URL alternativa para baixar simulador.jar")
	cmd.Flags().StringVar(&checksum, "checksum", "", "SHA-256 esperado para validar o simulador.jar")
	cmd.Flags().IntVar(&waitSeconds, "wait", 30, "Tempo maximo, em segundos, para aguardar readiness")
	return cmd
}

func newStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Para o simulador registrado",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := simulador.Stop()
			if err != nil {
				return err
			}
			printStatus(status)
			return nil
		},
	}
}

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Mostra o estado atual do simulador",
		RunE: func(cmd *cobra.Command, args []string) error {
			status, err := simulador.CurrentStatus(context.Background(), options())
			if err != nil {
				return err
			}
			printStatus(status)
			if !status.Running {
				os.Exit(2)
			}
			return nil
		},
	}
}

func options() simulador.Options {
	return simulador.Options{
		Port:       port,
		SourceURL:  sourceURL,
		Checksum:   checksum,
		HealthPath: healthPath,
		Timeout:    time.Duration(waitSeconds) * time.Second,
	}
}

func printStatus(status simulador.Status) {
	fmt.Println(status.Message)
	if status.Info == nil {
		return
	}
	fmt.Printf("PID: %d\n", status.Info.PID)
	fmt.Printf("Porta: %d\n", status.Info.Port)
	fmt.Printf("JAR: %s\n", status.Info.JarPath)
	fmt.Printf("Pronto: %t\n", status.Ready)
}
