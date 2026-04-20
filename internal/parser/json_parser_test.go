package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParseJSONLines_ValidEntries(t *testing.T) {
	input := `{"time":"2024-03-01T10:00:00Z","level":"info","msg":"started"}
{"ts":"2024-03-01T10:01:00Z","level":"error","msg":"failed"}
`
	entries, err := ParseJSONLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	expected0, _ := time.Parse(time.RFC3339, "2024-03-01T10:00:00Z")
	if !entries[0].Timestamp.Equal(expected0) {
		t.Errorf("entry 0 timestamp: got %v, want %v", entries[0].Timestamp, expected0)
	}
	if entries[0].Fields["level"] != "info" {
		t.Errorf("entry 0 level: got %v, want info", entries[0].Fields["level"])
	}
}

func TestParseJSONLines_SkipsInvalidLines(t *testing.T) {
	input := `{"time":"2024-03-01T12:00:00Z","msg":"ok"}
not json at all
{"time":"2024-03-01T12:01:00Z","msg":"alsotentries, err := ParseJSONLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries after skipping bad line, got %d", len(entries))
	}
}

func TestParseJSONLines_UnixTimestamp(t *testing.T) {
	input := `{"ts":1709290800,"msg":"unix time"}
`
	entries, err := ParseJSONLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	expected := time.Unix(1709290800, 0).UTC()
	if !entries[0].Timestamp.Equal(expected) {
		t.Errorf("unix ts: got %v, want %v", entries[0].Timestamp, expected)
	}
}

func TestParseJSONLines_NoTimestampField(t *testing.T) {
	input := `{"level":"debug","msg":"no time here"}
`
	entries, err := ParseJSONLines(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !entries[0].Timestamp.IsZero() {
		t.Errorf("expected zero timestamp when no time field present")
	}
}

func TestParseJSONLines_EmptyInput(t *testing.T) {
	entries, err := ParseJSONLines(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries for empty input, got %d", len(entries))
	}
}
