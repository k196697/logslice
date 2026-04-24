package aggregate

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
)

func makeEntry(fields map[string]interface{}) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

func TestRun_GroupsByField(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(map[string]interface{}{"level": "info", "msg": "a"}),
		makeEntry(map[string]interface{}{"level": "error", "msg": "b"}),
		makeEntry(map[string]interface{}{"level": "info", "msg": "c"}),
	}
	opts := Options{GroupBy: "level", CountField: "count"}
	result, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result))
	}
	// Groups are sorted alphabetically: error, info
	if result[0].Fields["level"] != "error" {
		t.Errorf("expected first group 'error', got %v", result[0].Fields["level"])
	}
	if result[0].Fields["count"] != 1 {
		t.Errorf("expected count 1 for error, got %v", result[0].Fields["count"])
	}
	if result[1].Fields["count"] != 2 {
		t.Errorf("expected count 2 for info, got %v", result[1].Fields["count"])
	}
}

func TestRun_EmptyGroupByReturnsError(t *testing.T) {
	_, err := Run([]filter.Entry{makeEntry(map[string]interface{}{"x": "1"})}, Options{})
	if err == nil {
		t.Fatal("expected error for empty GroupBy")
	}
}

func TestRun_MissingFieldGroupedAsEmpty(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(map[string]interface{}{"msg": "no level"}),
		makeEntry(map[string]interface{}{"level": "warn"}),
	}
	opts := Options{GroupBy: "level", CountField: "count"}
	result, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect two groups: "" and "warn"
	if len(result) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result))
	}
}

func TestRun_EmptyInput(t *testing.T) {
	opts := Options{GroupBy: "level", CountField: "count"}
	result, err := Run([]filter.Entry{}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}

func TestRun_DefaultCountField(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(map[string]interface{}{"level": "debug"}),
	}
	// CountField left empty — should default to "count"
	opts := Options{GroupBy: "level"}
	result, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result[0].Fields["count"]; !ok {
		t.Errorf("expected 'count' field in result entry")
	}
}
