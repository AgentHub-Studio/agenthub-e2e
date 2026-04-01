package tests

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AgentHub-Studio/agenthub-e2e/internal/testutil"
)

func observabilityURL() string {
	if v := os.Getenv("OBSERVABILITY_URL"); v != "" {
		return v
	}
	return "http://localhost:8086"
}

// TestObservability_MetricRecord verifies recording a metric event.
func TestObservability_MetricRecord(t *testing.T) {
	skipIfNotRunning(t)
	if os.Getenv("E2E_OBSERVABILITY") == "" {
		t.Skip("set E2E_OBSERVABILITY=1 to run (requires agenthub-observability running)")
	}

	tenantID := os.Getenv("TEST_TENANT_ID")
	if tenantID == "" {
		t.Skip("set TEST_TENANT_ID")
	}

	client := testutil.NewClient(observabilityURL(), authToken())

	event := map[string]any{
		"tenantId":   tenantID,
		"metricName": "test.counter",
		"metricType": "counter",
		"value":      1.0,
	}
	var result map[string]any
	status, err := client.Post("/api/v1/metrics/events", event, &result)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, status)
}

// TestObservability_TracesQuery verifies listing execution traces.
func TestObservability_TracesQuery(t *testing.T) {
	skipIfNotRunning(t)
	if os.Getenv("E2E_OBSERVABILITY") == "" {
		t.Skip("set E2E_OBSERVABILITY=1 to run")
	}

	tenantID := os.Getenv("TEST_TENANT_ID")
	if tenantID == "" {
		t.Skip("set TEST_TENANT_ID")
	}

	client := testutil.NewClient(observabilityURL(), "")

	var result any
	status, err := client.Get("/api/v1/traces/executions?tenantId="+tenantID+"&limit=10", &result)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
}

// TestObservability_MetricsQuery verifies listing metric events.
func TestObservability_MetricsQuery(t *testing.T) {
	skipIfNotRunning(t)
	if os.Getenv("E2E_OBSERVABILITY") == "" {
		t.Skip("set E2E_OBSERVABILITY=1 to run")
	}

	tenantID := os.Getenv("TEST_TENANT_ID")
	if tenantID == "" {
		t.Skip("set TEST_TENANT_ID")
	}

	client := testutil.NewClient(observabilityURL(), "")

	var result any
	status, err := client.Get("/api/v1/metrics/events?tenantId="+tenantID+"&limit=10", &result)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, status)
}
