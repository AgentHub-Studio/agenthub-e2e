package contracts

import (
	"net/http"
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractAgentMemory validates agent memory CRUD and recall endpoints.
func TestContractAgentMemory(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "memory_list_empty", javaPath: "/api/agents/00000000-0000-0000-0000-000000000000/memory?page=0&size=5", goPath: "/api/agents/00000000-0000-0000-0000-000000000000/memory?page=0&size=5", statusCode: 200},
		{name: "memory_get_key_not_found", javaPath: "/api/agents/00000000-0000-0000-0000-000000000000/memory/nonexistent-key", goPath: "/api/agents/00000000-0000-0000-0000-000000000000/memory/nonexistent-key", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractAgentStateTransitions validates that agent publish/archive/clone action
// routes exist and reject wrong-method GET with 405 consistently between Java and Go.
func TestContractAgentStateTransitions(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	// These are POST-only endpoints; GET validates the 405 contract.
	cases := []contractCase{
		{name: "agent_publish_method_not_allowed", javaPath: "/api/agents/00000000-0000-0000-0000-000000000000/publish", goPath: "/api/agents/00000000-0000-0000-0000-000000000000/publish", statusCode: http.StatusMethodNotAllowed},
		{name: "agent_archive_method_not_allowed", javaPath: "/api/agents/00000000-0000-0000-0000-000000000000/archive", goPath: "/api/agents/00000000-0000-0000-0000-000000000000/archive", statusCode: http.StatusMethodNotAllowed},
		{name: "agent_clone_method_not_allowed", javaPath: "/api/agents/00000000-0000-0000-0000-000000000000/clone", goPath: "/api/agents/00000000-0000-0000-0000-000000000000/clone", statusCode: http.StatusMethodNotAllowed},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractChatSessionDetail validates chat session detail and message endpoints.
func TestContractChatSessionDetail(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "chat_session_not_found", javaPath: "/api/chat/sessions/00000000-0000-0000-0000-000000000000", goPath: "/api/chat/sessions/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "chat_session_messages_not_found", javaPath: "/api/chat/sessions/00000000-0000-0000-0000-000000000000/messages?page=0&size=5", goPath: "/api/chat/sessions/00000000-0000-0000-0000-000000000000/messages?page=0&size=5", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractExecutionToolTraces validates execution node tool trace endpoints.
func TestContractExecutionToolTraces(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "execution_node_tools_not_found", javaPath: "/api/executions/00000000-0000-0000-0000-000000000000/nodes/00000000-0000-0000-0000-000000000001/tools", goPath: "/api/executions/00000000-0000-0000-0000-000000000000/nodes/00000000-0000-0000-0000-000000000001/tools", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractKnowledgeBaseTransitions validates KB activate/pause and document list endpoints.
func TestContractKnowledgeBaseTransitions(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		// POST-only endpoints: GET validates the 405 contract.
		{name: "kb_activate_method_not_allowed", javaPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000/activate", goPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000/activate", statusCode: http.StatusMethodNotAllowed},
		{name: "kb_pause_method_not_allowed", javaPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000/pause", goPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000/pause", statusCode: http.StatusMethodNotAllowed},
		{name: "kb_documents_list_not_found", javaPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000/documents?page=0&size=5", goPath: "/api/knowledge-bases/00000000-0000-0000-0000-000000000000/documents?page=0&size=5", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractMCPServers validates MCP server config CRUD endpoints.
func TestContractMCPServers(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "mcp_servers_list", javaPath: "/api/mcp/servers?page=0&size=5", goPath: "/api/mcp/servers?page=0&size=5", statusCode: 200},
		{name: "mcp_server_not_found", javaPath: "/api/mcp/servers/00000000-0000-0000-0000-000000000000", goPath: "/api/mcp/servers/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractPipelineDetail validates pipeline get-by-ID endpoint.
func TestContractPipelineDetail(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "pipeline_not_found", javaPath: "/api/pipelines/00000000-0000-0000-0000-000000000000", goPath: "/api/pipelines/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractSettings validates settings list and get-by-key endpoints.
func TestContractSettings(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "settings_list", javaPath: "/api/settings", goPath: "/api/settings", statusCode: 200},
		{name: "settings_key_not_found", javaPath: "/api/settings/nonexistent.key", goPath: "/api/settings/nonexistent.key", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractSkillToolBindings validates skill-to-tool binding endpoints.
func TestContractSkillToolBindings(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "skill_tools_list_not_found", javaPath: "/api/skills/00000000-0000-0000-0000-000000000000/tools", goPath: "/api/skills/00000000-0000-0000-0000-000000000000/tools", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractWebhookDeliveries validates webhook delivery log and test endpoints.
func TestContractWebhookDeliveries(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "webhook_deliveries_not_found", javaPath: "/api/webhooks/00000000-0000-0000-0000-000000000000/deliveries?page=0&size=5", goPath: "/api/webhooks/00000000-0000-0000-0000-000000000000/deliveries?page=0&size=5", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractPackageRegistry validates public package registry endpoints.
func TestContractPackageRegistry(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "packages_list_public", javaPath: "/api/packages?page=0&size=5", goPath: "/api/packages?page=0&size=5", statusCode: 200},
		{name: "package_not_found", javaPath: "/api/packages/00000000-0000-0000-0000-000000000000", goPath: "/api/packages/00000000-0000-0000-0000-000000000000", statusCode: http.StatusNotFound},
		{name: "package_slug_not_found", javaPath: "/api/packages/slug/nonexistent-package-slug", goPath: "/api/packages/slug/nonexistent-package-slug", statusCode: http.StatusNotFound},
	}
	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractPackageVersions validates package version endpoints.
func TestContractPackageVersions(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "package_versions_list_not_found", javaPath: "/api/packages/00000000-0000-0000-0000-000000000000/versions", goPath: "/api/packages/00000000-0000-0000-0000-000000000000/versions", statusCode: 200},
	}
	runContractCasesWithSnapshot(t, store, cases)
}
