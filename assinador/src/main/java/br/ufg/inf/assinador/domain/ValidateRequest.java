package br.ufg.inf.assinador.domain;

/**
 * Representa os dados para validação de assinatura digital.
 */
public class ValidateRequest {

    /**
     * Conteúdo original.
     */
    private String content;

    /**
     * Assinatura a ser validada.
     */
    private String signature;

    public ValidateRequest() {}

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getSignature() {
        return signature;
    }

    public void setSignature(String signature) {
        this.signature = signature;
    }
}