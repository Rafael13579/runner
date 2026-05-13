package br.ufg.inf.assinador.controller;

import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;
import br.ufg.inf.assinador.service.SignatureService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api")
public class SignatureController {

    private final SignatureService service;

    public SignatureController(SignatureService service) {
        this.service = service;
    }

    @PostMapping("/sign")
    public ResponseEntity<SignatureResponse> sign(@RequestBody SignRequest request) {
        return ResponseEntity.ok(service.sign(request));
    }

    @PostMapping("/validate")
    public ResponseEntity<SignatureResponse> validate(@RequestBody ValidateRequest request) {
        return ResponseEntity.ok(service.validate(request));
    }
}