# Sistema Runner

Sistema de execução de aplicações Java via linha de comandos, desenvolvido como trabalho prático da disciplina de Implementação e Integração de Software (2026/1) — Bacharelado em Engenharia de Software, UFG.

## Visão Geral

O **Sistema Runner** facilita o acesso à funcionalidade de execução de aplicações Java sem necessidade de conhecimento técnico profundo. Oferece:

- **CLI multiplataforma** (Windows, Linux, macOS) para invocar operações de assinatura digital
- **Modo local**: invocação direta do `assinador.jar` (cold start)
- **Modo servidor**: invocação via HTTP com menor latência (warm start)
- **JDK automático**: detecção e provisionamento automático do Java quando necessário
- **Simulação de assinatura digital** com validação rigorosa de parâmetros

## Requisitos

- **Go 1.25+** (para compilação do CLI)
- **Java 21+** (ou será provisionado automaticamente)

## Instalação

### Via GitHub Releases (Recomendado)

Baixe a versão pré-compilada para sua plataforma em [GitHub Releases](../../releases):

```bash
# Linux
wget https://github.com/kyriosdata/runner/releases/download/v0.1.0/assinatura-v0.1.0-linux-amd64
chmod +x assinatura-v0.1.0-linux-amd64
./assinatura-v0.1.0-linux-amd64 --version

# Windows
# Baixe assinatura-v0.1.0-windows-amd64.exe

# macOS
# Baixe assinatura-v0.1.0-darwin-amd64
```

### Compilação Local

```bash
git clone https://github.com/kyriosdata/runner
cd runner
go build ./cmd/assinatura
go build ./cmd/simulador
./assinatura --version
./simulador --version
```

## Como Usar

### Comando version

```bash
assinatura --version
# Saída: assinatura version dev (local) ou assinatura version v0.1.0 (release)
```

### Comando help

```bash
assinatura --help
```

### Exemplos de Uso (Sprint 2+)

```bash
# Criar assinatura (modo local)
assinatura sign --content "documento.pdf" --local

# Validar assinatura (modo servidor)
assinatura validate --content "documento.pdf" --signature "sig123"

# Iniciar servidor de assinatura
assinatura start

# Parar servidor
assinatura stop
```

### Simulador do HubSaúde

```bash
# Iniciar o simulador na porta padrão 8090
simulador start

# Usar uma porta específica
simulador start --port 9090

# Baixar o simulador.jar de uma URL alternativa e validar SHA-256
simulador start --source https://example.org/simulador.jar --checksum <sha256>

# Consultar status
simulador status

# Parar o simulador registrado
simulador stop
```

O CLI procura `simulador.jar` em `~/.hubsaude/simulador/simulador.jar`.
Quando o arquivo não existe localmente, o comando consulta o último GitHub Release
do projeto e baixa o asset `simulador.jar`; se existir `simulador.jar.sha256`,
o checksum é validado antes da execução.

## Desenvolvimento

### Build local

```bash
go build ./...
go vet ./...
go test ./...
```

### Decisões de arquitetura

As decisões técnicas específicas desta implementação ficam em [`docs/adr`](docs/adr):

- ADR 0001: CLIs em Go
- ADR 0002: Registro de processos em `~/.hubsaude`
- ADR 0003: Obtenção dinâmica do `simulador.jar`

### Testes

```bash
go test ./...
```

### Cross-compilation

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o dist/assinatura-linux-amd64 ./cmd/assinatura
GOOS=linux GOARCH=amd64 go build -o dist/simulador-linux-amd64 ./cmd/simulador

# Windows
GOOS=windows GOARCH=amd64 go build -o dist/assinatura-windows-amd64.exe ./cmd/assinatura
GOOS=windows GOARCH=amd64 go build -o dist/simulador-windows-amd64.exe ./cmd/simulador

# macOS
GOOS=darwin GOARCH=amd64 go build -o dist/assinatura-darwin-amd64 ./cmd/assinatura
GOOS=darwin GOARCH=amd64 go build -o dist/simulador-darwin-amd64 ./cmd/simulador
```

## Verificação de Integridade

Cada release inclui artefatos assinados com **Cosign**:

```bash
# Verificar integridade de um binário
cosign verify-blob \
  --certificate assinatura-v0.1.0-linux-amd64.pem \
  --signature assinatura-v0.1.0-linux-amd64.sig \
  assinatura-v0.1.0-linux-amd64
```

## Estrutura do Projeto

```
runner/
├── cmd/
│   ├── assinatura/       # CLI principal
│   └── simulador/        # CLI do simulador (Sprint 4)
├── internal/
│   ├── cli/              # Parsing Cobra
│   └── invoker/          # Invocação do assinador.jar
├── assinador/            # Projeto Java (Sprint 2)
├── .github/workflows/
│   ├── build.yml         # CI para pull requests e push
│   └── release.yml       # Publicação de releases
├── .goreleaser.yaml      # Configuração GoReleaser
└── go.mod
```

## CI/CD

- **build.yml**: Executa em `push` para `main` e em PRs
  - Testa em Windows, Linux, macOS
  - Gera artifacts de build
- **release.yml**: Executa ao criar tag `v*`
  - Compila para todas as plataformas
  - Gera checksums SHA256
  - Assina com Cosign
  - Publica no GitHub Releases

## Especificação

Consulte a [especificação completa](https://github.com/kyriosdata/runner/blob/25e12058b62da03f106653cd95dc509baa265341/SPECIFICATION.md) para:
- Histórias de usuário (US-01 a US-05)
- Critérios de aceitação
- Requisitos funcionais e não-funcionais
- Referências técnicas

## Contribuindo

Para contribuir:

1. Crie uma branch feature: `git checkout -b feat/minha-feature`
2. Commit com mensagens descritivas: `git commit -m "feat: descrição"`
3. Push para sua branch: `git push origin feat/minha-feature`
4. Abra um Pull Request

Ver [CONTRIBUTING.md](CONTRIBUTING.md) para mais detalhes.

## Status da Implementação

### Sprint 1 ✅
- [x] Estrutura base do CLI em Go
- [x] Comando `version`
- [x] Pipeline CI/CD multiplataforma
- [x] GoReleaser + Cosign

### Sprint 2
- [x] Assinador.jar com simulação
- [x] Invocação local (modo CLI)
- [x] Comandos `sign` e `validate`

### Sprint 3
- [x] Modo servidor HTTP
- [ ] Integração PKCS#11 real com SoftHSM2

### Sprint 4
- [x] CLI `simulador` com comandos `start`, `stop` e `status`
- [x] Registro de PID, porta e JAR em `~/.hubsaude/simulador.json`
- [x] Health check/readiness via HTTP
- [x] Download automático de `simulador.jar` via GitHub Releases ou `--source`
- [x] Cache local e verificação SHA-256
- [x] Build e release multiplataforma do CLI `simulador`

## Licença

Este projeto é distribuído sob a licença **MIT**. Ver [LICENSE](LICENSE).

## Contato

Para dúvidas ou sugestões, abra uma [issue](../../issues) no GitHub.
