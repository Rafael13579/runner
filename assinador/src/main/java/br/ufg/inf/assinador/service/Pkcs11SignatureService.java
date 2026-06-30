package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.domain.SignRequest;
import br.ufg.inf.assinador.domain.SignatureResponse;
import br.ufg.inf.assinador.domain.ValidateRequest;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.core.env.Environment;
import org.springframework.stereotype.Service;

import java.nio.charset.StandardCharsets;
import java.security.Key;
import java.security.KeyStore;
import java.security.PrivateKey;
import java.security.Provider;
import java.security.Security;
import java.security.Signature;
import java.security.cert.Certificate;
import java.util.Base64;
import java.util.Enumeration;

@Service
@ConditionalOnProperty(name = "assinador.mode", havingValue = "pkcs11")
public class Pkcs11SignatureService implements SignatureService {

    private final Pkcs11Settings settings;

    public Pkcs11SignatureService(Environment environment) {
        this(Pkcs11Settings.fromEnvironment(environment));
    }

    Pkcs11SignatureService(Pkcs11Settings settings) {
        this.settings = settings;
    }

    public Pkcs11SignatureService() {
        this(Pkcs11Settings.fromSystemEnvironment());
    }

    @Override
    public SignatureResponse sign(SignRequest request) {
        if (request == null || request.getContent() == null || request.getContent().trim().isEmpty()) {
            return new SignatureResponse(null, false, "Erro do Usuario: o parametro '--content' e obrigatorio.", true);
        }

        String settingsError = settings.validationError();
        if (settingsError != null) {
            return new SignatureResponse(null, false, settingsError, true);
        }

        try {
            KeyMaterial material = loadKeyMaterial();
            Signature signer = Signature.getInstance(settings.algorithm(), material.provider());
            signer.initSign(material.privateKey());
            signer.update(request.getContent().getBytes(StandardCharsets.UTF_8));
            String signature = Base64.getEncoder().encodeToString(signer.sign());
            return new SignatureResponse(signature, true, "Assinatura PKCS#11 gerada com sucesso", false);
        } catch (Exception e) {
            return new SignatureResponse(null, false, "Erro de sistema ao assinar via PKCS#11: " + e.getMessage(), true);
        }
    }

    @Override
    public SignatureResponse validate(ValidateRequest request) {
        if (request == null || request.getContent() == null || request.getContent().trim().isEmpty()) {
            return new SignatureResponse(null, false, "Erro do Usuario: o parametro '--content' e obrigatorio.", true);
        }
        if (request.getSignature() == null || request.getSignature().trim().isEmpty()) {
            return new SignatureResponse(null, false, "Erro do Usuario: o parametro '--signature' e obrigatorio.", true);
        }

        String settingsError = settings.validationError();
        if (settingsError != null) {
            return new SignatureResponse(null, false, settingsError, true);
        }

        try {
            KeyMaterial material = loadKeyMaterial();
            Signature verifier = Signature.getInstance(settings.algorithm());
            verifier.initVerify(material.certificate());
            verifier.update(request.getContent().getBytes(StandardCharsets.UTF_8));
            boolean valid = verifier.verify(Base64.getDecoder().decode(request.getSignature()));
            return new SignatureResponse(request.getSignature(), valid,
                    valid ? "Assinatura PKCS#11 valida" : "Assinatura PKCS#11 invalida",
                    false);
        } catch (IllegalArgumentException e) {
            return new SignatureResponse(request.getSignature(), false, "Erro do Usuario: assinatura nao esta em Base64 valido.", true);
        } catch (Exception e) {
            return new SignatureResponse(request.getSignature(), false, "Erro de sistema ao validar via PKCS#11: " + e.getMessage(), true);
        }
    }

    private KeyMaterial loadKeyMaterial() throws Exception {
        Provider baseProvider = Security.getProvider("SunPKCS11");
        if (baseProvider == null) {
            throw new IllegalStateException("Provider SunPKCS11 nao esta disponivel nesta JVM.");
        }

        Provider provider = baseProvider.configure(settings.configPath());
        if (Security.getProvider(provider.getName()) == null) {
            Security.addProvider(provider);
        }

        KeyStore keyStore = KeyStore.getInstance("PKCS11", provider);
        keyStore.load(null, settings.pinChars());

        String alias = resolveAlias(keyStore);
        Key key = keyStore.getKey(alias, null);
        if (!(key instanceof PrivateKey privateKey)) {
            throw new IllegalStateException("Alias PKCS#11 nao possui chave privada: " + alias);
        }

        Certificate certificate = keyStore.getCertificate(alias);
        if (certificate == null) {
            throw new IllegalStateException("Alias PKCS#11 nao possui certificado: " + alias);
        }

        return new KeyMaterial(provider, privateKey, certificate);
    }

    private String resolveAlias(KeyStore keyStore) throws Exception {
        if (settings.alias() != null) {
            if (!keyStore.containsAlias(settings.alias())) {
                throw new IllegalStateException("Alias PKCS#11 nao encontrado: " + settings.alias());
            }
            return settings.alias();
        }

        Enumeration<String> aliases = keyStore.aliases();
        while (aliases.hasMoreElements()) {
            String candidate = aliases.nextElement();
            if (keyStore.isKeyEntry(candidate)) {
                return candidate;
            }
        }
        throw new IllegalStateException("Nenhuma chave privada encontrada no token PKCS#11.");
    }

    private record KeyMaterial(Provider provider, PrivateKey privateKey, Certificate certificate) {
    }
}
