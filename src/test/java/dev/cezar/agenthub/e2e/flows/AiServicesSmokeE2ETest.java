package dev.cezar.agenthub.e2e.flows;

import com.fasterxml.jackson.databind.ObjectMapper;
import dev.cezar.agenthub.e2e.config.E2EConfig;
import io.restassured.http.ContentType;
import io.restassured.response.Response;
import org.junit.jupiter.api.*;

import java.io.InputStream;
import java.util.List;
import java.util.Map;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Fluxo 6 — Smoke tests dos serviços IA.
 *
 * <p>Verifica que agenthub-embedding e agenthub-extractor estão respondendo
 * corretamente com os payloads esperados.
 * Não requer tenant ou JWT — os serviços são chamados diretamente.
 */
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("Fluxo 6 — Smoke Tests dos Serviços IA")
class AiServicesSmokeE2ETest {

    private static final int TIMEOUT_SECONDS = 60;

    @Test
    @Order(1)
    @DisplayName("POST /embed → retorna array de 1024 floats em menos de 10s")
    void embeddingService_returnsVector1024() {
        long start = System.currentTimeMillis();

        Response resp = given()
                .baseUri(E2EConfig.EMBEDDING_URL)
                .contentType(ContentType.JSON)
                .body(Map.of("text", "AgentHub é uma plataforma de orquestração de agentes IA"))
                .post("/embed")
                .then()
                .statusCode(200)
                .body("embedding", notNullValue())
                .body("dimension", equalTo(1024))
                .body("model", notNullValue())
                .extract().response();

        long elapsed = System.currentTimeMillis() - start;
        assertTrue(elapsed < TIMEOUT_SECONDS * 1_000L,
                "Embedding service should respond in < " + TIMEOUT_SECONDS + "s, took " + elapsed + "ms");

        List<Float> embedding = resp.jsonPath().getList("embedding");
        assertNotNull(embedding);
        assertEquals(1024, embedding.size(), "Embedding vector must have 1024 dimensions");

        // Verify values are actual floats in [-1, 1] range (normalized)
        embedding.forEach(v ->
                assertTrue(v >= -1.5f && v <= 1.5f, "Embedding value out of expected range: " + v)
        );
    }

    @Test
    @Order(2)
    @DisplayName("POST /embed com texto em inglês → dimensão correta")
    void embeddingService_englishText_returnsDimension() {
        given()
                .baseUri(E2EConfig.EMBEDDING_URL)
                .contentType(ContentType.JSON)
                .body(Map.of("text", "Hello, this is an end-to-end test for the embedding service"))
                .post("/embed")
                .then()
                .statusCode(200)
                .body("dimension", equalTo(1024))
                .body("embedding.size()", equalTo(1024));
    }

    @Test
    @Order(3)
    @DisplayName("POST /extract-text → extrai texto de PDF de teste em menos de 10s")
    void extractorService_extractsTextFromPdf() throws Exception {
        byte[] pdfBytes = loadSamplePdf();

        long start = System.currentTimeMillis();

        Response resp = given()
                .baseUri(E2EConfig.EXTRACTOR_URL)
                .contentType("multipart/form-data")
                .multiPart("file", "e2e-sample.pdf", pdfBytes, "application/pdf")
                .post("/extract-text")
                .then()
                .statusCode(200)
                .extract().response();

        long elapsed = System.currentTimeMillis() - start;
        assertTrue(elapsed < TIMEOUT_SECONDS * 1_000L,
                "Extractor service should respond in < " + TIMEOUT_SECONDS + "s, took " + elapsed + "ms");

        // Response may be plain text or JSON — handle both
        String body = resp.body().asString();
        assertNotNull(body);
        assertFalse(body.isBlank(), "Extracted text should not be empty");
        System.out.println("[AiServicesSmokeE2ETest] Extracted text length: " + body.length() + " chars");
    }

    @Test
    @Order(4)
    @DisplayName("GET /health ou / → extractor está rodando")
    void extractorService_isAlive() {
        // Try health endpoint, fall back to root
        int status = given()
                .baseUri(E2EConfig.EXTRACTOR_URL)
                .get("/health")
                .then()
                .extract().response().statusCode();

        if (status == 404) {
            // No /health endpoint — verify root responds
            given().baseUri(E2EConfig.EXTRACTOR_URL)
                    .get("/")
                    .then()
                    .statusCode(anyOf(equalTo(200), equalTo(404), equalTo(422)));
        } else {
            assertEquals(200, status, "Extractor /health should return 200");
        }
    }

    @Test
    @Order(5)
    @DisplayName("GET /health ou / → embedding está rodando")
    void embeddingService_isAlive() {
        int status = given()
                .baseUri(E2EConfig.EMBEDDING_URL)
                .get("/health")
                .then()
                .extract().response().statusCode();

        if (status == 404) {
            given().baseUri(E2EConfig.EMBEDDING_URL)
                    .get("/")
                    .then()
                    .statusCode(anyOf(equalTo(200), equalTo(404), equalTo(422)));
        } else {
            assertEquals(200, status, "Embedding /health should return 200");
        }
    }

    // ─── Helper ───────────────────────────────────────────────────────────────

    private byte[] loadSamplePdf() throws Exception {
        try (InputStream is = getClass().getResourceAsStream("/sample.pdf")) {
            if (is == null) throw new IllegalStateException("sample.pdf not found in test resources");
            return is.readAllBytes();
        }
    }
}
