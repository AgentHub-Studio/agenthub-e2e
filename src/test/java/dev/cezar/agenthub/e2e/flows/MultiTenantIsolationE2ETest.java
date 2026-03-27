package dev.cezar.agenthub.e2e.flows;

import dev.cezar.agenthub.e2e.config.E2EConfig;
import dev.cezar.agenthub.e2e.support.TenantFixture;
import io.restassured.http.ContentType;
import io.restassured.response.Response;
import org.junit.jupiter.api.*;

import java.util.List;
import java.util.Map;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Fluxo 5 — Isolamento multi-tenant.
 *
 * <p>Verifica que recursos criados em um tenant não são visíveis nem acessíveis
 * por outro tenant, mesmo com um JWT válido.
 */
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("Fluxo 5 — Isolamento Multi-Tenant")
class MultiTenantIsolationE2ETest {

    private static TenantFixture tenantA;
    private static TenantFixture tenantB;
    private static String agentIdInTenantA;

    @BeforeAll
    static void setUp() {
        // Provision two separate tenants in parallel would be nice, but sequential is safer
        tenantA = TenantFixture.create();
        tenantB = TenantFixture.create();
    }

    @AfterAll
    static void tearDown() {
        try {
            if (agentIdInTenantA != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", tenantA.bearerToken())
                        .delete("/api/agents/" + agentIdInTenantA);
            }
        } catch (Exception ignored) {}
        try { if (tenantA != null) tenantA.destroy(); } catch (Exception ignored) {}
        try { if (tenantB != null) tenantB.destroy(); } catch (Exception ignored) {}
    }

    @Test
    @Order(1)
    @DisplayName("Tenant A cria um agent exclusivo")
    void tenantA_createsAgent() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantA.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "agent-exclusivo-a",
                        "description", "Belongs only to tenant A",
                        "enabled", true,
                        "model", Map.of("provider", "ollama", "name", "llama3.2")
                ))
                .post("/api/agents")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)))
                .body("name", equalTo("agent-exclusivo-a"))
                .extract().response();

        agentIdInTenantA = resp.jsonPath().getString("id");
        assertNotNull(agentIdInTenantA);
    }

    @Test
    @Order(2)
    @DisplayName("Tenant A vê o próprio agent na listagem")
    void tenantA_seesOwnAgent() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantA.bearerToken())
                .get("/api/agents")
                .then()
                .statusCode(200)
                .body("content.find { it.id == '" + agentIdInTenantA + "' }", notNullValue());
    }

    @Test
    @Order(3)
    @DisplayName("Tenant B NÃO vê o agent do tenant A na listagem")
    void tenantB_doesNotSeeAgentFromTenantA() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantB.bearerToken())
                .get("/api/agents")
                .then()
                .statusCode(200)
                .extract().response();

        List<Map<String, Object>> content = resp.jsonPath().getList("content");
        assertNotNull(content);

        boolean found = content.stream()
                .anyMatch(a -> agentIdInTenantA.equals(a.get("id")));
        assertFalse(found, "Tenant B should NOT see tenant A's agent in the listing");
    }

    @Test
    @Order(4)
    @DisplayName("Tenant B NÃO consegue acessar diretamente o agent do tenant A → 404")
    void tenantB_cannotAccessAgentFromTenantA() {
        // Backend should return 404 (schema isolation) rather than 403
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantB.bearerToken())
                .get("/api/agents/" + agentIdInTenantA)
                .then()
                .statusCode(anyOf(equalTo(404), equalTo(403)));
        // 404 expected because the agent simply doesn't exist in tenant B's schema
    }

    @Test
    @Order(5)
    @DisplayName("Tenant B tem schema isolado — listagem retorna apenas recursos do próprio tenant")
    void tenantB_hasIsolatedSchema() {
        // Tenant B creates its own agent
        String bAgentId = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantB.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "agent-exclusivo-b",
                        "description", "Belongs only to tenant B",
                        "enabled", true,
                        "model", Map.of("provider", "ollama", "name", "llama3.2")
                ))
                .post("/api/agents")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)))
                .extract().jsonPath().getString("id");

        // Verify B sees its agent
        given().baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantB.bearerToken())
                .get("/api/agents/" + bAgentId)
                .then()
                .statusCode(200)
                .body("name", equalTo("agent-exclusivo-b"));

        // Verify A does NOT see B's agent
        given().baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantA.bearerToken())
                .get("/api/agents/" + bAgentId)
                .then()
                .statusCode(anyOf(equalTo(404), equalTo(403)));

        // Cleanup B's agent
        given().baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenantB.bearerToken())
                .delete("/api/agents/" + bAgentId);
    }

    @Test
    @Order(6)
    @DisplayName("Request sem JWT → 401 em endpoint protegido")
    void noJwt_returns401() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .get("/api/agents")
                .then()
                .statusCode(401);
    }
}
