package dev.cezar.agenthub.e2e.flows;

import dev.cezar.agenthub.e2e.config.E2EConfig;
import dev.cezar.agenthub.e2e.support.TenantFixture;
import io.restassured.http.ContentType;
import io.restassured.response.Response;
import org.junit.jupiter.api.*;

import java.util.Map;
import java.util.UUID;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Fluxo 2 — Ciclo completo de Skill + Tool.
 *
 * <p>Cria, lê, lista e deleta uma skill com tool vinculada,
 * verificando paginação e respostas corretas em cada etapa.
 */
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("Fluxo 2 — Ciclo Completo de Skill + Tool")
class SkillToolCrudE2ETest {

    private static TenantFixture tenant;
    private static String skillId;
    private static String toolId;

    @BeforeAll
    static void setUp() {
        tenant = TenantFixture.create();
    }

    @AfterAll
    static void tearDown() {
        try {
            // Cleanup tool and skill if test failed mid-way
            if (toolId != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", tenant.bearerToken())
                        .delete("/api/tools/" + toolId);
            }
            if (skillId != null) {
                given().baseUri(E2EConfig.BACKEND_URL)
                        .header("Authorization", tenant.bearerToken())
                        .delete("/api/skills/" + skillId);
            }
        } finally {
            if (tenant != null) tenant.destroy();
        }
    }

    @Test
    @Order(1)
    @DisplayName("POST /api/skills → 201 e skill criada com campos corretos")
    void createSkill_returns201() {
        String slug = "e2e-skill-" + UUID.randomUUID().toString().substring(0, 8);

        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Test Skill",
                        "slug", slug,
                        "description", "Created by e2e test",
                        "type", "DATA",
                        "category", "DATA",
                        "version", "1.0.0"
                ))
                .post("/api/skills")
                .then()
                .statusCode(201)
                .body("id", notNullValue())
                .body("slug", equalTo(slug))
                .body("status", equalTo("ACTIVE"))
                .extract().response();

        skillId = resp.jsonPath().getString("id");
        assertNotNull(skillId, "Skill ID should not be null");
    }

    @Test
    @Order(2)
    @DisplayName("POST /api/tools → 200/201 e tool criada com tipo HTTP")
    void createTool_returns201() {
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "E2E Test Tool",
                        "description", "HTTP tool created by e2e test",
                        "type", "HTTP",
                        "config", Map.of(
                                "url", "https://httpbin.org/get",
                                "method", "GET"
                        ),
                        "labels", new String[]{"e2e", "test"}
                ))
                .post("/api/tools")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)))
                .body("id", notNullValue())
                .body("type", notNullValue())
                .extract().response();

        toolId = resp.jsonPath().getString("id");
        assertNotNull(toolId, "Tool ID should not be null");
    }

    @Test
    @Order(3)
    @DisplayName("POST /api/skills/{id}/tools → vincula tool à skill")
    void bindToolToSkill_returns201() {
        assertNotNull(skillId, "Skill must be created first");
        assertNotNull(toolId, "Tool must be created first");

        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of("toolId", toolId, "priority", 0))
                .post("/api/skills/" + skillId + "/tools")
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(201)));
    }

    @Test
    @Order(4)
    @DisplayName("GET /api/skills?page=0&size=20 → estrutura Page<T> correta")
    void listSkills_returnsPageStructure() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .queryParam("page", 0)
                .queryParam("size", 20)
                .get("/api/skills")
                .then()
                .statusCode(200)
                .body("content", notNullValue())
                .body("totalElements", greaterThanOrEqualTo(1))
                .body("size", equalTo(20))
                .body("number", equalTo(0))
                .body("content.find { it.id == '" + skillId + "' }", notNullValue());
    }

    @Test
    @Order(5)
    @DisplayName("GET /api/skills/{id} → retorna skill com campos corretos")
    void getSkill_returnsCorrectFields() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .get("/api/skills/" + skillId)
                .then()
                .statusCode(200)
                .body("id", equalTo(skillId))
                .body("name", equalTo("E2E Test Skill"))
                .body("category", equalTo("DATA"))
                .body("version", equalTo("1.0.0"))
                .body("status", equalTo("ACTIVE"));
    }

    @Test
    @Order(6)
    @DisplayName("GET /api/tools/{id} → retorna tool com campos corretos")
    void getTool_returnsCorrectFields() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .get("/api/tools/" + toolId)
                .then()
                .statusCode(200)
                .body("id", equalTo(toolId))
                .body("name", equalTo("E2E Test Tool"))
                .body("type", notNullValue());
    }

    @Test
    @Order(7)
    @DisplayName("DELETE /api/tools/{id} → 200/204")
    void deleteTool_returns204() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .delete("/api/tools/" + toolId)
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(204)));
        toolId = null; // Prevent double-delete in tearDown
    }

    @Test
    @Order(8)
    @DisplayName("DELETE /api/skills/{id} → 200/204")
    void deleteSkill_returns204() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .delete("/api/skills/" + skillId)
                .then()
                .statusCode(anyOf(equalTo(200), equalTo(204)));
        skillId = null; // Prevent double-delete in tearDown
    }

    @Test
    @Order(9)
    @DisplayName("GET /api/skills/{id} após deleção → 404")
    void getSkill_afterDeletion_returns404() {
        // Re-create and immediately delete to test 404 (skillId is null from test 8)
        String slug = "e2e-del-" + UUID.randomUUID().toString().substring(0, 8);
        String id = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .contentType(ContentType.JSON)
                .body(Map.of(
                        "name", "To Delete",
                        "slug", slug,
                        "version", "1.0.0",
                        "type", "DATA",
                        "category", "DATA"
                ))
                .post("/api/skills")
                .then().statusCode(201)
                .extract().jsonPath().getString("id");

        given().baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .delete("/api/skills/" + id)
                .then().statusCode(anyOf(equalTo(200), equalTo(204)));

        given().baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .get("/api/skills/" + id)
                .then().statusCode(anyOf(equalTo(404), equalTo(400)));
    }
}
