package br.ufg.inf.assinador.service;

import org.springframework.stereotype.Service;

@Service
public class FakeSignatureService implements SignatureService {

    private static final String PREFIX = "ASSINATURA_SIMULADA_HUB_SAUDE_";

    @Override
    public String sign(String message, String pin) {
        return PREFIX + Math.abs(message.hashCode());
    }

    @Override
    public boolean validate(String message, String signature,  String publicKey) {
        if (signature == null || message == null) {
            return false;
        }

        String expected = PREFIX + Math.abs(message.hashCode());

        return signature.equals(expected);
    }
}