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
@DisplayName("Fluxo — CRUD de integração HTTP simplificada")
class HttpIntegrationCrudE2ETest {

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
                    .delete("/api/integrations/http/" + integrationId);
        }
    }

    @Test
    @Order(1)
    @DisplayName("POST /api/integrations/http cria integração HTTP simplificada")
    void createHttpIntegration_returns201() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E HTTP Integration",
                        "description", "Created by e2e",
                        "method", "POST",
                        "url", "https://api.example.com/customers",
                        "headers", Map.of("X-Test", "1"),
                        "bodyTemplate", "{\"name\":\"{name}\"}",
                        "readOnly", false
                ))
                .post("/api/integrations/http")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .body("name", equalTo("E2E HTTP Integration"))
                .body("method", equalTo("POST"))
                .extract().response();

        integrationId = resp.jsonPath().getString("id");
    }

    @Test
    @Order(2)
    @DisplayName("GET /api/integrations/http/{id} retorna integração")
    void getHttpIntegration_returns200() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/integrations/http/" + integrationId)
                .then()
                .statusCode(200)
                .body("id", equalTo(integrationId))
                .body("url", equalTo("https://api.example.com/customers"));
    }

    @Test
    @Order(3)
    @DisplayName("PUT /api/integrations/http/{id} atualiza integração")
    void updateHttpIntegration_returns200() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E HTTP Integration Updated",
                        "description", "Updated by e2e",
                        "method", "PATCH",
                        "url", "https://api.example.com/customers",
                        "headers", Map.of("X-Test", "2"),
                        "bodyTemplate", "{\"status\":\"updated\"}",
                        "readOnly", true
                ))
                .put("/api/integrations/http/" + integrationId)
                .then()
                .statusCode(200)
                .body("id", equalTo(integrationId))
                .body("name", equalTo("E2E HTTP Integration Updated"))
                .body("method", equalTo("PATCH"))
                .body("readOnly", equalTo(true));
    }

    @Test
    @Order(4)
    @DisplayName("DELETE /api/integrations/http/{id} remove integração")
    void deleteHttpIntegration_returns204() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .delete("/api/integrations/http/" + integrationId)
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(204)));

        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/integrations/http/" + integrationId)
                .then()
                .statusCode(404);

        integrationId = null;
    }

    private static String fakeJwtForTenant(String tenantRealm) {
        String header = base64Url("{\"alg\":\"none\",\"typ\":\"JWT\"}");
        String payload = base64Url(("{" +
                "\"iss\":\"http://keycloak.cezar.dev/realms/" + tenantRealm + "\"," +
                "\"sub\":\"http-integration-e2e\"," +
                "\"preferred_username\":\"http-integration-e2e\"" +
                "}").replace("\n", ""));
        return header + "." + payload + ".signature";
    }

    private static String base64Url(String value) {
        return Base64.getUrlEncoder()
                .withoutPadding()
                .encodeToString(value.getBytes(StandardCharsets.UTF_8));
    }
}
