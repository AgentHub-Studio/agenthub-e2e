// Package testutil provides shared Testcontainers setup for AgentHub E2E tests.
package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Containers holds all running infrastructure containers for a test suite.
type Containers struct {
	Postgres *postgres.PostgresContainer
	RabbitMQ *rabbitmq.RabbitMQContainer
	Keycloak testcontainers.Container
	MinIO    testcontainers.Container

	PostgresDSN  string
	RabbitMQURL  string
	KeycloakURL  string
	MinIOAddress string
}

// StartContainers starts all required infrastructure containers.
// Call Terminate when done to clean up.
func StartContainers(ctx context.Context) (*Containers, error) {
	c := &Containers{}

	pgc, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("agenthub"),
		postgres.WithUsername("agenthub"),
		postgres.WithPassword("agenthub"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("start postgres: %w", err)
	}
	c.Postgres = pgc
	dsn, err := pgc.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("postgres connection string: %w", err)
	}
	c.PostgresDSN = dsn

	rmq, err := rabbitmq.Run(ctx, "rabbitmq:3.13-management-alpine",
		rabbitmq.WithAdminUsername("agenthub"),
		rabbitmq.WithAdminPassword("agenthub"),
	)
	if err != nil {
		_ = pgc.Terminate(ctx)
		return nil, fmt.Errorf("start rabbitmq: %w", err)
	}
	c.RabbitMQ = rmq
	amqpURL, err := rmq.AmqpURL(ctx)
	if err != nil {
		_ = pgc.Terminate(ctx)
		_ = rmq.Terminate(ctx)
		return nil, fmt.Errorf("rabbitmq amqp url: %w", err)
	}
	c.RabbitMQURL = amqpURL

	kc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "quay.io/keycloak/keycloak:24.0",
			Env: map[string]string{
				"KEYCLOAK_ADMIN":          "admin",
				"KEYCLOAK_ADMIN_PASSWORD": "admin",
				"KC_HOSTNAME_STRICT":      "false",
				"KC_HTTP_ENABLED":         "true",
			},
			Cmd:          []string{"start-dev"},
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForLog("Listening on:").WithStartupTimeout(120 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		_ = pgc.Terminate(ctx)
		_ = rmq.Terminate(ctx)
		return nil, fmt.Errorf("start keycloak: %w", err)
	}
	c.Keycloak = kc
	kcHost, err := kc.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("keycloak host: %w", err)
	}
	kcPort, err := kc.MappedPort(ctx, "8080")
	if err != nil {
		return nil, fmt.Errorf("keycloak port: %w", err)
	}
	c.KeycloakURL = fmt.Sprintf("http://%s:%s", kcHost, kcPort.Port())

	minio, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "minio/minio:latest",
			Env: map[string]string{
				"MINIO_ROOT_USER":     "minioadmin",
				"MINIO_ROOT_PASSWORD": "minioadmin",
			},
			Cmd:          []string{"server", "/data"},
			ExposedPorts: []string{"9000/tcp"},
			WaitingFor:   wait.ForHTTP("/minio/health/ready").WithPort("9000/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		_ = pgc.Terminate(ctx)
		_ = rmq.Terminate(ctx)
		_ = kc.Terminate(ctx)
		return nil, fmt.Errorf("start minio: %w", err)
	}
	c.MinIO = minio
	minioHost, err := minio.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("minio host: %w", err)
	}
	minioPort, err := minio.MappedPort(ctx, "9000")
	if err != nil {
		return nil, fmt.Errorf("minio port: %w", err)
	}
	c.MinIOAddress = fmt.Sprintf("%s:%s", minioHost, minioPort.Port())

	return c, nil
}

// Terminate stops and removes all containers.
func (c *Containers) Terminate(ctx context.Context) {
	if c.Postgres != nil {
		_ = c.Postgres.Terminate(ctx)
	}
	if c.RabbitMQ != nil {
		_ = c.RabbitMQ.Terminate(ctx)
	}
	if c.Keycloak != nil {
		_ = c.Keycloak.Terminate(ctx)
	}
	if c.MinIO != nil {
		_ = c.MinIO.Terminate(ctx)
	}
}
