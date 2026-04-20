package parser

import (
	"strings"
	"testing"
)

func TestParseLogfmtLines_ValidEntries(t *testing.T) {
	input := `time="2024-01-15T10:00:00Z" level=info msg="server started" port=8080
time="2024-01-15T10:01:00Z" level=warn msg="high memory" used=92%
`
	entries, err := ParseLogfmtLines(strings.NewReader(input), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0]["level"] != "info" {
		t.Errorf("expected level=info, got %v", entries[0]["level"])
	}
	if entries[1]["msg"] != "high memory" {
		t.Errorf("expected msg='high memory', got %v", entries[1]["msg"])
	}
}

func TestParseLogfmtLines_TimestampParsed(t *testing.T) {
	input := `time="2024-03-10T08:30:00Z" level=debug msg=ping
`
	entries, err := ParseLogfmtLines(strings.NewReader(input), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	// timestamp should be normalised to RFC3339Nano
	ts, ok := entries[0]["time"].(string)
	if !ok || ts == "" {
		t.Errorf("expected parsed timestamp string, got %v", entries[0]["time"])
	}
}

func TestParseLogfmtLines_EmptyInput(t *testing.T) {
	entries, err := ParseLogfmtLines(strings.NewReader(""), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseLogfmtLines_SkipsMalformedRows(t *testing.T) {
	input := `time="2024-01-01T00:00:00Z" level=info msg=ok
   
time="2024-01-01T00:01:00Z" level=error msg="bad \x00 bytes"
`
	entries, err := ParseLogfmtLines(strings.NewReader(input), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// at minimum the first valid line should be present
	if len(entries) < 1 {
		t.Errorf("expected at least 1 entry, got %d", len(entries))
	}
}

func TestParseLogfmtLines_DefaultTimestampField(t *testing.T) {
	input := `time="2024-06-01T12:00:00Z" level=info msg=hello
`
	// passing empty string should default to "time"
	entries, err := ParseLogfmtLines(strings.NewReader(input), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if _, ok := entries[0]["time"]; !ok {
		t.Error("expected 'time' field in entry")
	}
}

func TestParseLogfmtLines_BareKey(t *testing.T) {
	input := `time="2024-01-01T00:00:00Z" level=info verbose
`
	entries, err := ParseLogfmtLines(strings.NewReader(input), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0]["verbose"] != true {
		t.Errorf("expected bare key 'verbose'=true, got %v", entries[0]["verbose"])
	}
}
