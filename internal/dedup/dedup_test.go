package dedup

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/filter"
)

func makeEntry(fields map[string]string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

func TestRun_RemovesGlobalDuplicates(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(map[string]string{"level": "error", "msg": "fail"}),
		makeEntry(map[string]string{"level": "info", "msg": "ok"}),
		makeEntry(map[string]string{"level": "error", "msg": "fail"}),
	}
	result := Run(entries, DefaultOptions())
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestRun_ConsecutiveOnly(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(map[string]string{"msg": "dup"}),
		makeEntry(map[string]string{"msg": "dup"}),
		makeEntry(map[string]string{"msg": "other"}),
		makeEntry(map[string]string{"msg": "dup"}),
	}
	opts := Options{Consecutive: true}
	result := Run(entries, opts)
	// second dup removed, but last dup kept (not consecutive with first)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}

func TestRun_FieldSubset(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(map[string]string{"level": "error", "msg": "a", "svc": "api"}),
		makeEntry(map[string]string{"level": "error", "msg": "b", "svc": "api"}),
		makeEntry(map[string]string{"level": "info", "msg": "c", "svc": "api"}),
	}
	// Dedup only on "level" — first two collapse to one.
	opts := Options{Fields: []string{"level"}}
	result := Run(entries, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
}

func TestRun_EmptyInput(t *testing.T) {
	result := Run(nil, DefaultOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty result, got %d", len(result))
	}
}

func TestRun_NoDuplicates(t *testing.T) {
	entries := []filter.Entry{
		makeEntry(map[string]string{"msg": "a"}),
		makeEntry(map[string]string{"msg": "b"}),
		makeEntry(map[string]string{"msg": "c"}),
	}
	result := Run(entries, DefaultOptions())
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
}
