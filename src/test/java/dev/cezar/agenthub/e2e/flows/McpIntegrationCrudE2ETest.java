package dev.cezar.agenthub.e2e.flows;

import dev.cezar.agenthub.e2e.config.E2EConfig;
import io.restassured.http.ContentType;
import io.restassured.response.Response;
import org.junit.jupiter.api.*;

import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.Map;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;

@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("Fluxo — CRUD de integração MCP simplificada")
class McpIntegrationCrudE2ETest {

    private static final String TENANT_REALM = "test";

    private static String bearerToken;
    private static String integrationId;

    @BeforeAll
    static void setUp() {
        bearerToken = "Bearer " + fakeJwtForTenant(TENANT_REALM);
    }

    @AfterAll
    static void tearDown() {
        if (integrationId != null) {
            given()
                    .baseUri(E2EConfig.BACKEND_URL)
                    .header("Authorization", bearerToken)
                    .delete("/api/integrations/mcp/" + integrationId);
        }
    }

    @Test
    @Order(1)
    @DisplayName("POST /api/integrations/mcp cria integração MCP")
    void createMcpIntegration_returns201() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E MCP Integration",
                        "transportType", "stdio",
                        "command", "echo",
                        "args", new String[]{"hello"},
                        "env", Map.of(),
                        "autoStart", false,
                        "enabled", true
                ))
                .post("/api/integrations/mcp")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .body("name", equalTo("E2E MCP Integration"))
                .extract().response();

        integrationId = resp.jsonPath().getString("id");
    }

    @Test
    @Order(2)
    @DisplayName("GET /api/integrations/mcp/{id} retorna integração")
    void getMcpIntegration_returns200() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/integrations/mcp/" + integrationId)
                .then()
                .statusCode(200)
                .body("id", equalTo(integrationId))
                .body("transportType", equalTo("stdio"));
    }

    @Test
    @Order(3)
    @DisplayName("PATCH /api/integrations/mcp/{id} atualiza integração")
    void updateMcpIntegration_returns200() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "enabled", false
                ))
                .patch("/api/integrations/mcp/" + integrationId)
                .then()
                .statusCode(200)
                .body("id", equalTo(integrationId))
                .body("enabled", equalTo(false));
    }

    @Test
    @Order(4)
    @DisplayName("DELETE /api/integrations/mcp/{id} remove integração")
    void deleteMcpIntegration_returns204() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .delete("/api/integrations/mcp/" + integrationId)
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(204)));

        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/integrations/mcp/" + integrationId)
                .then()
                .statusCode(404);

        integrationId = null;
    }

    private static String fakeJwtForTenant(String tenantRealm) {
        String header = base64Url("{\"alg\":\"none\",\"typ\":\"JWT\"}");
        String payload = base64Url(("{" +
                "\"iss\":\"http://keycloak.cezar.dev/realms/" + tenantRealm + "\"," +
                "\"sub\":\"mcp-integration-e2e\"," +
                "\"preferred_username\":\"mcp-integration-e2e\"" +
                "}").replace("\n", ""));
        return header + "." + payload + ".signature";
    }

    private static String base64Url(String value) {
        return Base64.getUrlEncoder()
                .withoutPadding()
                .encodeToString(value.getBytes(StandardCharsets.UTF_8));
    }
}
