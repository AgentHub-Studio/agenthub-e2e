package testutil

import "github.com/google/uuid"

// TenantFixture holds test tenant data.
type TenantFixture struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// NewTenantFixture returns a unique test tenant fixture.
func NewTenantFixture() TenantFixture {
	id := "test-" + uuid.New().String()[:8]
	return TenantFixture{
		ID:   id,
		Name: "Test Tenant " + id,
	}
}

// AgentFixture holds test agent data.
type AgentFixture struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	TenantID    string `json:"tenantId,omitempty"`
}

// NewAgentFixture returns a unique test agent fixture.
func NewAgentFixture() AgentFixture {
	return AgentFixture{
		Name:        "Test Agent " + uuid.New().String()[:8],
		Description: "Created by E2E tests",
	}
}

// SkillFixture holds test skill data.
type SkillFixture struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// NewSkillFixture returns a unique test skill fixture.
func NewSkillFixture() SkillFixture {
	suffix := uuid.New().String()[:8]
	return SkillFixture{
		Name:        "Test Skill " + suffix,
		Slug:        "test-skill-" + suffix,
		Description: "Created by E2E tests",
	}
}
