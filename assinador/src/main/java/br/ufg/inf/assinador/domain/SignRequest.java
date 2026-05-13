package br.ufg.inf.assinador.domain;

/**
 * Representa os dados para criação de assinatura digital.
 */
public class SignRequest {

    /**
     * Conteúdo a ser assinado.
     */
    private String content;

    /**
     * Token, PIN ou credencial.
     */
    private String token;

    public SignRequest() {}

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }
}