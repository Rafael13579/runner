package br.ufg.inf.assinador.controller;

import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;
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

    @GetMapping("/health")
    public ResponseEntity<Map<String, String>> health() {
        return ResponseEntity.ok(Map.of("status", "UP", "message", "Assinador pronto"));
    }

    @PostMapping("/sign")
    public ResponseEntity<SignatureResponse> sign(@RequestBody SignRequest request) {
        SignatureResponse response = service.sign(request);
        if (response.isError()) {
            return ResponseEntity.badRequest().body(response); // Retorna HTTP 400 se houver erro de validação (Critério E3)
        }
        return ResponseEntity.ok(response);
    }

    @PostMapping("/validate")
    public ResponseEntity<SignatureResponse> validate(@RequestBody ValidateRequest request) {
        SignatureResponse response = service.validate(request);
        if (response.isError()) {
            return ResponseEntity.badRequest().body(response); // Retorna HTTP 400 se houver erro de validação (Critério E3)
        }
        return ResponseEntity.ok(response);
    }
}