package contracts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractAgents validates all agent + pipeline endpoints.
func TestContractAgents(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "agents_list", javaPath: "/api/agents?page=0&size=5", goPath: "/api/agents?page=0&size=5", statusCode: 200},
		{name: "agents_list_by_status", javaPath: "/api/agents?status=PUBLISHED&page=0&size=5", goPath: "/api/agents?status=PUBLISHED&page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractPipelines validates pipeline endpoints.
func TestContractPipelines(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "pipeline_node_types", javaPath: "/api/pipelines/node-types", goPath: "/api/pipelines/node-types", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractAgentPagination validates Page[T] structure matches between Java and Go.
func TestContractAgentPagination(t *testing.T) {
	skipIfNoServices(t)

	javaClient := newJavaClient()
	goClient := newGoClient()

	javaBody, javaStatus, err := javaClient.Get("/api/agents?page=0&size=20")
	require.NoError(t, err)
	assert.Equal(t, 200, javaStatus)

	goBody, goStatus, err := goClient.Get("/api/agents?page=0&size=20")
	require.NoError(t, err)
	assert.Equal(t, 200, goStatus)

	assertPaginationStructure(t, javaBody, goBody)
}
