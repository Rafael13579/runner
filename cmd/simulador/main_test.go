package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	root := newRootCommand()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"--version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("comando --version falhou: %v", err)
	}

	if !strings.Contains(out.String(), "simulador") {
		t.Fatalf("saida de versao inesperada: %q", out.String())
	}
}
