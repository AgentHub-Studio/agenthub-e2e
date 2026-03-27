package dev.cezar.agenthub.e2e.support;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import dev.cezar.agenthub.e2e.config.E2EConfig;

import java.net.URI;
import java.net.URLEncoder;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.nio.charset.StandardCharsets;
import java.util.Map;

/**
 * Thin wrapper around the Keycloak Admin REST API and token endpoint.
 * Uses java.net.http to avoid extra dependencies.
 */
public class KeycloakClient {

    private static final HttpClient HTTP = HttpClient.newHttpClient();
    private static final ObjectMapper MAPPER = new ObjectMapper();

    private final String keycloakUrl;
    private final String adminUser;
    private final String adminPassword;

    public KeycloakClient() {
        this.keycloakUrl = E2EConfig.KEYCLOAK_URL;
        this.adminUser   = E2EConfig.KEYCLOAK_ADMIN_USER;
        this.adminPassword = E2EConfig.KEYCLOAK_ADMIN_PASSWORD;
    }

    // ─── Master realm admin token ─────────────────────────────────────────────

    public String getAdminToken() {
        return fetchToken("master", "admin-cli", adminUser, adminPassword);
    }

    // ─── Token for a specific realm ───────────────────────────────────────────

    /**
     * Obtains a JWT via Resource Owner Password Grant for the given tenant realm.
     * The Keycloak client must have "Direct Access Grants" enabled.
     */
    public String getUserToken(String realm, String username, String password) {
        return fetchToken(realm, "agenthub-frontend", username, password);
    }

    // ─── User management ──────────────────────────────────────────────────────

    /**
     * Creates a user in the given realm and assigns them the "admin" realm role.
     * Returns the created user's ID.
     */
    public String createAdminUser(String realm, String username, String password) throws Exception {
        String adminToken = getAdminToken();

        // 1. Create user
        String body = MAPPER.writeValueAsString(Map.of(
                "username", username,
                "enabled", true,
                "emailVerified", true,
                "email", username + "@e2e.test",
                "credentials", new Object[]{
                        Map.of("type", "password", "value", password, "temporary", false)
                }
        ));

        HttpRequest createReq = HttpRequest.newBuilder()
                .uri(URI.create(keycloakUrl + "/admin/realms/" + realm + "/users"))
                .header("Authorization", "Bearer " + adminToken)
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(body))
                .build();

        HttpResponse<String> createResp = HTTP.send(createReq, HttpResponse.BodyHandlers.ofString());
        if (createResp.statusCode() != 201) {
            throw new RuntimeException("Failed to create user in realm " + realm + ": " + createResp.statusCode() + " " + createResp.body());
        }

        // 2. Get user ID from Location header
        String location = createResp.headers().firstValue("Location")
                .orElseThrow(() -> new RuntimeException("No Location header in user creation response"));
        String userId = location.substring(location.lastIndexOf('/') + 1);

        // 3. Get "admin" role representation
        HttpRequest roleReq = HttpRequest.newBuilder()
                .uri(URI.create(keycloakUrl + "/admin/realms/" + realm + "/roles/admin"))
                .header("Authorization", "Bearer " + adminToken)
                .GET()
                .build();

        HttpResponse<String> roleResp = HTTP.send(roleReq, HttpResponse.BodyHandlers.ofString());
        if (roleResp.statusCode() != 200) {
            throw new RuntimeException("Failed to get admin role in realm " + realm + ": " + roleResp.statusCode());
        }

        // 4. Assign "admin" role to user
        String rolesBody = "[" + roleResp.body() + "]";
        HttpRequest assignReq = HttpRequest.newBuilder()
                .uri(URI.create(keycloakUrl + "/admin/realms/" + realm + "/users/" + userId + "/role-mappings/realm"))
                .header("Authorization", "Bearer " + adminToken)
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(rolesBody))
                .build();

        HttpResponse<String> assignResp = HTTP.send(assignReq, HttpResponse.BodyHandlers.ofString());
        if (assignResp.statusCode() != 204) {
            throw new RuntimeException("Failed to assign admin role: " + assignResp.statusCode());
        }

        return userId;
    }

    /**
     * Deletes a realm from Keycloak. No-op if the realm doesn't exist.
     */
    public void deleteRealm(String realm) {
        try {
            String adminToken = getAdminToken();
            HttpRequest req = HttpRequest.newBuilder()
                    .uri(URI.create(keycloakUrl + "/admin/realms/" + realm))
                    .header("Authorization", "Bearer " + adminToken)
                    .DELETE()
                    .build();
            HTTP.send(req, HttpResponse.BodyHandlers.ofString());
        } catch (Exception e) {
            System.err.println("[KeycloakClient] Warning: could not delete realm " + realm + ": " + e.getMessage());
        }
    }

    // ─── Internal helpers ─────────────────────────────────────────────────────

    private String fetchToken(String realm, String clientId, String username, String password) {
        try {
            String form = "grant_type=password"
                    + "&client_id=" + URLEncoder.encode(clientId, StandardCharsets.UTF_8)
                    + "&username=" + URLEncoder.encode(username, StandardCharsets.UTF_8)
                    + "&password=" + URLEncoder.encode(password, StandardCharsets.UTF_8);

            HttpRequest req = HttpRequest.newBuilder()
                    .uri(URI.create(keycloakUrl + "/realms/" + realm + "/protocol/openid-connect/token"))
                    .header("Content-Type", "application/x-www-form-urlencoded")
                    .POST(HttpRequest.BodyPublishers.ofString(form))
                    .build();

            HttpResponse<String> resp = HTTP.send(req, HttpResponse.BodyHandlers.ofString());
            if (resp.statusCode() != 200) {
                throw new RuntimeException("Token request failed (" + resp.statusCode() + "): " + resp.body());
            }

            JsonNode node = MAPPER.readTree(resp.body());
            return node.get("access_token").asText();
        } catch (RuntimeException e) {
            throw e;
        } catch (Exception e) {
            throw new RuntimeException("fetchToken failed", e);
        }
    }
}
