package contracts

import (
	"net/http"
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractPipelinesList validates pipeline list endpoints.
func TestContractPipelinesList(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "pipelines_list", javaPath: "/api/pipelines?page=0&size=5", goPath: "/api/pipelines?page=0&size=5", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractLLMPresetsExtended validates LLM preset filter endpoints.
func TestContractLLMPresetsExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "llm_presets_by_provider_openai", javaPath: "/api/llm-config-presets/by-provider/openai?page=0&size=5", goPath: "/api/llm-config-presets/by-provider/openai?page=0&size=5", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractExecutionDetails validates execution hierarchy endpoint.
func TestContractExecutionDetails(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		// Non-existent ID validates the error contract for 404 responses.
		{name: "execution_details_not_found", javaPath: "/api/executions/00000000-0000-0000-0000-000000000000/details", goPath: "/api/executions/00000000-0000-0000-0000-000000000000/details", statusCode: http.StatusNotFound},
		{name: "execution_nodes_not_found", javaPath: "/api/executions/00000000-0000-0000-0000-000000000000/nodes", goPath: "/api/executions/00000000-0000-0000-0000-000000000000/nodes", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractWebhooksExtended validates webhook list + 404 contract.
func TestContractWebhooksExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "webhook_not_found", javaPath: "/api/webhooks/00000000-0000-0000-0000-000000000000", goPath: "/api/webhooks/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "webhook_delivery_logs_not_found", javaPath: "/api/webhooks/00000000-0000-0000-0000-000000000000/delivery-logs?page=0&size=5", goPath: "/api/webhooks/00000000-0000-0000-0000-000000000000/delivery-logs?page=0&size=5", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractKnowledgeBasesExtended validates knowledge base + document list.
func TestContractKnowledgeBasesExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "knowledgebase_not_found", javaPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000", goPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "documents_list_all", javaPath: "/api/documents?page=0&size=5", goPath: "/api/documents?page=0&size=5", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractExperimentsExtended validates experiment list and results.
func TestContractExperimentsExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "experiment_not_found", javaPath: "/api/experiments/00000000-0000-0000-0000-000000000000", goPath: "/api/experiments/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractMetricsExtended validates cost-breakdown endpoint.
func TestContractMetricsExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "metrics_cost_breakdown", javaPath: "/api/metrics/cost-breakdown", goPath: "/api/metrics/cost-breakdown", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractAgentSingleResource validates GET by ID returns a 404 contract
// (using a nil UUID that will never match a real resource).
func TestContractAgentSingleResource(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "agent_not_found", javaPath: "/api/agents/00000000-0000-0000-0000-000000000000", goPath: "/api/agents/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "skill_not_found", javaPath: "/api/skills/00000000-0000-0000-0000-000000000000", goPath: "/api/skills/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "tool_not_found", javaPath: "/api/tools/00000000-0000-0000-0000-000000000000", goPath: "/api/tools/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "llm_preset_not_found", javaPath: "/api/llm-config-presets/00000000-0000-0000-0000-000000000000", goPath: "/api/llm-config-presets/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "user_not_found", javaPath: "/api/users/nonexistent-user", goPath: "/api/users/nonexistent-user", statusCode: http.StatusNotFound},
		{name: "mcp_config_not_found", javaPath: "/api/mcp-server-configs/00000000-0000-0000-0000-000000000000", goPath: "/api/mcp-server-configs/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "datasource_not_found", javaPath: "/api/datasources/00000000-0000-0000-0000-000000000000", goPath: "/api/datasources/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "vpn_resource_not_found", javaPath: "/api/vpn-resources/00000000-0000-0000-0000-000000000000", goPath: "/api/vpn-resources/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractMemoryExtended validates memory list endpoint.
func TestContractMemoryExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "memory_agent_not_found", javaPath: "/api/memories/00000000-0000-0000-0000-000000000000?page=0&size=5", goPath: "/api/memories/00000000-0000-0000-0000-000000000000?page=0&size=5", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractChatExtended validates chat session list.
func TestContractChatExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "chat_sessions_list", javaPath: "/api/chat/sessions?page=0&size=5", goPath: "/api/chat/sessions?page=0&size=5", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractAuditLogExtended validates audit log 404 contract.
func TestContractAuditLogExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "audit_log_not_found", javaPath: "/api/audit-logs/00000000-0000-0000-0000-000000000000", goPath: "/api/audit-logs/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractTenantExtended validates tenant exists endpoint.
func TestContractTenantExtended(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "tenant_exists_nonexistent", javaPath: "/public/tenants/nonexistent-tenant-xyz/exists", goPath: "/public/tenants/nonexistent-tenant-xyz/exists", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}
