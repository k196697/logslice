package transform

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeEntry(fields map[string]string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Fields:    fields,
	}
}

func TestParseRule_Rename(t *testing.T) {
	r, err := ParseRule("rename:level=severity")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Op != OpRename || r.Field != "level" || r.Value != "severity" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Drop(t *testing.T) {
	r, err := ParseRule("drop:debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Op != OpDrop || r.Field != "debug" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Set(t *testing.T) {
	r, err := ParseRule("set:env=production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Op != OpSet || r.Field != "env" || r.Value != "production" {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_Invalid(t *testing.T) {
	cases := []string{"nodot", "unknown:x", "drop:", "rename:noequals", "set:=val"}
	for _, c := range cases {
		_, err := ParseRule(c)
		if err == nil {
			t.Errorf("expected error for %q, got nil", c)
		}
	}
}

func TestApply_Rename(t *testing.T) {
	e := makeEntry(map[string]string{"level": "info", "msg": "hello"})
	out := Apply(e, []Rule{{Op: OpRename, Field: "level", Value: "severity"}})
	if _, ok := out.Fields["level"]; ok {
		t.Error("old field 'level' should be removed")
	}
	if out.Fields["severity"] != "info" {
		t.Errorf("expected severity=info, got %q", out.Fields["severity"])
	}
}

func TestApply_Drop(t *testing.T) {
	e := makeEntry(map[string]string{"level": "debug", "msg": "trace"})
	out := Apply(e, []Rule{{Op: OpDrop, Field: "level"}})
	if _, ok := out.Fields["level"]; ok {
		t.Error("field 'level' should have been dropped")
	}
	if out.Fields["msg"] != "trace" {
		t.Error("msg field should be preserved")
	}
}

func TestApply_Set(t *testing.T) {
	e := makeEntry(map[string]string{"msg": "hello"})
	out := Apply(e, []Rule{{Op: OpSet, Field: "env", Value: "staging"}})
	if out.Fields["env"] != "staging" {
		t.Errorf("expected env=staging, got %q", out.Fields["env"])
	}
}

func TestApply_PreservesOriginal(t *testing.T) {
	e := makeEntry(map[string]string{"level": "info"})
	Apply(e, []Rule{{Op: OpDrop, Field: "level"}})
	if e.Fields["level"] != "info" {
		t.Error("original entry should not be mutated")
	}
}

func TestRun_EmptyRules(t *testing.T) {
	entries := []filter.Entry{makeEntry(map[string]string{"a": "1"})}
	out := Run(entries, nil)
	if &out[0] != &entries[0] {
		// same slice returned
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}
