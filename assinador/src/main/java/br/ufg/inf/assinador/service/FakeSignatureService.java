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
        if (request == null || request.getContent() == null || request.getContent().trim().isEmpty()) {
            return new SignatureResponse(null, false, "Erro do Usuário: O parâmetro '--content' é obrigatório e não pode estar vazio.", true);
        }
        if (request.getToken() == null || request.getToken().trim().isEmpty()) {
            return new SignatureResponse(null, false, "Erro do Usuário: O parâmetro '--token' é obrigatório.", true);
        }

        return new SignatureResponse(FAKE_SIGNATURE, true, "Assinatura gerada com sucesso", false);
    }

    @Override
    public SignatureResponse validate(ValidateRequest request) {
        // Validação estrita de parâmetros da US-02.3
        if (request == null || request.getContent() == null || request.getContent().trim().isEmpty()) {
            return new SignatureResponse(null, false, "Erro do Usuário: O parâmetro '--content' original é obrigatório para validação.", true);
        }
        if (request.getSignature() == null || request.getSignature().trim().isEmpty()) {
            return new SignatureResponse(null, false, "Erro do Usuário: O parâmetro '--signature' é obrigatório para validação.", true);
        }

        boolean isValid = FAKE_SIGNATURE.equals(request.getSignature());

        return new SignatureResponse(
                request.getSignature(),
                isValid,
                isValid ? "Assinatura válida" : "Assinatura inválida: divergência de hash",
                false
        );
    }
}