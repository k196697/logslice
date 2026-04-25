package truncate

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeFilterEntry(fields map[string]string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestFromFilterEntries_Converts(t *testing.T) {
	src := []filter.Entry{
		makeFilterEntry(map[string]string{"msg": "hello", "level": "info"}),
	}
	out := FromFilterEntries(src)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["msg"] != "hello" {
		t.Errorf("expected 'hello', got %q", out[0].Fields["msg"])
	}
}

func TestToFilterEntries_PreservesTimestamp(t *testing.T) {
	now := time.Now()
	orig := []filter.Entry{
		{Timestamp: now, Fields: map[string]string{"msg": "hi"}, Raw: "raw"},
	}
	truncated := []Entry{
		{Fields: map[string]string{"msg": "h..."}},
	}
	out := ToFilterEntries(truncated, orig)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if !out[0].Timestamp.Equal(now) {
		t.Errorf("timestamp not preserved")
	}
	if out[0].Fields["msg"] != "h..." {
		t.Errorf("expected truncated value, got %q", out[0].Fields["msg"])
	}
	if out[0].Raw != "raw" {
		t.Errorf("raw not preserved")
	}
}

func TestParseConfig_NilWhenEmpty(t *testing.T) {
	opts, err := ParseConfig("", 0, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts != nil {
		t.Errorf("expected nil opts for empty config")
	}
}

func TestParseConfig_ValidMaxLength(t *testing.T) {
	opts, err := ParseConfig("msg,level", 64, "…")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts == nil {
		t.Fatal("expected non-nil opts")
	}
	if opts.MaxLength != 64 {
		t.Errorf("MaxLength: got %d, want 64", opts.MaxLength)
	}
	if len(opts.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(opts.Fields))
	}
}

func TestRunFromConfig_NilOptsPassthrough(t *testing.T) {
	entries := []filter.Entry{makeFilterEntry(map[string]string{"msg": "unchanged"})}
	out, err := RunFromConfig(entries, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 || out[0].Fields["msg"] != "unchanged" {
		t.Errorf("expected passthrough, got %v", out)
	}
}
