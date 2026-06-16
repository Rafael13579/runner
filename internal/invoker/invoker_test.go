package invoker

import (
	"os"
	"strings"
	"testing"
)

func checkJar(t *testing.T) {
	if _, err := os.Stat("assinador.jar"); os.IsNotExist(err) {
		t.Skip("⚠️ Teste pulado: 'assinador.jar' não encontrado neste diretório.")
	}
}

func TestSign_Success(t *testing.T) {
	checkJar(t)

	signature, err := Sign("conteudo_teste", "token123")
	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}

	expected := "MOCKED_SIGNATURE_BASE64_=="
	if signature != expected {
		t.Errorf("Esperava %q, obteve %q", expected, signature)
	}
}

func TestSign_MissingContent(t *testing.T) {
	checkJar(t)

	_, err := Sign("", "token123")
	if err == nil {
		t.Fatal("Esperava erro por falta de conteúdo, mas a assinatura foi gerada.")
	}

	if !strings.Contains(err.Error(), "obrigatório") {
		t.Errorf("Mensagem de erro inesperada: %v", err)
	}
}
