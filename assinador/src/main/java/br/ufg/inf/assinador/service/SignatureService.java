package br.ufg.inf.assinador.service;

public interface SignatureService {
    String sign(String message, String pin);
    boolean validate(String message, String signature,  String publicKey);
}