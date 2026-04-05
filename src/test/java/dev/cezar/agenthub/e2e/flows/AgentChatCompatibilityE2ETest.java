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
@DisplayName("Fluxo — Compatibilidade entre agent e chat")
class AgentChatCompatibilityE2ETest {

    private static final String TENANT_REALM = "test";

    private static String bearerToken;
    private static String agentId;
    private static String sessionId;

    @BeforeAll
    static void setUp() {
        bearerToken = "Bearer " + fakeJwtForTenant(TENANT_REALM);
    }

    @AfterAll
    static void tearDown() {
        if (sessionId != null) {
            given()
                    .baseUri(E2EConfig.BACKEND_URL)
                    .header("Authorization", bearerToken)
                    .delete("/api/chat/sessions/" + sessionId);
        }
        if (agentId != null) {
            given()
                    .baseUri(E2EConfig.BACKEND_URL)
                    .header("Authorization", bearerToken)
                    .delete("/api/agents/" + agentId);
        }
    }

    @Test
    @Order(1)
    @DisplayName("POST /api/agents cria agente")
    void createAgent_returns201() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Compatibility Agent",
                        "description", "Created to validate agent/chat compatibility",
                        "systemPrompt", "You are a compatibility validation agent."
                ))
                .post("/api/agents")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)))
                .body("id", notNullValue())
                .extract().response();

        agentId = resp.jsonPath().getString("id");
    }

    @Test
    @Order(2)
    @DisplayName("GET /api/agents lista o agente criado")
    void listAgents_containsCreatedAgent() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/agents?page=0&size=20")
                .then()
                .statusCode(200)
                .body("content.find { it.id == '" + agentId + "' }", notNullValue());
    }

    @Test
    @Order(3)
    @DisplayName("POST /api/chat/sessions cria sessão associada ao agente")
    void createSession_returns201() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "title", "Compatibility Session",
                        "agentId", agentId
                ))
                .post("/api/chat/sessions")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)))
                .body("id", notNullValue())
                .extract().response();

        sessionId = resp.jsonPath().getString("id");
    }

    @Test
    @Order(4)
    @DisplayName("GET /api/chat/sessions e /messages continuam compatíveis")
    void listSessionsAndMessages_returnsPageShape() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/chat/sessions?page=0&size=20")
                .then()
                .statusCode(200)
                .body("content.find { it.id == '" + sessionId + "' }", notNullValue());

        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", bearerToken)
                .get("/api/chat/sessions/" + sessionId + "/messages?page=0&size=20")
                .then()
                .statusCode(200)
                .body("content", notNullValue())
                .body("totalElements", greaterThanOrEqualTo(0));
    }

    private static String fakeJwtForTenant(String tenantRealm) {
        String header = base64Url("{\"alg\":\"none\",\"typ\":\"JWT\"}");
        String payload = base64Url(("{" +
                "\"iss\":\"http://keycloak.cezar.dev/realms/" + tenantRealm + "\"," +
                "\"sub\":\"agent-chat-compat-e2e\"," +
                "\"preferred_username\":\"agent-chat-compat-e2e\"" +
                "}").replace("\n", ""));
        return header + "." + payload + ".signature";
    }

    private static String base64Url(String value) {
        return Base64.getUrlEncoder()
                .withoutPadding()
                .encodeToString(value.getBytes(StandardCharsets.UTF_8));
    }
}
