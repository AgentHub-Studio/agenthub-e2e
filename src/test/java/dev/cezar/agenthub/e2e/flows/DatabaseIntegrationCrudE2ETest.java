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
@DisplayName("Fluxo — CRUD de integração de banco simplificada")
class DatabaseIntegrationCrudE2ETest {

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
                    .delete("/api/integrations/database/" + integrationId);
        }
    }

    @Test
    @Order(1)
    @DisplayName("POST /api/integrations/database cria integração de banco")
    void createDatabaseIntegration_returns201() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Database Integration",
                        "description", "Created by e2e",
                        "type", "POSTGRESQL",
                        "host", "pg.internal",
                        "port", 5432,
                        "database", "orders",
                        "dbUser", "orders_user",
                        "dbPassword", "orders_secret",
                        "query", "SELECT * FROM orders LIMIT 10",
                        "allowWrite", false
                ))
                .post("/api/integrations/database")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .body("name", equalTo("E2E Database Integration"))
                .extract().response();

        integrationId = resp.jsonPath().getString("id");
    }

    @Test
    @Order(2)
    @DisplayName("GET /api/integrations/database/{id} retorna integração")
    void getDatabaseIntegration_returns200() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/integrations/database/" + integrationId)
                .then()
                .statusCode(200)
                .body("id", equalTo(integrationId))
                .body("query", containsString("SELECT"));
    }

    @Test
    @Order(3)
    @DisplayName("PUT /api/integrations/database/{id} atualiza integração")
    void updateDatabaseIntegration_returns200() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Database Integration Updated",
                        "description", "Updated by e2e",
                        "type", "POSTGRESQL",
                        "host", "pg2.internal",
                        "port", 5432,
                        "database", "orders",
                        "dbUser", "orders_user",
                        "dbPassword", "",
                        "query", "UPDATE orders SET synced = true",
                        "allowWrite", true
                ))
                .put("/api/integrations/database/" + integrationId)
                .then()
                .statusCode(200)
                .body("id", equalTo(integrationId))
                .body("allowWrite", equalTo(true));
    }

    @Test
    @Order(4)
    @DisplayName("DELETE /api/integrations/database/{id} remove integração")
    void deleteDatabaseIntegration_returns204() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .delete("/api/integrations/database/" + integrationId)
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(204)));

        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/integrations/database/" + integrationId)
                .then()
                .statusCode(404);

        integrationId = null;
    }

    private static String fakeJwtForTenant(String tenantRealm) {
        String header = base64Url("{\"alg\":\"none\",\"typ\":\"JWT\"}");
        String payload = base64Url(("{" +
                "\"iss\":\"http://keycloak.cezar.dev/realms/" + tenantRealm + "\"," +
                "\"sub\":\"database-integration-e2e\"," +
                "\"preferred_username\":\"database-integration-e2e\"" +
                "}").replace("\n", ""));
        return header + "." + payload + ".signature";
    }

    private static String base64Url(String value) {
        return Base64.getUrlEncoder()
                .withoutPadding()
                .encodeToString(value.getBytes(StandardCharsets.UTF_8));
    }
}
