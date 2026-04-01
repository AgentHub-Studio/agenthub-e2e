package contracts

import (
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/snapshot"
)

// TestContractSkills validates skill and tool endpoints.
func TestContractSkills(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "skills_list", javaPath: "/api/skills?page=0&size=5", goPath: "/api/skills?page=0&size=5", statusCode: 200},
		{name: "tools_list", javaPath: "/api/tools?page=0&size=5", goPath: "/api/tools?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractLlmPresets validates LLM config preset endpoints.
func TestContractLlmPresets(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "llm_presets_list", javaPath: "/api/llm-config-presets?page=0&size=5", goPath: "/api/llm-config-presets?page=0&size=5", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}

// TestContractSettings validates tenant settings endpoints.
func TestContractSettings(t *testing.T) {
	skipIfNoServices(t)
	store := snapshot.DefaultStore()

	cases := []contractCase{
		{name: "settings_get", javaPath: "/api/settings", goPath: "/api/settings", statusCode: 200},
	}

	runContractCasesWithSnapshot(t, store, cases)
}
