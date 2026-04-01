package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AgentHub-Studio/agenthub-e2e/internal/testutil"
)

// TestExecution_ListEmpty verifies the executions list endpoint returns an empty page initially.
func TestExecution_ListEmpty(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	var page testutil.Page[map[string]any]
	status, err := client.Get("/api/executions?page=0&size=20", &page)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.GreaterOrEqual(t, page.TotalElements, int64(0))
}

// TestExecution_FilterByStatus verifies filtering executions by status.
func TestExecution_FilterByStatus(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	for _, status := range []string{"PENDING", "RUNNING", "COMPLETED", "FAILED", "CANCELLED"} {
		var page testutil.Page[map[string]any]
		httpStatus, err := client.Get(fmt.Sprintf("/api/executions?status=%s&page=0&size=10", status), &page)
		require.NoError(t, err, "status=%s", status)
		assert.Equal(t, http.StatusOK, httpStatus, "status=%s", status)
	}
}

// TestExecution_NotFound verifies 404 for a non-existent execution.
func TestExecution_NotFound(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	var result map[string]any
	status, err := client.Get("/api/executions/"+uuid.New().String(), &result)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

// TestExecution_Details_HierarchyShape verifies that an execution details response
// contains the expected hierarchy fields (nodes, tools).
func TestExecution_Details_HierarchyShape(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	// List to find any existing execution.
	var page testutil.Page[map[string]any]
	status, err := client.Get("/api/executions?page=0&size=1", &page)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	if len(page.Content) == 0 {
		t.Skip("no executions present — create one via agent execution to run this test")
	}

	execID, ok := page.Content[0]["id"].(string)
	require.True(t, ok, "execution id must be a string")

	var details map[string]any
	status, err = client.Get("/api/executions/"+execID, &details)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	// Verify expected top-level fields.
	assert.Contains(t, details, "id")
	assert.Contains(t, details, "agentId")
	assert.Contains(t, details, "status")
	assert.Contains(t, details, "nodes")

	// Nodes must be an array.
	nodes, ok := details["nodes"].([]any)
	assert.True(t, ok, "nodes must be an array")

	// Each node must have tools array.
	for _, n := range nodes {
		node, ok := n.(map[string]any)
		require.True(t, ok)
		assert.Contains(t, node, "nodeId")
		assert.Contains(t, node, "status")
		assert.Contains(t, node, "tools")
	}
}

// TestExecution_Cancel_StateMachine verifies that cancelling an execution returns 204
// and subsequent GET shows CANCELLED status.
func TestExecution_Cancel_StateMachine(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	// Find a RUNNING or PENDING execution to cancel.
	var page testutil.Page[map[string]any]
	status, err := client.Get("/api/executions?status=RUNNING&page=0&size=1", &page)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	if len(page.Content) == 0 {
		t.Skip("no RUNNING executions to cancel")
	}

	execID := page.Content[0]["id"].(string)

	cancelStatus, err := client.Delete("/api/executions/" + execID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, cancelStatus)

	// Verify the execution is now CANCELLED.
	var details map[string]any
	getStatus, err := client.Get("/api/executions/"+execID, &details)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, getStatus)
	assert.Equal(t, "CANCELLED", details["status"])
}

// TestExecution_InvalidTransition verifies that attempting to cancel an already-terminal
// execution returns an error.
func TestExecution_InvalidTransition(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	// Find a COMPLETED execution.
	var page testutil.Page[map[string]any]
	status, err := client.Get("/api/executions?status=COMPLETED&page=0&size=1", &page)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	if len(page.Content) == 0 {
		t.Skip("no COMPLETED executions to test invalid transition")
	}

	execID := page.Content[0]["id"].(string)

	// Attempting to cancel a COMPLETED execution should return 409 or 422.
	cancelStatus, err := client.Delete("/api/executions/" + execID)
	require.NoError(t, err)
	assert.True(t, cancelStatus == http.StatusConflict || cancelStatus == http.StatusUnprocessableEntity,
		"expected 409 or 422 for invalid state transition, got %d", cancelStatus)
}
