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

// TestKnowledgeBase_CRUD verifies create, get, list, and delete of knowledge bases.
func TestKnowledgeBase_CRUD(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	// Create
	req := map[string]any{
		"name":        fmt.Sprintf("kb-test-%s", uuid.New().String()[:8]),
		"description": "E2E test knowledge base",
	}
	var created map[string]any
	status, err := client.Post("/api/knowledge-bases", req, &created)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	require.NotEmpty(t, created["id"])

	kbID := created["id"].(string)

	t.Cleanup(func() {
		_, _ = client.Delete("/api/knowledge-bases/" + kbID)
	})

	// Get by ID
	var fetched map[string]any
	status, err = client.Get("/api/knowledge-bases/"+kbID, &fetched)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, kbID, fetched["id"])
	assert.Equal(t, req["name"], fetched["name"])

	// List — knowledge base must appear
	var page testutil.Page[map[string]any]
	status, err = client.Get("/api/knowledge-bases?page=0&size=20", &page)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	found := false
	for _, item := range page.Content {
		if item["id"] == kbID {
			found = true
			break
		}
	}
	assert.True(t, found, "created knowledge base must appear in list")

	// Delete
	status, err = client.Delete("/api/knowledge-bases/" + kbID)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, status)

	// Verify deletion
	var after map[string]any
	status, err = client.Get("/api/knowledge-bases/"+kbID, &after)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

// TestKnowledgeBase_NotFound verifies 404 for a non-existent knowledge base.
func TestKnowledgeBase_NotFound(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	var result map[string]any
	status, err := client.Get("/api/knowledge-bases/"+uuid.New().String(), &result)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}

// TestKnowledgeBase_DocumentList verifies listing documents within a knowledge base.
func TestKnowledgeBase_DocumentList(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	// Create a KB first.
	req := map[string]any{
		"name": fmt.Sprintf("kb-docs-%s", uuid.New().String()[:8]),
	}
	var kb map[string]any
	status, err := client.Post("/api/knowledge-bases", req, &kb)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	kbID := kb["id"].(string)

	t.Cleanup(func() { _, _ = client.Delete("/api/knowledge-bases/" + kbID) })

	// List documents — expect empty page.
	var page testutil.Page[map[string]any]
	status, err = client.Get(fmt.Sprintf("/api/knowledge-bases/%s/documents?page=0&size=10", kbID), &page)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, int64(0), page.TotalElements)
}

// TestKnowledgeBase_Pagination verifies Page[T] structure for knowledge base list.
func TestKnowledgeBase_Pagination(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	var page testutil.Page[map[string]any]
	status, err := client.Get("/api/knowledge-bases?page=0&size=5", &page)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, status)

	assert.GreaterOrEqual(t, page.TotalElements, int64(0))
	assert.Equal(t, 5, page.PageSize)
	assert.Equal(t, 0, page.PageNumber)
}
