package testutil

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ServiceConfig holds connection info for a running agenthub service.
type ServiceConfig struct {
	APIURL          string
	ObservabilityURL string
}

// DefaultServiceConfig returns a ServiceConfig reading URLs from environment variables.
// These are set when tests run against already-deployed services.
//
// Env vars:
//
//	API_URL             — base URL for agenthub-api (default: http://localhost:8081)
//	OBSERVABILITY_URL   — base URL for agenthub-observability (default: http://localhost:8086)
func DefaultServiceConfig() ServiceConfig {
	return ServiceConfig{
		APIURL:          getEnvOrDefault("API_URL", "http://localhost:8081"),
		ObservabilityURL: getEnvOrDefault("OBSERVABILITY_URL", "http://localhost:8086"),
	}
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

// APIServiceContainer starts agenthub-api as a Testcontainer using the local Docker image.
// The image must exist locally (build via `cd agenthub-api && ./build.sh package`).
// Returns the container base URL and a cleanup function.
func APIServiceContainer(ctx context.Context, pgDSN, keycloakURL string) (string, func(), error) {
	req := testcontainers.ContainerRequest{
		Image: "agenthub-api:latest",
		Env: map[string]string{
			"DATABASE_URL":     pgDSN,
			"KEYCLOAK_BASE_URL": keycloakURL,
			"PORT":             "8081",
			"LOG_LEVEL":        "warn",
		},
		ExposedPorts: []string{"8081/tcp"},
		WaitingFor: wait.ForHTTP("/health").
			WithPort("8081/tcp").
			WithStartupTimeout(60 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, fmt.Errorf("start agenthub-api container: %w", err)
	}

	host, err := c.Host(ctx)
	if err != nil {
		_ = c.Terminate(ctx)
		return "", nil, fmt.Errorf("api container host: %w", err)
	}
	port, err := c.MappedPort(ctx, "8081")
	if err != nil {
		_ = c.Terminate(ctx)
		return "", nil, fmt.Errorf("api container port: %w", err)
	}

	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())
	cleanup := func() { _ = c.Terminate(ctx) }
	return baseURL, cleanup, nil
}

// ObservabilityServiceContainer starts agenthub-observability as a Testcontainer.
func ObservabilityServiceContainer(ctx context.Context, clickhouseDSN, rabbitMQURL string) (string, func(), error) {
	req := testcontainers.ContainerRequest{
		Image: "agenthub-observability:latest",
		Env: map[string]string{
			"CLICKHOUSE_DSN": clickhouseDSN,
			"RABBITMQ_URL":   rabbitMQURL,
			"PORT":           "8086",
			"LOG_LEVEL":      "warn",
		},
		ExposedPorts: []string{"8086/tcp"},
		WaitingFor: wait.ForHTTP("/health").
			WithPort("8086/tcp").
			WithStartupTimeout(60 * time.Second),
	}

	c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", nil, fmt.Errorf("start observability container: %w", err)
	}

	host, err := c.Host(ctx)
	if err != nil {
		_ = c.Terminate(ctx)
		return "", nil, fmt.Errorf("observability container host: %w", err)
	}
	port, err := c.MappedPort(ctx, "8086")
	if err != nil {
		_ = c.Terminate(ctx)
		return "", nil, fmt.Errorf("observability container port: %w", err)
	}

	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())
	cleanup := func() { _ = c.Terminate(ctx) }
	return baseURL, cleanup, nil
}
