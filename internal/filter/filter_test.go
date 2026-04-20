package filter_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeEntry(ts time.Time, fields map[string]interface{}) filter.LogEntry {
	return filter.LogEntry{Timestamp: ts, Fields: fields, Raw: "{}"}
}

var (
	t1 = time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	t3 = time.Date(2024, 1, 1, 14, 0, 0, 0, time.UTC)
)

func TestApply_TimeRange(t *testing.T) {
	entry := makeEntry(t2, nil)
	from := t1
	to := t3
	opts := filter.Options{From: &from, To: &to}
	if !filter.Apply(entry, opts) {
		t.Error("expected entry within range to pass")
	}
}

func TestApply_BeforeFrom(t *testing.T) {
	entry := makeEntry(t1, nil)
	from := t2
	opts := filter.Options{From: &from}
	if filter.Apply(entry, opts) {
		t.Error("expected entry before From to be filtered out")
	}
}

func TestApply_AfterTo(t *testing.T) {
	entry := makeEntry(t3, nil)
	to := t2
	opts := filter.Options{To: &to}
	if filter.Apply(entry, opts) {
		t.Error("expected entry after To to be filtered out")
	}
}

func TestApply_FieldKeyExists(t *testing.T) {
	entry := makeEntry(t1, map[string]interface{}{"level": "error"})
	opts := filter.Options{FieldKey: "level"}
	if !filter.Apply(entry, opts) {
		t.Error("expected entry with field key to pass")
	}
}

func TestApply_FieldKeyMissing(t *testing.T) {
	entry := makeEntry(t1, map[string]interface{}{})
	opts := filter.Options{FieldKey: "level"}
	if filter.Apply(entry, opts) {
		t.Error("expected entry missing field key to be filtered out")
	}
}

func TestApply_FieldKeyValue(t *testing.T) {
	entry := makeEntry(t1, map[string]interface{}{"level": "error"})
	opts := filter.Options{FieldKey: "level", FieldVal: "error"}
	if !filter.Apply(entry, opts) {
		t.Error("expected entry with matching field value to pass")
	}
}

func TestApply_FieldKeyValueMismatch(t *testing.T) {
	entry := makeEntry(t1, map[string]interface{}{"level": "info"})
	opts := filter.Options{FieldKey: "level", FieldVal: "error"}
	if filter.Apply(entry, opts) {
		t.Error("expected entry with mismatched field value to be filtered out")
	}
}

func TestRun_FiltersMultiple(t *testing.T) {
	entries := []filter.LogEntry{
		makeEntry(t1, map[string]interface{}{"level": "info"}),
		makeEntry(t2, map[string]interface{}{"level": "error"}),
		makeEntry(t3, map[string]interface{}{"level": "error"}),
	}
	opts := filter.Options{FieldKey: "level", FieldVal: "error"}
	result := filter.Run(entries, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}
