package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AgentHub-Studio/agenthub-e2e/internal/testutil"
)

// TestAgentCRUD verifies create, read, update, delete of an agent.
func TestAgentCRUD(t *testing.T) {
	skipIfNotRunning(t)
	client := testutil.NewClient(apiURL(), authToken())

	fixture := testutil.NewAgentFixture()

	// Create
	var created map[string]any
	status, err := client.Post("/api/agents", fixture, &created)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	id := created["id"].(string)
	assert.Equal(t, fixture.Name, created["name"])

	// Read
	var fetched map[string]any
	status, err = client.Get("/api/agents/"+id, &fetched)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, id, fetched["id"])

	// Update
	update := map[string]any{"name": fixture.Name + " Updated"}
	var updated map[string]any
	status, err = client.Put("/api/agents/"+id, update, &updated)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, fixture.Name+" Updated", updated["name"])

	// Delete
	status, err = client.Delete("/api/agents/" + id)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, status)

	// Verify 404 after delete
	status, err = client.Get("/api/agents/"+id, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

// TestAgentList_Pagination verifies that agent list returns a paginated response.
func TestAgentList_Pagination(t *testing.T) {
	skipIfNotRunning(t)
	client := testutil.NewClient(apiURL(), authToken())

	var page testutil.Page[map[string]any]
	status, err := client.Get("/api/agents?page=0&size=10", &page)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.NotNil(t, page.Content)
	assert.GreaterOrEqual(t, page.TotalElements, int64(0))
}

// TestAgentClone verifies that a cloned agent has the same structure as the original.
func TestAgentClone(t *testing.T) {
	skipIfNotRunning(t)
	client := testutil.NewClient(apiURL(), authToken())

	fixture := testutil.NewAgentFixture()
	var created map[string]any
	status, err := client.Post("/api/agents", fixture, &created)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	id := created["id"].(string)

	var cloned map[string]any
	status, err = client.Post("/api/agents/"+id+"/clone", nil, &cloned)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.NotEqual(t, id, cloned["id"], "clone must have a different ID")
	assert.Equal(t, fixture.Name, cloned["name"])
}

// TestPipelineCRUD verifies create and read of a pipeline linked to an agent.
func TestPipelineCRUD(t *testing.T) {
	skipIfNotRunning(t)
	client := testutil.NewClient(apiURL(), authToken())

	// Create agent first.
	fixture := testutil.NewAgentFixture()
	var agent map[string]any
	status, err := client.Post("/api/agents", fixture, &agent)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	agentID := agent["id"].(string)

	// Create pipeline for agent.
	pipeline := map[string]any{
		"name":    "Test Pipeline",
		"agentId": agentID,
	}
	var createdPipeline map[string]any
	status, err = client.Post("/api/pipelines", pipeline, &createdPipeline)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
	assert.Equal(t, "Test Pipeline", createdPipeline["name"])
}
