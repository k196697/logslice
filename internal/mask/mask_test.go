package mask

import (
	"testing"
)

func makeEntry(fields map[string]string) Entry {
	return Entry{Fields: fields}
}

func TestRun_MasksEntireField(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"password": "secret123"})}
	result, err := Run(entries, Options{Fields: []string{"password"}, Char: "*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result[0].Fields["password"]; got != "*********" {
		t.Errorf("expected *********, got %q", got)
	}
}

func TestRun_KeepsPrefixAndSuffix(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"token": "abcdef"})}
	opts := Options{Fields: []string{"token"}, Char: "*", KeepPrefix: 2, KeepSuffix: 1}
	result, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := result[0].Fields["token"]; got != "ab***f" {
		t.Errorf("expected ab***f, got %q", got)
	}
}

func TestRun_SkipsMissingField(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"user": "alice"})}
	opts := Options{Fields: []string{"password"}, Char: "*"}
	result, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result[0].Fields["password"]; ok {
		t.Error("expected missing field to remain absent")
	}
	if result[0].Fields["user"] != "alice" {
		t.Error("expected untouched field to be preserved")
	}
}

func TestRun_NoFieldsReturnsUnchanged(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"key": "value"})}
	result, err := Run(entries, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result[0].Fields["key"] != "value" {
		t.Error("expected entry to be unchanged")
	}
}

func TestRun_PrefixSufixOverlapLeavesUnchanged(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"pin": "1234"})}
	opts := Options{Fields: []string{"pin"}, Char: "*", KeepPrefix: 3, KeepSuffix: 2}
	result, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// pre+suf=5 >= len("1234")=4 → unchanged
	if got := result[0].Fields["pin"]; got != "1234" {
		t.Errorf("expected 1234 unchanged, got %q", got)
	}
}

func TestRun_EmptyInput(t *testing.T) {
	result, err := Run(nil, Options{Fields: []string{"x"}, Char: "*"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
