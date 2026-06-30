# ADR 0004: Serviço PKCS#11 opcional

## Status

Aceita.

## Contexto

O Runner precisa manter o fluxo simulado simples para desenvolvimento, mas também oferecer um caminho real para assinatura com token USB, smart card, HSM ou SoftHSM2.

## Decisão

O `assinador.jar` possui dois serviços de assinatura:

- `FakeSignatureService`, ativo por padrão.
- `Pkcs11SignatureService`, ativo quando `assinador.mode=pkcs11`.

O serviço PKCS#11 usa o provider `SunPKCS11`, carrega um `KeyStore` do tipo `PKCS11`, assina com a chave privada do token e valida usando o certificado associado.

## Consequências

- A suíte padrão continua reprodutível sem hardware criptográfico.
- Integrações reais dependem de biblioteca PKCS#11, PIN e alias configurados.
- Testes com SoftHSM2 são opcionais e só rodam quando o ambiente fornece as variáveis necessárias.
