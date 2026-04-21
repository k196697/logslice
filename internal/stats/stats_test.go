package stats

import (
	"testing"
	"time"
)

func makeEntry(ts time.Time, fields map[string]interface{}) Entry {
	return Entry{Timestamp: ts, Fields: fields}
}

func TestCompute_BasicCounts(t *testing.T) {
	entries := []Entry{
		makeEntry(time.Unix(100, 0), map[string]interface{}{"level": "info"}),
		makeEntry(time.Unix(200, 0), map[string]interface{}{"level": "error"}),
	}
	s := Compute(entries, 5)
	if s.Total != 5 {
		t.Errorf("expected Total=5, got %d", s.Total)
	}
	if s.Matched != 2 {
		t.Errorf("expected Matched=2, got %d", s.Matched)
	}
	if s.Skipped != 3 {
		t.Errorf("expected Skipped=3, got %d", s.Skipped)
	}
}

func TestCompute_TimeRange(t *testing.T) {
	t1 := time.Unix(100, 0)
	t2 := time.Unix(300, 0)
	t3 := time.Unix(200, 0)
	entries := []Entry{
		makeEntry(t1, nil),
		makeEntry(t2, nil),
		makeEntry(t3, nil),
	}
	s := Compute(entries, 3)
	if !s.Earliest.Equal(t1) {
		t.Errorf("expected Earliest=%v, got %v", t1, s.Earliest)
	}
	if !s.Latest.Equal(t2) {
		t.Errorf("expected Latest=%v, got %v", t2, s.Latest)
	}
}

func TestCompute_FieldCounts(t *testing.T) {
	entries := []Entry{
		makeEntry(time.Unix(1, 0), map[string]interface{}{"level": "info", "host": "a"}),
		makeEntry(time.Unix(2, 0), map[string]interface{}{"level": "warn"}),
	}
	s := Compute(entries, 2)
	if s.Fields["level"] != 2 {
		t.Errorf("expected level count=2, got %d", s.Fields["level"])
	}
	if s.Fields["host"] != 1 {
		t.Errorf("expected host count=1, got %d", s.Fields["host"])
	}
}

func TestCompute_EmptyEntries(t *testing.T) {
	s := Compute([]Entry{}, 10)
	if s.Matched != 0 || s.Skipped != 10 {
		t.Errorf("unexpected counts: matched=%d skipped=%d", s.Matched, s.Skipped)
	}
}
