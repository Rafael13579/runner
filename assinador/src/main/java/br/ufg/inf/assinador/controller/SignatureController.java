package br.ufg.inf.assinador.controller;

import br.ufg.inf.assinador.dtos.SignRequest;
import br.ufg.inf.assinador.dtos.ValidateRequest;
import br.ufg.inf.assinador.service.SignatureService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/api")
public class SignatureController {

    private final SignatureService service;

    public SignatureController(SignatureService service) {
        this.service = service;
    }

    @PostMapping("/sign")
    public ResponseEntity<?> sign(@RequestBody SignRequest request) {
        
        if (request.message == null || request.message.isBlank()) {
            return ResponseEntity.badRequest()
                    .body(Map.of("error", "Mensagem é obrigatória"));
        }

        if (request.pin == null || request.pin.isBlank()) {
            return ResponseEntity.badRequest()
                    .body(Map.of("error", "PIN é obrigatório"));
        }

        String signature = service.sign(request.message, request.pin);

        return ResponseEntity.ok(Map.of(
                "signature", signature,
                "status", "SUCCESS"
        ));
    }

    @PostMapping("/validate")
    public ResponseEntity<?> validate(@RequestBody ValidateRequest request) {

        if (request.message == null || request.message.isBlank()) {
            return ResponseEntity.badRequest()
                    .body(Map.of("error", "Mensagem é obrigatória"));
        }

        if (request.signature == null || request.signature.isBlank()) {
            return ResponseEntity.badRequest()
                    .body(Map.of("error", "Assinatura é obrigatória"));
        }

        boolean valid = service.validate(request.message, request.signature, request.publicKey);

        return ResponseEntity.ok(Map.of(
                "valid", valid
        ));
    }
}