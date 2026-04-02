package main

import (
	"flag"
	"fmt"
	"os"
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
	fmt.Println("(Simulação: O Java subiria aqui)")
}

func pararAssinador() {
	fmt.Println("(Simulação: O Java pararia aqui)")
}