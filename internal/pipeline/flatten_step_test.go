package pipeline

import (
	"testing"
	"time"

	"logslice/internal/filter"
)

func makeFlattenEntry(fields map[string]string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Fields:    fields,
	}
}

func TestParseFlattenConfig_NilWhenDisabled(t *testing.T) {
	opts := ParseFlattenConfig(".", 0, false)
	if opts != nil {
		t.Errorf("expected nil when disabled, got %+v", opts)
	}
}

func TestParseFlattenConfig_DefaultsWhenEnabled(t *testing.T) {
	opts := ParseFlattenConfig("", 0, true)
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
	if opts.Separator == "" {
		t.Error("expected default separator to be set")
	}
}

func TestParseFlattenConfig_CustomSeparator(t *testing.T) {
	opts := ParseFlattenConfig("_", 3, true)
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
	if opts.Separator != "_" {
		t.Errorf("expected separator '_', got %q", opts.Separator)
	}
	if opts.MaxDepth != 3 {
		t.Errorf("expected max depth 3, got %d", opts.MaxDepth)
	}
}

func TestApplyFlatten_NilOptsReturnsUnchanged(t *testing.T) {
	entries := []filter.Entry{
		makeFlattenEntry(map[string]string{"msg": "hello"}),
	}
	out, err := applyFlatten(entries, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != len(entries) {
		t.Errorf("expected %d entries, got %d", len(entries), len(out))
	}
}

func TestApplyFlatten_EmptyEntries(t *testing.T) {
	out, err := applyFlatten([]filter.Entry{}, ParseFlattenConfig(".", 0, true))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected 0 entries, got %d", len(out))
	}
}

func TestApplyFlatten_PassesThroughFlatEntries(t *testing.T) {
	entries := []filter.Entry{
		makeFlattenEntry(map[string]string{"level": "info", "msg": "ok"}),
	}
	opts := ParseFlattenConfig(".", 0, true)
	out, err := applyFlatten(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["level"] != "info" {
		t.Errorf("expected field 'level'='info', got %q", out[0].Fields["level"])
	}
}
