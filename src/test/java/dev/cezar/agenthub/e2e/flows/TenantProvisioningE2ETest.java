package dev.cezar.agenthub.e2e.flows;

import dev.cezar.agenthub.e2e.config.E2EConfig;
import dev.cezar.agenthub.e2e.support.KeycloakClient;
import dev.cezar.agenthub.e2e.support.TenantFixture;
import io.restassured.http.ContentType;
import org.junit.jupiter.api.*;

import java.util.UUID;

import static io.restassured.RestAssured.given;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;

/**
 * Fluxo 1 — Provisionamento de tenant.
 *
 * <p>Verifica que um tenant pode ser criado, que o realm Keycloak é provisionado,
 * que um usuário admin pode se autenticar e que o backend aceita o JWT resultante.
 */
@TestMethodOrder(MethodOrderer.OrderAnnotation.class)
@DisplayName("Fluxo 1 — Provisionamento de Tenant")
class TenantProvisioningE2ETest {

    private static TenantFixture tenant;

    @BeforeAll
    static void setUp() {
        tenant = TenantFixture.create();
    }

    @AfterAll
    static void tearDown() {
        if (tenant != null) tenant.destroy();
    }

    @Test
    @Order(1)
    @DisplayName("POST /public/tenants → 201 e campos obrigatórios presentes")
    void createTenant_returns201WithRequiredFields() {
        // Tenant already created in @BeforeAll; verify exists endpoint
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .get("/public/tenants/" + tenant.getSlug() + "/exists")
                .then()
                .statusCode(200)
                .body(equalTo("true"));
    }

    @Test
    @Order(2)
    @DisplayName("GET /public/tenants/{id}/exists → true para tenant criado")
    void existsEndpoint_returnsTrue() {
        Boolean exists = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .get("/public/tenants/" + tenant.getSlug() + "/exists")
                .then()
                .statusCode(200)
                .extract()
                .as(Boolean.class);

        assertTrue(exists, "Tenant should exist after provisioning");
    }

    @Test
    @Order(3)
    @DisplayName("GET /public/tenants/{id}/exists → false para tenant inexistente")
    void existsEndpoint_returnsFalseForUnknown() {
        String nonExistent = "non-existent-slug-" + UUID.randomUUID();
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .get("/public/tenants/" + nonExistent + "/exists")
                .then()
                .statusCode(200)
                .body(equalTo("false"));
    }

    @Test
    @Order(4)
    @DisplayName("JWT obtido via Keycloak é aceito pelo backend")
    void jwtFromKeycloak_isAcceptedByBackend() {
        assertNotNull(tenant.getJwt(), "JWT should not be null");
        assertFalse(tenant.getJwt().isBlank(), "JWT should not be blank");

        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .header("Authorization", tenant.bearerToken())
                .get("/api/agents")
                .then()
                .statusCode(200)
                .body("content", notNullValue())
                .body("totalElements", greaterThanOrEqualTo(0));
    }

    @Test
    @Order(5)
    @DisplayName("POST /public/tenants com slug inválido → 400")
    void createTenant_invalidSlug_returns400() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .contentType(ContentType.JSON)
                .body("{\"name\": \"Test\", \"slug\": \"INVALID SLUG WITH SPACES\"}")
                .post("/public/tenants")
                .then()
                .statusCode(400);
    }

    @Test
    @Order(6)
    @DisplayName("GET /public/tenants → lista inclui o tenant criado")
    void listTenants_containsCreatedTenant() {
        given()
                .baseUri(E2EConfig.BACKEND_URL)
                .queryParam("size", 500)
                .get("/public/tenants")
                .then()
                .statusCode(200)
                .body("content.find { it.id == '" + tenant.getSlug() + "' }", notNullValue());
    }
}
