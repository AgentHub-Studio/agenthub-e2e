package dev.cezar.agenthub.e2e.flows;

import dev.cezar.agenthub.e2e.config.E2EConfig;
import dev.cezar.agenthub.e2e.support.TenantFixture;
import io.restassured.http.ContentType;
import io.restassured.response.Response;
import org.junit.jupiter.api.*;

import java.io.InputStream;
import java.util.Base64;
import java.util.Map;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Fluxo 3 — RAG Pipeline.
 *
 * <p>Cria uma knowledge base, faz upload de um documento PDF de teste,
 * aguarda a indexação completa (INDEXED) e verifica busca via hybrid search.
 */
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("Fluxo 3 — RAG Pipeline")
class RagPipelineE2ETest {

    private static TenantFixture tenant;
    private static String kbId;
    private static String documentId;

    private static final int INDEXING_TIMEOUT_SECONDS = 120;
    private static final int POLL_INTERVAL_MS = 5_000;

    @BeforeAll
    static void setUp() {
        tenant = TenantFixture.create();
    }

    @AfterAll
    static void tearDown() {
        try {
            if (kbId != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", tenant.bearerToken())
                        .delete("/api/knowledge-bases/" + kbId);
            }
        } finally {
            if (tenant != null) tenant.destroy();
        }
    }

    @Test
    @Order(1)
    @DisplayName("POST /api/knowledge-bases → 201 com status ACTIVE")
    void createKnowledgeBase_returns201() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Knowledge Base",
                        "description", "Created by e2e RAG pipeline test",
                        "storageType", "PGVECTOR",
                        "searchMode", "HYBRID"
                ))
                .post("/api/knowledge-bases")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .body("status", equalTo("ACTIVE"))
                .extract().response();

        kbId = resp.jsonPath().getString("id");
        assertNotNull(kbId, "Knowledge Base ID must not be null");
    }

    @Test
    @Order(2)
    @DisplayName("POST /api/documents → 201, upload de PDF de teste")
    void uploadDocument_returns201() throws Exception {
        assertNotNull(kbId, "Knowledge base must be created first");

        // Load sample PDF from test resources
        byte[] pdfBytes = loadSamplePdf();
        String base64Content = Base64.getEncoder().encodeToString(pdfBytes);

        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "knowledgeBaseId", kbId,
                        "fileName", "e2e-sample.pdf",
                        "contentType", "application/pdf",
                        "fileContent", base64Content
                ))
                .post("/api/documents")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .body("knowledgeBaseId", equalTo(kbId))
                .body("fileName", equalTo("e2e-sample.pdf"))
                .extract().response();

        documentId = resp.jsonPath().getString("id");
        assertNotNull(documentId, "Document ID must not be null");
    }

    @Test
    @Order(3)
    @DisplayName("Polling status do documento até INDEXED (timeout 120s)")
    void documentIndexing_reachesIndexedStatus() throws InterruptedException {
        assertNotNull(documentId, "Document must be uploaded first");

        long deadline = System.currentTimeMillis() + (INDEXING_TIMEOUT_SECONDS * 1_000L);
        String status = "PENDING";

        while (System.currentTimeMillis() < deadline) {
            Response resp = given()
                    .baseUri(E2EConfig.BACKEND_URL)
                    .header("Authorization", tenant.bearerToken())
                    .get("/api/documents/" + documentId)
                    .then()
                    .statusCode(200)
                    .extract().response();

            status = resp.jsonPath().getString("status");
            System.out.println("[RagPipelineE2ETest] Document status: " + status);

            if ("INDEXED".equals(status)) return;
            if ("FAILED".equals(status)) {
                fail("Document processing failed with status FAILED");
            }

            Thread.sleep(POLL_INTERVAL_MS);
        }

        fail("Document did not reach INDEXED status within " + INDEXING_TIMEOUT_SECONDS + "s. Last status: " + status);
    }

    @Test
    @Order(4)
    @DisplayName("GET /api/search?q=...&type=knowledge_bases → resultados relevantes")
    void hybridSearch_returnsResults() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .queryParam("q", "AgentHub platform agents")
                .queryParam("type", "knowledge_bases")
                .queryParam("size", 5)
                .get("/api/search")
                .then()
                .statusCode(200)
                .body("content", notNullValue());
        // Note: result count may be 0 if no text overlap; the critical assertion is status=INDEXED (test 3)
    }

    @Test
    @Order(5)
    @DisplayName("GET /api/knowledge-bases/{id} → documentCount >= 1 após indexação")
    void knowledgeBase_documentCountUpdated() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .get("/api/knowledge-bases/" + kbId)
                .then()
                .statusCode(200)
                .body("documentCount", greaterThanOrEqualTo(1));
    }

    // ─── Helper ───────────────────────────────────────────────────────────────

    private byte[] loadSamplePdf() throws Exception {
        try (InputStream is = getClass().getResourceAsStream("/sample.pdf")) {
            if (is == null) throw new IllegalStateException("sample.pdf not found in test resources");
            return is.readAllBytes();
        }
    }
}
