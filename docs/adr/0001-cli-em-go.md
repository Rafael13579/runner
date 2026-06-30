# ADR 0001: CLIs em Go

## Status

Aceita.

## Contexto

O Sistema Runner precisa entregar executáveis para Windows, Linux e macOS, com comandos simples para assinatura digital e gestão do Simulador do HubSaúde.

## Decisão

Os CLIs `assinatura` e `simulador` são implementados em Go, usando Cobra para parsing de comandos.

## Consequências

- Cross-compilation é feita diretamente no CI e no GoReleaser.
- A interface de usuário fica isolada da implementação Java dos JARs.
- A versão de runtime Java pode ser provisionada separadamente pelo CLI.
