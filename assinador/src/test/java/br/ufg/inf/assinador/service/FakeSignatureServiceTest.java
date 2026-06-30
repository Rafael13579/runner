package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class FakeSignatureServiceTest {

    private final FakeSignatureService service = new FakeSignatureService();

    @Test
    void signReturnsFakeSignatureForValidRequest() {
        SignRequest request = new SignRequest();
        request.setContent("documento.pdf");
        request.setToken("token123");

        SignatureResponse response = service.sign(request);

        assertFalse(response.isError());
        assertTrue(response.isValid());
        assertEquals("MOCKED_SIGNATURE_BASE64_==", response.getSignature());
    }

    @Test
    void signRejectsMissingContent() {
        SignRequest request = new SignRequest();
        request.setToken("token123");

        SignatureResponse response = service.sign(request);

        assertTrue(response.isError());
        assertFalse(response.isValid());
        assertTrue(response.getMessage().contains("--content"));
    }

    @Test
    void signRejectsMissingToken() {
        SignRequest request = new SignRequest();
        request.setContent("documento.pdf");

        SignatureResponse response = service.sign(request);

        assertTrue(response.isError());
        assertFalse(response.isValid());
        assertTrue(response.getMessage().contains("--token"));
    }

    @Test
    void validateAcceptsKnownFakeSignature() {
        ValidateRequest request = new ValidateRequest();
        request.setContent("documento.pdf");
        request.setSignature("MOCKED_SIGNATURE_BASE64_==");

        SignatureResponse response = service.validate(request);

        assertFalse(response.isError());
        assertTrue(response.isValid());
    }

    @Test
    void validateRejectsDifferentSignatureWithoutSystemError() {
        ValidateRequest request = new ValidateRequest();
        request.setContent("documento.pdf");
        request.setSignature("outra-assinatura");

        SignatureResponse response = service.validate(request);

        assertFalse(response.isError());
        assertFalse(response.isValid());
        assertTrue(response.getMessage().contains("inv"));
    }

    @Test
    void validateRejectsMissingSignature() {
        ValidateRequest request = new ValidateRequest();
        request.setContent("documento.pdf");

        SignatureResponse response = service.validate(request);

        assertTrue(response.isError());
        assertFalse(response.isValid());
        assertTrue(response.getMessage().contains("--signature"));
    }
}
