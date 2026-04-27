package typecast

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeFilterEntry(ts time.Time, fields map[string]string) filter.Entry {
	return filter.Entry{Timestamp: ts, Fields: fields}
}

func TestFromFilterEntries_Converts(t *testing.T) {
	now := time.Unix(1700000000, 0)
	in := []filter.Entry{
		makeFilterEntry(now, map[string]string{"status": "200", "msg": "ok"}),
	}
	out := FromFilterEntries(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["status"] != "200" {
		t.Errorf("expected status=200, got %v", out[0].Fields["status"])
	}
}

func TestToFilterEntries_StringifiesValues(t *testing.T) {
	in := []Entry{
		{Timestamp: 0, Fields: map[string]interface{}{"status": int64(404), "latency": 1.5}},
	}
	out := ToFilterEntries(in)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["status"] != "404" {
		t.Errorf("expected status=404 as string, got %q", out[0].Fields["status"])
	}
	if out[0].Fields["latency"] != "1.5" {
		t.Errorf("expected latency=1.5 as string, got %q", out[0].Fields["latency"])
	}
}

func TestParseConfig_NilWhenEmpty(t *testing.T) {
	opts, err := ParseConfig(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts != nil {
		t.Error("expected nil opts for empty exprs")
	}
}

func TestParseConfig_ValidExprs(t *testing.T) {
	opts, err := ParseConfig([]string{"status:int", "latency:float"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts == nil {
		t.Fatal("expected non-nil opts")
	}
	if len(opts.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(opts.Rules))
	}
}

func TestParseConfig_InvalidExpr(t *testing.T) {
	_, err := ParseConfig([]string{"status:unknown"})
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestRunFromConfig_NilOptsPassthrough(t *testing.T) {
	now := time.Unix(1700000000, 0)
	in := []filter.Entry{makeFilterEntry(now, map[string]string{"x": "1"})}
	out := RunFromConfig(in, nil)
	if len(out) != 1 || out[0].Fields["x"] != "1" {
		t.Error("expected passthrough when opts is nil")
	}
}

func TestRunFromConfig_AppliesRules(t *testing.T) {
	now := time.Unix(1700000000, 0)
	in := []filter.Entry{makeFilterEntry(now, map[string]string{"count": "42"})}
	opts, _ := ParseConfig([]string{"count:int"})
	out := RunFromConfig(in, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["count"] != "42" {
		t.Errorf("expected count=42 (re-stringified), got %q", out[0].Fields["count"])
	}
}
