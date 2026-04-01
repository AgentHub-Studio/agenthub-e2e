package contracts

import (
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractKnowledgeBases validates knowledge base and document endpoints.
func TestContractKnowledgeBases(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "knowledgebases_list", javaPath: "/api/knowledge-bases?page=0&size=5", goPath: "/api/knowledge-bases?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractDocuments validates document endpoints.
func TestContractDocuments(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		// Documents are scoped to a knowledge base; list all is not exposed — skip.
		{name: "documents_status_enum", javaPath: "/api/documents?status=PENDING&page=0&size=5", goPath: "/api/documents?status=PENDING&page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractMcp validates MCP server config endpoints.
func TestContractMcp(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "mcp_configs_list", javaPath: "/api/mcp-server-configs?page=0&size=5", goPath: "/api/mcp-server-configs?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}
