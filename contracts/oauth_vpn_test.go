package contracts

import (
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractOAuth validates OAuth credential endpoints (secrets must be masked).
func TestContractOAuth(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "oauth_credentials_list", javaPath: "/api/oauth-credentials?page=0&size=5", goPath: "/api/oauth-credentials?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractExperiments validates prompt experiment endpoints.
func TestContractExperiments(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "experiments_list", javaPath: "/api/experiments?page=0&size=5", goPath: "/api/experiments?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractVpnDatasource validates VPN resource and datasource endpoints.
func TestContractVpnDatasource(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "vpn_resources_list", javaPath: "/api/vpn-resources?page=0&size=5", goPath: "/api/vpn-resources?page=0&size=5", statusCode: 200},
		{name: "datasources_list", javaPath: "/api/datasources?page=0&size=5", goPath: "/api/datasources?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractSearch validates global search endpoints.
func TestContractSearch(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "search_agents", javaPath: "/api/search?q=test&type=AGENT", goPath: "/api/search?q=test&type=AGENT", statusCode: 200},
		{name: "search_all", javaPath: "/api/search?q=test", goPath: "/api/search?q=test", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractMemory validates agent memory endpoints.
func TestContractMemory(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "memory_list", javaPath: "/api/memories?page=0&size=5", goPath: "/api/memories?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}
