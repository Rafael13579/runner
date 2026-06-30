# ADR 0002: Registro de processos em ~/.hubsaude

## Status

Aceita.

## Contexto

Os comandos precisam iniciar, reutilizar, consultar e encerrar processos Java executados em background.

## Decisão

O Runner grava metadados operacionais em `~/.hubsaude/`.

- O assinador usa `~/.hubsaude/assinador.pid`.
- O simulador usa `~/.hubsaude/simulador.json`, contendo PID, porta, caminho do JAR e data de início.

## Consequências

- O estado é simples de inspecionar e remover manualmente quando necessário.
- O health check HTTP continua sendo a fonte de verdade para readiness; o arquivo local apenas aponta o processo esperado.
