package dev.cezar.agenthub.e2e.config;

/**
 * Central configuration for e2e tests. All values are read from environment variables
 * with sensible defaults for the local Docker Compose stack.
 */
public final class E2EConfig {

    public static final String BACKEND_URL =
            env("BACKEND_URL", "http://localhost:8081");

    public static final String KEYCLOAK_URL =
            env("KEYCLOAK_URL", "http://localhost:8080");

    public static final String KEYCLOAK_ADMIN_USER =
            env("KEYCLOAK_ADMIN_USER", "admin");

    public static final String KEYCLOAK_ADMIN_PASSWORD =
            env("KEYCLOAK_ADMIN_PASSWORD", "@admin#");

    public static final String EMBEDDING_URL =
            env("EMBEDDING_URL", "http://localhost:8092");

    public static final String EXTRACTOR_URL =
            env("EXTRACTOR_URL", "http://localhost:8093");

    public static final String E2E_TENANT_PREFIX =
            env("E2E_TENANT_PREFIX", "e2e-test");

    /** Default admin user password created in new realms during e2e tests. */
    public static final String E2E_USER_PASSWORD =
            env("E2E_USER_PASSWORD", "E2eTestPass#1");

    private E2EConfig() {}

    private static String env(String key, String defaultValue) {
        String value = System.getenv(key);
        return (value != null && !value.isBlank()) ? value : defaultValue;
    }
}
