# PKCS#11 e SoftHSM2

O modo PKCS#11 permite que o `assinador.jar` use uma chave privada mantida em token USB, smart card, HSM ou simulador compatível. O modo padrão continua sendo `fake`; para usar PKCS#11, configure explicitamente `assinador.mode=pkcs11`.

## Propriedades

As opções podem ser passadas como propriedades Spring/JVM ou variáveis de ambiente:

| Propriedade | Variável | Descrição |
|-------------|----------|-----------|
| `assinador.mode` | `ASSINADOR_MODE` | Use `pkcs11` para ativar o serviço real |
| `assinador.pkcs11.config-path` | `PKCS11_CONFIG_PATH` | Caminho do arquivo de configuração SunPKCS11 |
| `assinador.pkcs11.pin` | `PKCS11_PIN` | PIN do token |
| `assinador.pkcs11.alias` | `PKCS11_ALIAS` | Alias da chave/certificado; opcional quando há apenas uma chave |
| `assinador.pkcs11.algorithm` | `PKCS11_ALGORITHM` | Algoritmo de assinatura; padrão `SHA256withRSA` |

## Arquivo SunPKCS11

Exemplo para Windows com SoftHSM2:

```text
name=HubSaudeSoftHSM
library=C:\Program Files\SoftHSM2\lib\softhsm2-x64.dll
slotListIndex=0
```

Em Linux, o caminho normalmente é semelhante a:

```text
name=HubSaudeSoftHSM
library=/usr/lib/softhsm/libsofthsm2.so
slotListIndex=0
```

## Execução em modo servidor

```bash
java -Dserver.port=8080 \
  -Dassinador.mode=pkcs11 \
  -Dassinador.pkcs11.config-path=./pkcs11.cfg \
  -Dassinador.pkcs11.pin=1234 \
  -Dassinador.pkcs11.alias=hubsaude \
  -jar assinador.jar
```

Também é possível usar variáveis de ambiente:

```bash
export ASSINADOR_MODE=pkcs11
export PKCS11_CONFIG_PATH=./pkcs11.cfg
export PKCS11_PIN=1234
export PKCS11_ALIAS=hubsaude
java -jar assinador.jar
```

## Execução local

```bash
ASSINADOR_MODE=pkcs11 \
PKCS11_CONFIG_PATH=./pkcs11.cfg \
PKCS11_PIN=1234 \
PKCS11_ALIAS=hubsaude \
java -jar assinador.jar sign --content documento.pdf
```

## Teste de integração opcional

Os testes normais não exigem token. O teste de integração PKCS#11 só roda quando `PKCS11_CONFIG_PATH` e `PKCS11_PIN` estão definidos:

```bash
cd assinador
PKCS11_CONFIG_PATH=../pkcs11.cfg PKCS11_PIN=1234 PKCS11_ALIAS=hubsaude ./mvnw test
```

Se essas variáveis não existirem, o teste é ignorado.
