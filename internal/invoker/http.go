package invoker

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"github.com/kyriosdata/assinatura/internal/jdk"
)

type SignResponse struct {
	Signature string `json:"signature"`
	Valid     bool   `json:"valid"`
	Message   string `json:"message"`
	Error     bool   `json:"error"`
}

type ExecResult struct {
	Stdout string
	Stderr string
	Err    error
}

func detectJava() string {
	javaPath, err := jdk.EnsureJava()
	if err == nil && javaPath != "" {
		return javaPath
	}

	fmt.Println("⚠️ Aviso: Falha ao provisionar JRE isolado. Tentando usar o Java do sistema...")
	javaCmd := "java"
	if runtime.GOOS == "windows" {
		javaCmd = "java.exe"
	}
	return javaCmd
}

func execJava(args ...string) ExecResult {
	cmd := exec.Command(detectJava(), args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return ExecResult{Stdout: strings.TrimSpace(stdout.String()), Stderr: strings.TrimSpace(stderr.String()), Err: err}
}

func Sign(content, token string) (string, error) {
	if content == "" {
		return "", errors.New("content é obrigatório")
	}
	if token == "" {
		return "", errors.New("token é obrigatório")
	}

	result := execJava("-jar", "assinador.jar", "sign", "--content", content, "--token", token)
	if result.Err != nil {
		return "", fmt.Errorf("erro ao executar assinador.jar: %w\n%s", result.Err, result.Stderr)
	}

	var response SignResponse
	if err := json.Unmarshal([]byte(result.Stdout), &response); err != nil {
		return "", fmt.Errorf("falha ao analisar resposta do assinador.jar: %w", err)
	}

	if response.Error {
		return "", errors.New(response.Message)
	}

	return response.Signature, nil
}

func Validate(content, signature string) (bool, error) {
	if content == "" {
		return false, errors.New("content é obrigatório")
	}
	if signature == "" {
		return false, errors.New("signature é obrigatório")
	}

	result := execJava("-jar", "assinador.jar", "validate", "--content", content, "--signature", signature)
	if result.Err != nil {
		return false, fmt.Errorf("erro ao executar assinador.jar: %w\n%s", result.Err, result.Stderr)
	}

	var response SignResponse
	if err := json.Unmarshal([]byte(result.Stdout), &response); err != nil {
		return false, fmt.Errorf("falha ao analisar resposta do assinador.jar: %w", err)
	}

	if response.Error {
		return false, errors.New(response.Message)
	}

	return response.Valid, nil
}
