# ADR 0003: Obtenção dinâmica do simulador.jar

## Status

Aceita.

## Contexto

O CLI `simulador` deve iniciar o Simulador do HubSaúde sem exigir que o usuário conheça comandos Java ou faça download manual do JAR.

## Decisão

O CLI procura o JAR nesta ordem:

1. `~/.hubsaude/simulador/simulador.jar`
2. `simulador.jar` no diretório atual
3. asset `simulador.jar` do último GitHub Release
4. URL alternativa informada por `--source`

Quando um checksum é informado por `--checksum`, ou quando existe o asset `simulador.jar.sha256`, o SHA-256 é validado antes da execução.

## Consequências

- O download é cacheado e não é repetido quando o JAR já existe.
- Releases precisam publicar `simulador.jar` e, idealmente, `simulador.jar.sha256`.
