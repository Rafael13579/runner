package br.ufg.inf.assinador.service;


import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;

/**
 * Define operações de assinatura e validação digital.
 */
public interface SignatureService {

    /**
     * Cria uma assinatura digital simulada.
     */
    SignatureResponse sign(SignRequest request);

    /**
     * Valida uma assinatura digital simulada.
     */
    SignatureResponse validate(ValidateRequest request);
}