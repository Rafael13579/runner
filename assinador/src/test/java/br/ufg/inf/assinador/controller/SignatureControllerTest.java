package br.ufg.inf.assinador.controller;

import org.junit.jupiter.api.Test;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.setup.MockMvcBuilders;

import br.ufg.inf.assinador.service.FakeSignatureService;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

class SignatureControllerTest {

    private final MockMvc mockMvc = MockMvcBuilders
            .standaloneSetup(new SignatureController(new FakeSignatureService()))
            .build();

    @Test
    void healthReturnsUp() throws Exception {
        mockMvc.perform(get("/api/health"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.status").value("UP"));
    }

    @Test
    void signReturnsSignatureForValidPayload() throws Exception {
        mockMvc.perform(post("/api/sign")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"content\":\"documento.pdf\",\"token\":\"token123\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.error").value(false))
                .andExpect(jsonPath("$.signature").value("MOCKED_SIGNATURE_BASE64_=="));
    }

    @Test
    void signReturnsBadRequestForMissingToken() throws Exception {
        mockMvc.perform(post("/api/sign")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"content\":\"documento.pdf\"}"))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.error").value(true))
                .andExpect(jsonPath("$.message").value(org.hamcrest.Matchers.containsString("--token")));
    }

    @Test
    void validateReturnsValidForKnownSignature() throws Exception {
        mockMvc.perform(post("/api/validate")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content("{\"content\":\"documento.pdf\",\"signature\":\"MOCKED_SIGNATURE_BASE64_==\"}"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.error").value(false))
                .andExpect(jsonPath("$.valid").value(true));
    }
}
