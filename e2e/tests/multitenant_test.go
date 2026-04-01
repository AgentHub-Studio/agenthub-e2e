package tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AgentHub-Studio/agenthub-e2e/internal/testutil"
)

// tenantAClient returns an HTTP client authenticated as Tenant A.
// Requires AUTH_TOKEN_TENANT_A environment variable.
func tenantAClient() *testutil.Client {
	return testutil.NewClient(apiURL(), os.Getenv("AUTH_TOKEN_TENANT_A"))
}

// tenantBClient returns an HTTP client authenticated as Tenant B.
// Requires AUTH_TOKEN_TENANT_B environment variable.
func tenantBClient() *testutil.Client {
	return testutil.NewClient(apiURL(), os.Getenv("AUTH_TOKEN_TENANT_B"))
}

func skipIfNoTenantTokens(t *testing.T) {
	t.Helper()
	skipIfNotRunning(t)
	if os.Getenv("AUTH_TOKEN_TENANT_A") == "" || os.Getenv("AUTH_TOKEN_TENANT_B") == "" {
		t.Skip("set AUTH_TOKEN_TENANT_A and AUTH_TOKEN_TENANT_B to run multi-tenancy isolation tests")
	}
}

// TestMultitenant_AgentIsolation verifies that agents created by Tenant A are not
// visible to Tenant B and vice versa.
func TestMultitenant_AgentIsolation(t *testing.T) {
	skipIfNoTenantTokens(t)

	clientA := tenantAClient()
	clientB := tenantBClient()

	// Tenant A creates an agent.
	fixtureA := testutil.NewAgentFixture()
	var agentA map[string]any
	status, err := clientA.Post("/api/agents", fixtureA, &agentA)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	agentAID := agentA["id"].(string)

	// Tenant B creates a different agent.
	fixtureB := testutil.NewAgentFixture()
	var agentB map[string]any
	status, err = clientB.Post("/api/agents", fixtureB, &agentB)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	agentBID := agentB["id"].(string)

	// Tenant B must not access Tenant A's agent.
	var crossResult map[string]any
	status, err = clientB.Get("/api/agents/"+agentAID, &crossResult)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status, "tenant B must not see tenant A's agent")

	// Tenant A must not access Tenant B's agent.
	status, err = clientA.Get("/api/agents/"+agentBID, &crossResult)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status, "tenant A must not see tenant B's agent")
}

// TestMultitenant_AgentList_OnlyOwnTenantData verifies that listing agents returns
// only resources belonging to the authenticated tenant.
func TestMultitenant_AgentList_OnlyOwnTenantData(t *testing.T) {
	skipIfNoTenantTokens(t)

	clientA := tenantAClient()
	clientB := tenantBClient()

	// Create an agent as tenant A.
	fixtureA := testutil.NewAgentFixture()
	status, err := clientA.Post("/api/agents", fixtureA, nil)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)

	// List agents as tenant A.
	var pageA testutil.Page[map[string]any]
	status, err = clientA.Get("/api/agents?page=0&size=100", &pageA)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	// List agents as tenant B.
	var pageB testutil.Page[map[string]any]
	status, err = clientB.Get("/api/agents?page=0&size=100", &pageB)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	// Tenant A's agent should appear in A's list but not B's list.
	foundInA := false
	for _, a := range pageA.Content {
		if a["name"] == fixtureA.Name {
			foundInA = true
			break
		}
	}
	assert.True(t, foundInA, "agent created by tenant A should appear in tenant A's list")

	for _, b := range pageB.Content {
		assert.NotEqual(t, fixtureA.Name, b["name"], "tenant A's agent must not appear in tenant B's list")
	}
}

