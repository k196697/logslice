package parser

import (
	"strings"
	"testing"
	"time"
)

const csvInput = `time,level,message
2024-03-01T10:00:00Z,info,service started
2024-03-01T10:01:00Z,warn,high memory usage
2024-03-01T10:02:00Z,error,connection refused
`

func TestParseCSVLines_ValidEntries(t *testing.T) {
	entries, err := ParseCSVLines(strings.NewReader(csvInput), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", entries[0].Fields["level"])
	}
	if entries[2].Fields["message"] != "connection refused" {
		t.Errorf("unexpected message: %v", entries[2].Fields["message"])
	}
}

func TestParseCSVLines_TimestampParsed(t *testing.T) {
	entries, err := ParseCSVLines(strings.NewReader(csvInput), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected, _ := time.Parse(time.RFC3339, "2024-03-01T10:01:00Z")
	if !entries[1].Timestamp.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, entries[1].Timestamp)
	}
}

func TestParseCSVLines_EmptyInput(t *testing.T) {
	entries, err := ParseCSVLines(strings.NewReader(""), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestParseCSVLines_MissingTimestampField(t *testing.T) {
	input := "level,message\ninfo,hello\n"
	entries, err := ParseCSVLines(strings.NewReader(input), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if !entries[0].Timestamp.IsZero() {
		t.Errorf("expected zero timestamp when field absent")
	}
}

func TestParseCSVLines_SkipsMalformedRows(t *testing.T) {
	input := "time,level,message\n2024-03-01T10:00:00Z,info\n2024-03-01T10:01:00Z,warn,ok\n"
	entries, err := ParseCSVLines(strings.NewReader(input), "time")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("expected 1 valid entry, got %d", len(entries))
	}
}

func TestParseTimestampString_Formats(t *testing.T) {
	cases := []string{
		"2024-03-01T10:00:00Z",
		"2024-03-01 10:00:00",
		"2024/03/01 10:00:00",
	}
	for _, c := range cases {
		_, err := parseTimestampString(c)
		if err != nil {
			t.Errorf("failed to parse %q: %v", c, err)
		}
	}
}
