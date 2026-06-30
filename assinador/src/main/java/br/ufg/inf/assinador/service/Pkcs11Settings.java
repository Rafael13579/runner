package br.ufg.inf.assinador.service;

import org.springframework.core.env.Environment;

final class Pkcs11Settings {

    private final String configPath;
    private final String pin;
    private final String alias;
    private final String algorithm;

    private Pkcs11Settings(String configPath, String pin, String alias, String algorithm) {
        this.configPath = clean(configPath);
        this.pin = clean(pin);
        this.alias = clean(alias);
        this.algorithm = clean(algorithm);
    }

    static Pkcs11Settings fromEnvironment(Environment environment) {
        return new Pkcs11Settings(
                first(environment.getProperty("assinador.pkcs11.config-path"), getenv("PKCS11_CONFIG_PATH")),
                first(environment.getProperty("assinador.pkcs11.pin"), getenv("PKCS11_PIN")),
                first(environment.getProperty("assinador.pkcs11.alias"), getenv("PKCS11_ALIAS")),
                first(environment.getProperty("assinador.pkcs11.algorithm"), getenv("PKCS11_ALGORITHM"), "SHA256withRSA")
        );
    }

    static Pkcs11Settings fromSystemEnvironment() {
        return new Pkcs11Settings(
                first(System.getProperty("assinador.pkcs11.config-path"), getenv("PKCS11_CONFIG_PATH")),
                first(System.getProperty("assinador.pkcs11.pin"), getenv("PKCS11_PIN")),
                first(System.getProperty("assinador.pkcs11.alias"), getenv("PKCS11_ALIAS")),
                first(System.getProperty("assinador.pkcs11.algorithm"), getenv("PKCS11_ALGORITHM"), "SHA256withRSA")
        );
    }

    String configPath() {
        return configPath;
    }

    char[] pinChars() {
        return pin == null ? null : pin.toCharArray();
    }

    String alias() {
        return alias;
    }

    String algorithm() {
        return algorithm == null ? "SHA256withRSA" : algorithm;
    }

    String validationError() {
        if (configPath == null) {
            return "Configuracao PKCS#11 invalida: informe assinador.pkcs11.config-path ou PKCS11_CONFIG_PATH.";
        }
        if (pin == null) {
            return "Configuracao PKCS#11 invalida: informe assinador.pkcs11.pin ou PKCS11_PIN.";
        }
        return null;
    }

    private static String getenv(String name) {
        return System.getenv(name);
    }

    private static String first(String... values) {
        for (String value : values) {
            String cleaned = clean(value);
            if (cleaned != null) {
                return cleaned;
            }
        }
        return null;
    }

    private static String clean(String value) {
        if (value == null || value.trim().isEmpty()) {
            return null;
        }
        return value.trim();
    }
}
