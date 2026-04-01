// Package tests provides end-to-end tests for AgentHub.
// Requires running agenthub-api and infrastructure.
//
// Run with:
//
//	API_URL=http://localhost:8081 go test ./tests/... -v
//
// Or with Testcontainers (full infra spin-up):
//
//	E2E_INFRA=1 go test ./tests/... -v -timeout 10m
package tests

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AgentHub-Studio/agenthub-e2e/internal/testutil"
)

func apiURL() string {
	if v := os.Getenv("API_URL"); v != "" {
		return v
	}
	return "http://localhost:8081"
}

func authToken() string {
	return os.Getenv("AUTH_TOKEN")
}

func skipIfNotRunning(t *testing.T) {
	t.Helper()
	if os.Getenv("E2E") == "" {
		t.Skip("set E2E=1 to run (requires agenthub-api running)")
	}
}

// TestTenantProvisioning_CreateAndList verifies the tenant creation flow.
func TestTenantProvisioning_CreateAndList(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())
	fixture := testutil.NewTenantFixture()

	var created map[string]any
	status, err := client.Post("/public/tenants", fixture, &created)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, fixture.ID, created["id"])
	assert.NotEmpty(t, created["id"])

	// Verify the tenant appears in the list.
	var page testutil.Page[map[string]any]
	status, err = client.Get("/public/tenants", &page)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.GreaterOrEqual(t, page.TotalElements, int64(1))
}

// TestTenantExists_AfterCreation verifies that GET /public/tenants/{id}/exists returns true.
func TestTenantExists_AfterCreation(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())
	fixture := testutil.NewTenantFixture()

	_, err := client.Post("/public/tenants", fixture, nil)
	require.NoError(t, err)

	var exists map[string]any
	status, err := client.Get("/public/tenants/"+fixture.ID+"/exists", &exists)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, true, exists["exists"])
}

// TestTenantExists_NotFound verifies that a non-existent tenant returns false.
func TestTenantExists_NotFound(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	var exists map[string]any
	status, err := client.Get("/public/tenants/nonexistent-tenant/exists", &exists)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, false, exists["exists"])
}

// TestMultitenantIsolation verifies that tenant schemas are isolated.
func TestMultitenantIsolation(t *testing.T) {
	skipIfNotRunning(t)
	ctx := context.Background()
	_ = ctx // used for future container checks

	clientA := testutil.NewClient(apiURL(), os.Getenv("AUTH_TOKEN_TENANT_A"))
	clientB := testutil.NewClient(apiURL(), os.Getenv("AUTH_TOKEN_TENANT_B"))

	if os.Getenv("AUTH_TOKEN_TENANT_A") == "" || os.Getenv("AUTH_TOKEN_TENANT_B") == "" {
		t.Skip("set AUTH_TOKEN_TENANT_A and AUTH_TOKEN_TENANT_B to run isolation tests")
	}

	agentA := testutil.NewAgentFixture()
	var createdA map[string]any
	status, err := clientA.Post("/api/agents", agentA, &createdA)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	agentID := createdA["id"].(string)

	// Tenant B should not see tenant A's agents.
	var agentFromB map[string]any
	status, err = clientB.Get("/api/agents/"+agentID, &agentFromB)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}
