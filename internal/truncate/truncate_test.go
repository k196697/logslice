package truncate

import (
	"testing"
)

func makeEntry(fields map[string]string) Entry {
	return Entry{Fields: fields}
}

func TestRun_TruncatesAllFields(t *testing.T) {
	entries := []Entry{
		makeEntry(map[string]string{"msg": "hello world", "level": "info"}),
	}
	opts := Options{MaxLength: 5, Suffix: "..."}
	out, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out[0].Fields["msg"]; got != "he..." {
		t.Errorf("msg: got %q, want %q", got, "he...")
	}
	if got := out[0].Fields["level"]; got != "info" {
		t.Errorf("level should not be truncated, got %q", got)
	}
}

func TestRun_TruncatesOnlySpecifiedFields(t *testing.T) {
	entries := []Entry{
		makeEntry(map[string]string{"msg": "hello world", "path": "/very/long/path/value"}),
	}
	opts := Options{Fields: []string{"msg"}, MaxLength: 7, Suffix: "~"}
	out, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out[0].Fields["msg"]; got != "hello w~" {
		t.Errorf("msg: got %q, want %q", got, "hello w~")
	}
	if got := out[0].Fields["path"]; got != "/very/long/path/value" {
		t.Errorf("path should be unchanged, got %q", got)
	}
}

func TestRun_ShortValuesUnchanged(t *testing.T) {
	entries := []Entry{
		makeEntry(map[string]string{"msg": "hi"}),
	}
	opts := Options{MaxLength: 100, Suffix: "..."}
	out, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out[0].Fields["msg"]; got != "hi" {
		t.Errorf("expected unchanged, got %q", got)
	}
}

func TestRun_InvalidMaxLength(t *testing.T) {
	_, err := Run([]Entry{}, Options{MaxLength: 0, Suffix: "..."})
	if err == nil {
		t.Fatal("expected error for MaxLength=0")
	}
}

func TestRun_EmptyEntries(t *testing.T) {
	out, err := Run([]Entry{}, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}

func TestParseFields_CommaSeparated(t *testing.T) {
	fields := ParseFields("msg, level ,path")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[1] != "level" {
		t.Errorf("expected 'level', got %q", fields[1])
	}
}

func TestParseFields_Empty(t *testing.T) {
	if fields := ParseFields(""); fields != nil {
		t.Errorf("expected nil for empty input, got %v", fields)
	}
}
