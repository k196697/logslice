package rename

import (
	"testing"
)

func makeEntry(ts int64, fields map[string]string) Entry {
	return Entry{Timestamp: ts, Fields: fields}
}

func TestParseRule_Valid(t *testing.T) {
	r, err := ParseRule("old:new")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.From != "old" || r.To != "new" {
		t.Errorf("expected old->new, got %q->%q", r.From, r.To)
	}
}

func TestParseRule_MissingColon(t *testing.T) {
	_, err := ParseRule("oldnew")
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRule_EmptyFrom(t *testing.T) {
	_, err := ParseRule(":new")
	if err == nil {
		t.Fatal("expected error for empty from field")
	}
}

func TestParseRule_EmptyTo(t *testing.T) {
	_, err := ParseRule("old:")
	if err == nil {
		t.Fatal("expected error for empty to field")
	}
}

func TestRun_RenamesField(t *testing.T) {
	entries := []Entry{
		makeEntry(1, map[string]string{"msg": "hello", "lvl": "info"}),
	}
	opts := Options{Rules: []Rule{{From: "lvl", To: "level"}}}
	out := Run(entries, opts)
	if _, ok := out[0].Fields["lvl"]; ok {
		t.Error("old field 'lvl' should be removed")
	}
	if out[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %q", out[0].Fields["level"])
	}
}

func TestRun_SkipsMissingField(t *testing.T) {
	entries := []Entry{
		makeEntry(1, map[string]string{"msg": "hello"}),
	}
	opts := Options{Rules: []Rule{{From: "missing", To: "renamed"}}}
	out := Run(entries, opts)
	if _, ok := out[0].Fields["renamed"]; ok {
		t.Error("renamed field should not appear when source is missing")
	}
	if out[0].Fields["msg"] != "hello" {
		t.Error("unrelated field should be unchanged")
	}
}

func TestRun_PreservesTimestamp(t *testing.T) {
	entries := []Entry{
		makeEntry(9999, map[string]string{"a": "1"}),
	}
	opts := Options{Rules: []Rule{{From: "a", To: "b"}}}
	out := Run(entries, opts)
	if out[0].Timestamp != 9999 {
		t.Errorf("expected timestamp 9999, got %d", out[0].Timestamp)
	}
}

func TestRun_EmptyInput(t *testing.T) {
	out := Run(nil, DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}

func TestParseRules_MultipleValid(t *testing.T) {
	rules, err := ParseRules([]string{"a:b", "c:d"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_InvalidStopsEarly(t *testing.T) {
	_, err := ParseRules([]string{"a:b", "bad"})
	if err == nil {
		t.Fatal("expected error for invalid rule")
	}
}
