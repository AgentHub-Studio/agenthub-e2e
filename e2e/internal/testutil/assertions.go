package testutil

import (
	"encoding/json"
	"fmt"
	"testing"
)

// AssertPage validates that a Page response has the expected total count.
func AssertPage[T any](t *testing.T, page Page[T], wantTotal int64) {
	t.Helper()
	if page.TotalElements != wantTotal {
		t.Errorf("page.TotalElements = %d, want %d", page.TotalElements, wantTotal)
	}
}

// AssertStatusCode validates an HTTP status code.
func AssertStatusCode(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("HTTP status = %d, want %d", got, want)
	}
}

// AssertField extracts a string field from a JSON response map.
func AssertField(t *testing.T, body []byte, field, want string) {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("AssertField: cannot parse JSON: %v", err)
	}
	got, ok := m[field]
	if !ok {
		t.Errorf("field %q not found in response", field)
		return
	}
	if fmt.Sprintf("%v", got) != want {
		t.Errorf("field %q = %v, want %q", field, got, want)
	}
}
