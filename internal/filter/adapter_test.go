package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func TestFromParsed_ConvertsEntries(t *testing.T) {
	now := time.Now().UTC()
	parsed := []filter.ParsedEntry{
		{
			Timestamp: now,
			Fields:    map[string]interface{}{"level": "info"},
			Raw:       `{"level":"info"}`,
		},
	}

	entries := filter.FromParsed(parsed)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !entries[0].Timestamp.Equal(now) {
		t.Error("timestamp mismatch after conversion")
	}
	if entries[0].Fields["level"] != "info" {
		t.Error("field mismatch after conversion")
	}
	if entries[0].Raw != `{"level":"info"}` {
		t.Error("raw mismatch after conversion")
	}
}

func TestToParsed_ConvertsBack(t *testing.T) {
	now := time.Now().UTC()
	entries := []filter.LogEntry{
		{
			Timestamp: now,
			Fields:    map[string]interface{}{"svc": "api"},
			Raw:       `{"svc":"api"}`,
		},
	}

	parsed := filter.ToParsed(entries)
	if len(parsed) != 1 {
		t.Fatalf("expected 1 parsed entry, got %d", len(parsed))
	}
	if !parsed[0].Timestamp.Equal(now) {
		t.Error("timestamp mismatch after reverse conversion")
	}
	if parsed[0].Fields["svc"] != "api" {
		t.Error("field mismatch after reverse conversion")
	}
}

func TestFromParsed_EmptySlice(t *testing.T) {
	result := filter.FromParsed([]filter.ParsedEntry{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
