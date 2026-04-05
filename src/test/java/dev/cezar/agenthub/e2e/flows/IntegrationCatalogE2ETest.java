package dev.cezar.agenthub.e2e.flows;

import dev.cezar.agenthub.e2e.config.E2EConfig;
import io.restassured.http.ContentType;
import io.restassured.response.Response;
import org.junit.jupiter.api.*;

import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.Map;
import java.util.UUID;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertNotNull;

/**
 * Fluxo — Catálogo unificado de integrações.
 *
 * <p>Cria recursos legados reais (HTTP tool, datasource + database tool, MCP config)
 * e valida que o backend os expõe de forma unificada via GET /api/integrations.
 */
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("Fluxo — Catálogo Unificado de Integrações")
class IntegrationCatalogE2ETest {

    private static final String TENANT_REALM = "test";

    private static String bearerToken;
    private static String httpToolId;
    private static String databaseToolId;
    private static String dataSourceId;
    private static String mcpId;

    @BeforeAll
    static void setUp() {
        bearerToken = "Bearer " + fakeJwtForTenant(TENANT_REALM);
    }

    @AfterAll
    static void tearDown() {
        try {
            if (httpToolId != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", bearerToken)
                        .delete("/api/tools/" + httpToolId);
            }
            if (databaseToolId != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", bearerToken)
                        .delete("/api/tools/" + databaseToolId);
            }
            if (mcpId != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", bearerToken)
                        .delete("/api/mcp-server-configs/" + mcpId);
            }
            if (dataSourceId != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", bearerToken)
                        .delete("/api/datasources/" + dataSourceId);
            }
        } finally {
            // No-op.
        }
    }

    @Test
    @Order(1)
    @DisplayName("Criar datasource, tools e MCP legados")
    void createLegacyResources() {
        Response dsResp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Orders DB",
                        "type", "POSTGRESQL",
                        "host", "pg.internal",
                        "port", 5432,
                        "database", "orders",
                        "dbUser", "orders_user",
                        "dbPassword", "orders_secret"
                ))
                .post("/api/datasources")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .extract().response();
        dataSourceId = dsResp.jsonPath().getString("id");

        Response httpResp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E ERP API",
                        "description", "HTTP integration for ERP sync",
                        "type", "HTTP",
                        "config", Map.of(
                                "url", "https://erp.example.com/customers",
                                "method", "POST"
                        )
                ))
                .post("/api/tools")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)))
                .body("id", notNullValue())
                .extract().response();
        httpToolId = httpResp.jsonPath().getString("id");

        Response dbToolResp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Orders Query",
                        "description", "Database integration for orders",
                        "type", "DATABASE",
                        "config", Map.of(
                                "dataSourceId", dataSourceId,
                                "sql", "select * from orders limit 10"
                        )
                ))
                .post("/api/tools")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)))
                .body("id", notNullValue())
                .extract().response();
        databaseToolId = dbToolResp.jsonPath().getString("id");

        Response mcpResp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "e2e-filesystem",
                        "transportType", "stdio",
                        "command", "npx",
                        "args", new String[]{"-y", "@modelcontextprotocol/server-filesystem", "/tmp"},
                        "autoStart", false,
                        "enabled", true
                ))
                .post("/api/mcp-server-configs")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .extract().response();
        mcpId = mcpResp.jsonPath().getString("id");

        assertNotNull(dataSourceId);
        assertNotNull(httpToolId);
        assertNotNull(databaseToolId);
        assertNotNull(mcpId);
    }

    @Test
    @Order(2)
    @DisplayName("GET /api/integrations retorna catálogo unificado paginado")
    void listIntegrations_returnsUnifiedCatalog() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .queryParam("page", 0)
                .queryParam("size", 100)
                .get("/api/integrations")
                .then()
                .statusCode(200)
                .body("content", notNullValue())
                .body("totalElements", greaterThanOrEqualTo(3))
                .body("page", equalTo(0))
                .body("size", equalTo(100))
                .body("content.find { it.id == '" + httpToolId + "' }.type", equalTo("HTTP_API"))
                .body("content.find { it.id == '" + httpToolId + "' }.sourceKind", equalTo("tool"))
                .body("content.find { it.id == '" + httpToolId + "' }.advanced", equalTo(false))
                .body("content.find { it.id == '" + dataSourceId + "' }.type", equalTo("DATABASE_QUERY"))
                .body("content.find { it.id == '" + dataSourceId + "' }.sourceKind", equalTo("datasource"))
                .body("content.find { it.id == '" + dataSourceId + "' }.advanced", equalTo(false))
                .body("content.find { it.id == '" + mcpId + "' }.type", equalTo("MCP"))
                .body("content.find { it.id == '" + mcpId + "' }.sourceKind", equalTo("mcp_server_config"));
    }

    @Test
    @Order(3)
    @DisplayName("GET /api/integrations?type=MCP filtra corretamente")
    void listIntegrations_filtersByType() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .queryParam("type", "MCP")
                .get("/api/integrations")
                .then()
                .statusCode(200)
                .body("content.size()", greaterThanOrEqualTo(1))
                .body("content.type.flatten().unique()", contains("MCP"))
                .body("content.find { it.id == '" + mcpId + "' }", notNullValue());
    }

    @Test
    @Order(4)
    @DisplayName("GET /api/integrations com filtro inválido retorna 400")
    void listIntegrations_invalidFilter_returns400() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .queryParam("type", "INVALID")
                .get("/api/integrations")
                .then()
                .statusCode(400)
                .body("error", containsString("invalid type filter"));
    }

    private static String fakeJwtForTenant(String tenantRealm) {
        String header = base64Url("{\"alg\":\"none\",\"typ\":\"JWT\"}");
        String payload = base64Url(("{" +
                "\"iss\":\"http://keycloak.cezar.dev/realms/" + tenantRealm + "\"," +
                "\"sub\":\"integration-catalog-e2e\"," +
                "\"preferred_username\":\"integration-e2e\"" +
                "}").replace("\n", ""));
        return header + "." + payload + ".signature";
    }

    private static String base64Url(String value) {
        return Base64.getUrlEncoder()
                .withoutPadding()
                .encodeToString(value.getBytes(StandardCharsets.UTF_8));
    }
}
