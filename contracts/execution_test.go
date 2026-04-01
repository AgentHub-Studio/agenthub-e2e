package contracts

import (
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractExecutions validates execution tracking endpoints.
func TestContractExecutions(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "executions_list", javaPath: "/api/executions?page=0&size=5", goPath: "/api/executions?page=0&size=5", statusCode: 200},
		{name: "executions_filter_status", javaPath: "/api/executions?status=COMPLETED&page=0&size=5", goPath: "/api/executions?status=COMPLETED&page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractWebhooks validates webhook config and delivery log endpoints.
func TestContractWebhooks(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "webhooks_list", javaPath: "/api/webhooks?page=0&size=5", goPath: "/api/webhooks?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractAuditLogs validates audit log endpoints.
func TestContractAuditLogs(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "audit_logs_list", javaPath: "/api/audit-logs?page=0&size=5", goPath: "/api/audit-logs?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractMetrics validates agent metrics endpoints.
func TestContractMetrics(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "metrics_tenant_summary", javaPath: "/api/metrics/summary", goPath: "/api/metrics/summary", statusCode: 200},
		{name: "metrics_top_agents", javaPath: "/api/metrics/top-agents?limit=5", goPath: "/api/metrics/top-agents?limit=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}
