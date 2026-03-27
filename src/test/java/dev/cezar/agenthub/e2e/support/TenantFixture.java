package dev.cezar.agenthub.e2e.support;

import dev.cezar.agenthub.e2e.config.E2EConfig;
import io.restassured.http.ContentType;
import io.restassured.response.Response;

import java.util.Map;
import java.util.UUID;

import static io.restassured.RestAssured.given;

/**
 * Creates and destroys test tenants. Each fixture instance represents one tenant.
 *
 * <p>Usage:
 * <pre>{@code
 *   TenantFixture tenant = TenantFixture.create();
 *   String jwt = tenant.getJwt();
 *   // ... run tests ...
 *   tenant.destroy();
 * }</pre>
 */
public class TenantFixture {

    private final String slug;
    private final String name;
    private final KeycloakClient keycloak;
    private String jwt;
    private boolean destroyed = false;

    private TenantFixture(String slug, String name) {
        this.slug     = slug;
        this.name     = name;
        this.keycloak = new KeycloakClient();
    }

    /**
     * Creates a new tenant via the backend and provisions a Keycloak admin user.
     * Returns a ready-to-use fixture with a valid JWT.
     */
    public static TenantFixture create() {
        String shortId = UUID.randomUUID().toString().substring(0, 8);
        String slug    = E2EConfig.E2E_TENANT_PREFIX + "-" + shortId;
        String name    = "E2E Test Tenant " + shortId;

        TenantFixture fixture = new TenantFixture(slug, name);
        fixture.provision();
        return fixture;
    }

    public String getSlug()  { return slug; }
    public String getName()  { return name; }
    public String getJwt()   { return jwt;  }

    /** Authorization header value ready for RestAssured. */
    public String bearerToken() { return "Bearer " + jwt; }

    /**
     * Refreshes the JWT (e.g. if it has expired during a long-running test).
     */
    public void refreshJwt() {
        try {
            jwt = keycloak.getUserToken(slug, "e2e-admin", E2EConfig.E2E_USER_PASSWORD);
        } catch (Exception e) {
            throw new RuntimeException("Failed to refresh JWT for tenant " + slug, e);
        }
    }

    /**
     * Deletes the Keycloak realm. The backend tenant record is left in place
     * (no public delete endpoint) but the realm cleanup is the critical part.
     */
    public void destroy() {
        if (!destroyed) {
            destroyed = true;
            keycloak.deleteRealm(slug);
        }
    }

    // ─── Private provisioning ─────────────────────────────────────────────────

    private void provision() {
        // 1. Create tenant via backend
        Response resp = given()
                .baseUri(E2EConfig.BACKEND_URL)
                .contentType(ContentType.JSON)
                .body(Map.of("name", name, "slug", slug))
                .post("/public/tenants")
                .then()
                .statusCode(201)
                .extract().response();

        System.out.println("[TenantFixture] Created tenant: " + slug + " (id=" + resp.jsonPath().getString("id") + ")");

        // 2. Wait briefly for Keycloak realm to be fully provisioned
        waitForRealm();

        // 3. Create admin user in the new realm
        try {
            keycloak.createAdminUser(slug, "e2e-admin", E2EConfig.E2E_USER_PASSWORD);
        } catch (Exception e) {
            throw new RuntimeException("Failed to create admin user in realm " + slug, e);
        }

        // 4. Obtain JWT
        refreshJwt();
    }

    private void waitForRealm() {
        int maxAttempts = 20;
        for (int i = 0; i < maxAttempts; i++) {
            try {
                // Try to fetch realm metadata
                Response r = given()
                        .baseUri(E2EConfig.KEYCLOAK_URL)
                        .get("/realms/" + slug)
                        .then()
                        .extract().response();
                if (r.statusCode() == 200) return;
            } catch (Exception ignored) {}
            try { Thread.sleep(500); } catch (InterruptedException e) { Thread.currentThread().interrupt(); return; }
        }
        throw new RuntimeException("Keycloak realm '" + slug + "' was not ready after " + maxAttempts + " attempts");
    }
}
