package unique

import (
	"testing"

	"github.com/user/logslice/internal/filter"
)

func makeFilterEntry(fields map[string]string) filter.Entry {
	return filter.Entry{Fields: fields}
}

func TestFromFilterEntries_Converts(t *testing.T) {
	input := []filter.Entry{
		makeFilterEntry(map[string]string{"level": "info", "svc": "api"}),
		makeFilterEntry(map[string]string{"level": "error"}),
	}
	out := FromFilterEntries(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Fields["level"] != "info" {
		t.Errorf("expected level=info, got %s", out[0].Fields["level"])
	}
	if out[1].Fields["level"] != "error" {
		t.Errorf("expected level=error, got %s", out[1].Fields["level"])
	}
}

func TestParseConfig_NilWhenEmpty(t *testing.T) {
	if ParseConfig("", 0) != nil {
		t.Error("expected nil for empty fields string")
	}
}

func TestParseConfig_ValidFields(t *testing.T) {
	opts := ParseConfig("level,svc", 0)
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
	if len(opts.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(opts.Fields))
	}
	if opts.Fields[0] != "level" || opts.Fields[1] != "svc" {
		t.Errorf("unexpected fields: %v", opts.Fields)
	}
}

func TestParseConfig_SetsLimit(t *testing.T) {
	opts := ParseConfig("level", 5)
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
	if opts.Limit != 5 {
		t.Errorf("expected limit=5, got %d", opts.Limit)
	}
}

func TestParseConfig_TrimsWhitespace(t *testing.T) {
	opts := ParseConfig(" level , svc ", 0)
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
	if opts.Fields[0] != "level" || opts.Fields[1] != "svc" {
		t.Errorf("expected trimmed fields, got %v", opts.Fields)
	}
}

func TestParseConfig_NilWhenOnlyCommas(t *testing.T) {
	if ParseConfig(" , , ", 0) != nil {
		t.Error("expected nil when all parts are empty after trim")
	}
}
