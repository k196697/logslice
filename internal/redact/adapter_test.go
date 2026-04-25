package redact

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
)

func makeFilterEntry(fields map[string]string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "",
	}
}

func TestFromFilterEntries_Converts(t *testing.T) {
	in := []filter.Entry{makeFilterEntry(map[string]string{"a": "1", "b": "2"})}
	out := FromFilterEntries(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["a"] != "1" || out[0].Fields["b"] != "2" {
		t.Error("field values not preserved")
	}
}

func TestToFilterEntries_PreservesTimestamp(t *testing.T) {
	ts := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	orig := []filter.Entry{{Timestamp: ts, Fields: map[string]string{"x": "y"}, Raw: "raw"}}
	redacted := []Entry{{Fields: map[string]string{"x": "[REDACTED]"}}}
	out := ToFilterEntries(redacted, orig)
	if !out[0].Timestamp.Equal(ts) {
		t.Error("timestamp not preserved")
	}
	if out[0].Fields["x"] != "[REDACTED]" {
		t.Error("redacted value not carried through")
	}
	if out[0].Raw != "raw" {
		t.Error("raw not preserved")
	}
}

func TestParseConfig_NilWhenEmpty(t *testing.T) {
	opts, err := ParseConfig(nil)
	if err != nil || opts != nil {
		t.Errorf("expected nil opts and nil err, got %v %v", opts, err)
	}
}

func TestParseConfig_ValidExprs(t *testing.T) {
	opts, err := ParseConfig([]string{"password", "token"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts == nil || len(opts.Rules) != 2 {
		t.Errorf("expected 2 rules, got %+v", opts)
	}
}

func TestParseConfig_InvalidRegex(t *testing.T) {
	_, err := ParseConfig([]string{"field:[bad"})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestRunFromConfig_NilOpts(t *testing.T) {
	entries := []filter.Entry{makeFilterEntry(map[string]string{"secret": "val"})}
	out := RunFromConfig(entries, nil)
	if out[0].Fields["secret"] != "val" {
		t.Error("nil opts should leave entries unchanged")
	}
}

func TestRunFromConfig_AppliesRedaction(t *testing.T) {
	entries := []filter.Entry{makeFilterEntry(map[string]string{"password": "hunter2"})}
	opts, _ := ParseConfig([]string{"password"})
	out := RunFromConfig(entries, opts)
	if out[0].Fields["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out[0].Fields["password"])
	}
}
