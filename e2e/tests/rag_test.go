package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AgentHub-Studio/agenthub-e2e/internal/testutil"
)

// TestRAG_DocumentUploadAndIndex verifies the RAG pipeline:
// knowledge base creation → document upload → document appears in list.
//
// Full embedding and vector search require the agenthub-embedding service,
// which is not started by this suite.  The test validates the storage and
// status-tracking layer only.
func TestRAG_DocumentUploadAndIndex(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	// Create knowledge base.
	kbReq := map[string]any{
		"name":        fmt.Sprintf("rag-test-%s", uuid.New().String()[:8]),
		"description": "E2E RAG test knowledge base",
	}
	var kb map[string]any
	status, err := client.Post("/api/knowledge-bases", kbReq, &kb)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	kbID := kb["id"].(string)
	t.Cleanup(func() { _, _ = client.Delete("/api/knowledge-bases/" + kbID) })

	// Upload a minimal plain-text document.
	docContent := []byte("AgentHub is a multi-tenant AI orchestration platform.")
	var doc map[string]any
	status, err = client.PostMultipart(
		fmt.Sprintf("/api/knowledge-bases/%s/documents", kbID),
		"file", "test-doc.txt", docContent,
		map[string]string{"name": "Test Document"},
		&doc,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status, "document upload should return 201")
	require.NotEmpty(t, doc["id"], "response must contain document id")

	docID := doc["id"].(string)
	t.Cleanup(func() {
		_, _ = client.Delete(fmt.Sprintf("/api/knowledge-bases/%s/documents/%s", kbID, docID))
	})

	// Document must have an initial status (PENDING or EXTRACTING).
	initialStatus, _ := doc["status"].(string)
	assert.NotEmpty(t, initialStatus, "document must have a status")
	validStatuses := []string{"PENDING", "EXTRACTING", "CHUNKING", "EMBEDDING", "INDEXED"}
	assert.Contains(t, validStatuses, strings.ToUpper(initialStatus),
		"document status %q must be one of %v", initialStatus, validStatuses)

	// Document must appear in the list.
	var page testutil.Page[map[string]any]
	status, err = client.Get(
		fmt.Sprintf("/api/knowledge-bases/%s/documents?page=0&size=20", kbID), &page)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.GreaterOrEqual(t, page.TotalElements, int64(1), "uploaded document must appear in list")

	found := false
	for _, item := range page.Content {
		if item["id"] == docID {
			found = true
			break
		}
	}
	assert.True(t, found, "uploaded document must be retrievable by id in list")

	// Document must be retrievable by ID.
	var fetchedDoc map[string]any
	status, err = client.Get(
		fmt.Sprintf("/api/knowledge-bases/%s/documents/%s", kbID, docID), &fetchedDoc)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, docID, fetchedDoc["id"])
}

// TestRAG_KnowledgeBase_NotFound verifies that requesting a non-existent document returns 404.
func TestRAG_Document_NotFound(t *testing.T) {
	skipIfNotRunning(t)

	client := testutil.NewClient(apiURL(), authToken())

	// Create KB.
	var kb map[string]any
	status, err := client.Post("/api/knowledge-bases",
		map[string]any{"name": "kb-notfound-" + uuid.New().String()[:8]}, &kb)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, status)
	kbID := kb["id"].(string)
	t.Cleanup(func() { _, _ = client.Delete("/api/knowledge-bases/" + kbID) })

	// Non-existent document.
	var result map[string]any
	status, err = client.Get(
		fmt.Sprintf("/api/knowledge-bases/%s/documents/%s", kbID, uuid.New().String()), &result)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, status)
}
