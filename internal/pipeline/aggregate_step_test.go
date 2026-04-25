package pipeline

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
)

func makeAggEntry(fields map[string]string) filter.Entry {
	f := make(map[string]interface{})
	for k, v := range fields {
		f[k] = v
	}
	return filter.Entry{Fields: f}
}

func TestParseAggregateConfig_NilWhenNoGroupBy(t *testing.T) {
	cfg, err := ParseAggregateConfig("", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg != nil {
		t.Errorf("expected nil config, got %+v", cfg)
	}
}

func TestParseAggregateConfig_ValidGroupBy(t *testing.T) {
	cfg, err := ParseAggregateConfig("level", "count")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.GroupBy != "level" {
		t.Errorf("expected GroupBy=level, got %s", cfg.GroupBy)
	}
}

func TestParseAggregateConfig_InvalidMetric(t *testing.T) {
	_, err := ParseAggregateConfig("level", "unknown_metric")
	if err == nil {
		t.Fatal("expected error for unknown metric")
	}
	if !strings.Contains(err.Error(), "unsupported metric") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestApplyAggregate_GroupsAndCounts(t *testing.T) {
	entries := []filter.Entry{
		makeAggEntry(map[string]string{"level": "info"}),
		makeAggEntry(map[string]string{"level": "info"}),
		makeAggEntry(map[string]string{"level": "error"}),
	}

	cfg, err := ParseAggregateConfig("level", "count")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := applyAggregate(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(result))
	}
}

func TestApplyAggregate_NilConfigPassesThrough(t *testing.T) {
	entries := []filter.Entry{
		makeAggEntry(map[string]string{"level": "info"}),
	}

	result, err := applyAggregate(entries, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 entry, got %d", len(result))
	}
}
