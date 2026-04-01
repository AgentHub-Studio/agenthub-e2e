package contracts

import (
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractTenants validates tenant provisioning and listing endpoints.
func TestContractTenants(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "tenants_list", javaPath: "/public/tenants", goPath: "/public/tenants", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractUsers validates user management endpoints.
func TestContractUsers(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "users_list", javaPath: "/api/users?page=0&size=5", goPath: "/api/users?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}
