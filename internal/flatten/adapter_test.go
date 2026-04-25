package flatten

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
)

func makeFilterEntry(ts time.Time, fields map[string]string) filter.Entry {
	return filter.Entry{Timestamp: ts, Fields: fields}
}

func TestFromFilterEntries_Converts(t *testing.T) {
	now := time.Unix(1000, 0)
	input := []filter.Entry{
		makeFilterEntry(now, map[string]string{"level": "info", "msg": "ok"}),
	}
	out := FromFilterEntries(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", out[0].Fields["level"])
	}
	if out[0].Timestamp != now.UnixNano() {
		t.Errorf("unexpected timestamp %d", out[0].Timestamp)
	}
}

func TestToFilterEntries_StringifiesValues(t *testing.T) {
	input := []Entry{
		{Timestamp: 0, Fields: map[string]interface{}{"code": 404, "path": "/api"}},
	}
	out := ToFilterEntries(input)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["code"] != "404" {
		t.Errorf("expected code=404, got %v", out[0].Fields["code"])
	}
	if out[0].Fields["path"] != "/api" {
		t.Errorf("expected path=/api, got %v", out[0].Fields["path"])
	}
}

func TestParseConfig_NilWhenEmpty(t *testing.T) {
	if ParseConfig("", 0) != nil {
		t.Error("expected nil config when no options provided")
	}
}

func TestParseConfig_CustomSeparator(t *testing.T) {
	opts := ParseConfig("_", 0)
	if opts == nil {
		t.Fatal("expected non-nil opts")
	}
	if opts.Separator != "_" {
		t.Errorf("expected separator=_, got %s", opts.Separator)
	}
}

func TestRunFromConfig_NilOptsPassthrough(t *testing.T) {
	input := []filter.Entry{
		makeFilterEntry(time.Now(), map[string]string{"a": "b"}),
	}
	out := RunFromConfig(input, nil)
	if len(out) != len(input) {
		t.Errorf("expected passthrough, got %d entries", len(out))
	}
}

func TestRunFromConfig_FlattensWhenConfigSet(t *testing.T) {
	now := time.Unix(500, 0)
	input := []filter.Entry{
		makeFilterEntry(now, map[string]string{"level": "debug"}),
	}
	opts := ParseConfig(".", 0)
	out := RunFromConfig(input, opts)
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}
