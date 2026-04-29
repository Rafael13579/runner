package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	startCmd := flag.NewFlagSet("start", flag.ExitOnError)
	stopCmd := flag.NewFlagSet("stop", flag.ExitOnError)
	signCmd := flag.NewFlagSet("sign", flag.ExitOnError)

	signData := signCmd.String("data", "", "Conteúdo a ser assinado")

	if len(os.Args) < 2 {
		fmt.Println("Uso esperado: assinatura [comando]")
		fmt.Println("Comandos disponíveis: start, stop, sign")
		os.Exit(1)
	}

	switch os.Args[1] {

	case "start":
		startCmd.Parse(os.Args[2:])
		fmt.Println("Iniciando o Assinador Java...")
		iniciarAssinador()

	case "stop":
		stopCmd.Parse(os.Args[2:])
		fmt.Println("Parando o Assinador Java...")
		pararAssinador()

	case "sign":
		signCmd.Parse(os.Args[2:])
		if *signData == "" {
			fmt.Println("Erro: Use --data para enviar o conteúdo.")
			signCmd.PrintDefaults()
			os.Exit(1)
		}
		fmt.Printf("Assinando dados: %s\n", *signData)
	

	default:
		fmt.Println("Comando desconhecido.")
		os.Exit(1)
	}
}

func iniciarAssinador() {
	cmd := exec.Command("java", "-jar", "assinador.jar")

	err := cmd.Start()
	if err != nil {
		fmt.Printf("Falha ao iniciar o Java: %v\n", err)
		return
	}

	pid := strconv.Itoa(cmd.Process.Pid)
	err = ioutil.WriteFile("assinador.pid", []byte(pid), 0644)
	if err != nil {
		fmt.Printf("Java iniciou (PID %s), mas falhei ao salvar o arquivo de controle.\n", pid)
	}

	fmt.Printf("Assinador rodando em background! (PID: %s)\n", pid)
	fmt.Println("Acesse em: http://localhost:8080")
}


func pararAssinador() {
	// 1. Lê o arquivo de PID
	data, err := ioutil.ReadFile("assinador.pid")
	if err != nil {
		fmt.Println("❌ Erro: O Assinador não parece estar rodando (arquivo .pid não encontrado).")
		return
	}

	pid, _ := strconv.Atoi(string(data))

	// 2. Tenta encontrar o processo e matá-lo
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("❌ Não encontrei o processo com PID %d.\n", pid)
		return
	}

	err = process.Kill()
	if err != nil {
		fmt.Printf("❌ Falha ao encerrar o processo: %v\n", err)
	} else {
		fmt.Println("✅ Assinador encerrado com sucesso.")
		os.Remove("assinador.pid") // Limpa o rastro
	}
}
