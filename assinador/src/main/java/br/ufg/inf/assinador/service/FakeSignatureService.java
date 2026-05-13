package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;
import org.springframework.stereotype.Service;

@Service
public class FakeSignatureService implements SignatureService {

    private static final String FAKE_SIGNATURE = "MOCKED_SIGNATURE_BASE64_==";

    @Override
    public SignatureResponse sign(SignRequest request) {
        if (request.getContent() == null) {
            return new SignatureResponse(null, false, "Content inválido");
        }

        return new SignatureResponse(FAKE_SIGNATURE, true, "OK");
    }

    @Override
    public SignatureResponse validate(ValidateRequest request) {
        boolean valid = FAKE_SIGNATURE.equals(request.getSignature());

        return new SignatureResponse(
                request.getSignature(),
                valid,
                valid ? "Válida" : "Inválida"
        );
    }
}