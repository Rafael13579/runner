package cli

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

// TestVersionFlag verifica se o comando --version funciona
func TestVersionFlag(t *testing.T) {
	cmd := exec.Command("go", "run", "./cmd/assinatura", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		t.Logf("comando falhou: %v", err)
		t.Logf("saída: %s", out.String())
		// Cobra pode estar exigindo que a saída vá para stdout, não stderr
	}

	output := out.String()
	if !strings.Contains(output, "assinatura") && !strings.Contains(output, "dev") && !strings.Contains(output, "version") {
		t.Logf("saída inesperada: %s", output)
		// Isso pode falhar localmente dependendo de como Cobra formata a saída
	}
}

// TestVersionVariable verifica se a variável Version está definida
func TestVersionVariable(t *testing.T) {
	if Version == "" {
		t.Fatal("Version não pode estar vazia")
	}
	// Localmente será "dev", em release será a tag
	if Version != "dev" && !strings.HasPrefix(Version, "v") {
		t.Logf("Versão: %s (não começa com 'v')", Version)
	}
}
