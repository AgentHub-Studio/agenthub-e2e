// Package snapshot provides file-based snapshot testing for REST contract tests.
// Snapshots are stored as JSON files under testdata/snapshots/.
//
// # Workflow
//
//  1. First run (UPDATE_SNAPSHOTS=true): captures live responses and writes snapshot files.
//  2. Subsequent runs: loads stored snapshots and compares new responses using the comparator.
//
// Activate via environment variables:
//
//	UPDATE_SNAPSHOTS=true  — overwrite existing snapshot files with live responses
//	CONTRACT_TESTS=1       — enable contract tests (otherwise tests are skipped)
package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AgentHub-Studio/agenthub-contracts/comparator"
)

const snapshotDir = "testdata/snapshots"

// Store manages snapshot files for a test suite.
type Store struct {
	dir    string
	update bool
}

// New creates a Store rooted at dir.
// When update is true the store writes (or overwrites) snapshots instead of comparing.
func New(dir string, update bool) *Store {
	return &Store{dir: dir, update: update}
}

// DefaultStore returns a Store using the conventional testdata/snapshots directory.
// Set UPDATE_SNAPSHOTS=true to enable update mode.
func DefaultStore() *Store {
	return New(snapshotDir, os.Getenv("UPDATE_SNAPSHOTS") == "true")
}

// Assert compares body against the stored snapshot for name.
// In update mode it writes body to disk instead and marks the test as passing.
// Dynamic fields listed in comparator.DefaultOptions are ignored.
func (s *Store) Assert(t *testing.T, name string, body []byte) {
	t.Helper()

	path := s.path(name)

	if s.update {
		if err := s.write(path, body); err != nil {
			t.Fatalf("snapshot: write %s: %v", path, err)
		}
		t.Logf("snapshot: updated %s", path)
		return
	}

	stored, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("snapshot: %s not found — run with UPDATE_SNAPSHOTS=true to create it", path)
		}
		t.Fatalf("snapshot: read %s: %v", path, err)
	}

	diffs := comparator.Diff(stored, body, comparator.DefaultOptions())
	if len(diffs) > 0 {
		t.Errorf("snapshot mismatch for %s:\n%s", name, strings.Join(diffs, "\n"))
	}
}

// AssertStructure compares only the JSON schema (field names + types) of body against the snapshot.
// Useful when values are volatile but structure must be stable.
func (s *Store) AssertStructure(t *testing.T, name string, body []byte) {
	t.Helper()

	path := s.path(name)

	schema := schemaOf(body)
	schemaBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		t.Fatalf("snapshot: marshal schema: %v", err)
	}

	if s.update {
		if err := s.write(path, schemaBytes); err != nil {
			t.Fatalf("snapshot: write schema %s: %v", path, err)
		}
		t.Logf("snapshot: updated schema %s", path)
		return
	}

	stored, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("snapshot: %s not found — run with UPDATE_SNAPSHOTS=true to create it", path)
		}
		t.Fatalf("snapshot: read %s: %v", path, err)
	}

	diffs := comparator.Diff(stored, schemaBytes, comparator.Options{})
	if len(diffs) > 0 {
		t.Errorf("schema mismatch for %s:\n%s", name, strings.Join(diffs, "\n"))
	}
}

func (s *Store) path(name string) string {
	safe := strings.ReplaceAll(name, "/", "_")
	safe = strings.ReplaceAll(safe, " ", "_")
	return filepath.Join(s.dir, safe+".json")
}

func (s *Store) write(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Pretty-print if possible.
	var v any
	if err := json.Unmarshal(data, &v); err == nil {
		if pretty, err := json.MarshalIndent(v, "", "  "); err == nil {
			data = pretty
		}
	}
	return os.WriteFile(path, data, 0o644)
}

// schemaOf converts a JSON value to a schema representation:
// objects become {"field": <type>}, arrays become [<element-schema>],
// scalars become their type name ("string", "number", "bool", "null").
func schemaOf(data []byte) any {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Sprintf("invalid JSON: %v", err)
	}
	return schemaValue(v)
}

func schemaValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, elem := range val {
			out[k] = schemaValue(elem)
		}
		return out
	case []any:
		if len(val) == 0 {
			return []any{"<empty>"}
		}
		return []any{schemaValue(val[0])}
	case string:
		return "string"
	case float64:
		return "number"
	case bool:
		return "bool"
	case nil:
		return "null"
	default:
		return fmt.Sprintf("unknown(%T)", v)
	}
}