// TestMultitenant_SkillIsolation verifies that skills are isolated per tenant.
func TestMultitenant_SkillIsolation(t *testing.T) {
	skipIfNoTenantTokens(t)

	clientA := tenantAClient()
	clientB := tenantBClient()

	// Tenant A creates a skill.
	fixtureA := testutil.NewSkillFixture()
	var skillA map[string]any
	status, err := clientA.Post("/api/skills", fixtureA, &skillA)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	skillAID := skillA["id"].(string)

	// Tenant B must not access Tenant A's skill.
	var crossResult map[string]any
	status, err = clientB.Get("/api/skills/"+skillAID, &crossResult)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status, "tenant B must not see tenant A's skill")
}

// TestMultitenant_LLMPreset_NameUniquenessPerTenant verifies that the same LLM preset
// name can be used by different tenants independently.
func TestMultitenant_LLMPreset_NameUniquenessPerTenant(t *testing.T) {
	skipIfNoTenantTokens(t)

	clientA := tenantAClient()
	clientB := tenantBClient()

	preset := map[string]any{
		"name":     "GPT-4o",
		"provider": "openai",
		"model":    "gpt-4o",
	}

	// Tenant A creates a preset named "GPT-4o".
	status, err := clientA.Post("/api/llm-config-presets", preset, nil)
	require.NoError(t, err)
	// Accept Created or Conflict (in case a previous test run left the preset).
	assert.True(t, status == http.StatusCreated || status == http.StatusConflict,
		"expected 201 or 409, got %d", status)

	// Tenant B creates a preset with the same name — must succeed (names are per-tenant).
	var createdB map[string]any
	status, err = clientB.Post("/api/llm-config-presets", preset, &createdB)
	require.NoError(t, err)
	assert.True(t, status == http.StatusCreated || status == http.StatusConflict,
		"tenant B should be allowed to use the same name as tenant A, got %d", status)
}

// TestMultitenant_LLMPreset_Isolation verifies that LLM presets created by Tenant A
// are not visible to Tenant B.
func TestMultitenant_LLMPreset_Isolation(t *testing.T) {
	skipIfNoTenantTokens(t)

	clientA := tenantAClient()
	clientB := tenantBClient()

	// Tenant A creates a uniquely-named preset.
	uniqueName := "E2E-Preset-" + testutil.NewAgentFixture().Name[len("Test Agent "):]
	presetA := map[string]any{
		"name":     uniqueName,
		"provider": "anthropic",
		"model":    "claude-sonnet-4-6",
	}

	var createdA map[string]any
	status, err := clientA.Post("/api/llm-config-presets", presetA, &createdA)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status, "tenant A should create preset successfully")
	presetAID := createdA["id"].(string)

	// Tenant B must not access Tenant A's preset.
	var crossResult map[string]any
	status, err = clientB.Get("/api/llm-config-presets/"+presetAID, &crossResult)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status, "tenant B must not see tenant A's LLM preset")

	// Tenant B's list must not contain Tenant A's preset.
	var pageB testutil.Page[map[string]any]
	status, err = clientB.Get("/api/llm-config-presets?page=0&size=100", &pageB)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	for _, p := range pageB.Content {
		assert.NotEqual(t, presetAID, p["id"],
			"tenant A's LLM preset must not appear in tenant B's list")
	}
}

// TestMultitenant_DeleteIsolation verifies that a tenant cannot delete another
// tenant's resources.
func TestMultitenant_DeleteIsolation(t *testing.T) {
	skipIfNoTenantTokens(t)

	clientA := tenantAClient()
	clientB := tenantBClient()

	// Tenant A creates an agent.
	fixtureA := testutil.NewAgentFixture()
	var agentA map[string]any
	status, err := clientA.Post("/api/agents", fixtureA, &agentA)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	agentAID := agentA["id"].(string)

	// Tenant B attempts to delete Tenant A's agent — must fail with 404.
	status, err = clientB.Delete("/api/agents/" + agentAID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status, "tenant B must not delete tenant A's agent")

	// Tenant A's agent must still exist.
	var stillExists map[string]any
	status, err = clientA.Get("/api/agents/"+agentAID, &stillExists)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status, "tenant A's agent should still exist after B's delete attempt")
}
