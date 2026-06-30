package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

class Pkcs11SignatureServiceTest {

    @Test
    void signFailsClearlyWhenConfigurationIsMissing() {
        Pkcs11SignatureService service = new Pkcs11SignatureService();
        SignRequest request = new SignRequest();
        request.setContent("documento.pdf");

        SignatureResponse response = service.sign(request);

        assertTrue(response.isError());
        assertFalse(response.isValid());
        assertTrue(response.getMessage().contains("PKCS#11"));
    }

    @Test
    void integrationSignsAndValidatesWhenTokenIsConfigured() {
        assumeTrue(hasPkcs11Environment(), "PKCS#11 environment not configured");

        Pkcs11SignatureService service = new Pkcs11SignatureService();
        SignRequest signRequest = new SignRequest();
        signRequest.setContent("documento.pdf");

        SignatureResponse signResponse = service.sign(signRequest);

        assertFalse(signResponse.isError(), signResponse.getMessage());
        assertNotNull(signResponse.getSignature());

        ValidateRequest validateRequest = new ValidateRequest();
        validateRequest.setContent("documento.pdf");
        validateRequest.setSignature(signResponse.getSignature());

        SignatureResponse validateResponse = service.validate(validateRequest);

        assertFalse(validateResponse.isError(), validateResponse.getMessage());
        assertTrue(validateResponse.isValid());
    }

    private boolean hasPkcs11Environment() {
        return exists("PKCS11_CONFIG_PATH") && exists("PKCS11_PIN");
    }

    private boolean exists(String name) {
        String value = System.getenv(name);
        return value != null && !value.trim().isEmpty();
    }
}
