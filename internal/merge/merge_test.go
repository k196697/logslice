package merge

import (
	"testing"
	"time"
)

func makeEntry(ts time.Time, raw string) Entry {
	return Entry{
		Timestamp: ts,
		Fields:    map[string]string{"msg": raw},
		Raw:       raw,
	}
}

var (
	t1 = time.Date(2024, 1, 1, 0, 0, 1, 0, time.UTC)
	t2 = time.Date(2024, 1, 1, 0, 0, 2, 0, time.UTC)
	t3 = time.Date(2024, 1, 1, 0, 0, 3, 0, time.UTC)
)

func TestRun_MergesAndSortsByTime(t *testing.T) {
	src1 := []Entry{makeEntry(t3, "c"), makeEntry(t1, "a")}
	src2 := []Entry{makeEntry(t2, "b")}

	result := Run([][]Entry{src1, src2}, DefaultOptions())

	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	expected := []string{"a", "b", "c"}
	for i, e := range result {
		if e.Raw != expected[i] {
			t.Errorf("pos %d: want %q, got %q", i, expected[i], e.Raw)
		}
	}
}

func TestRun_NoSort(t *testing.T) {
	src := []Entry{makeEntry(t3, "c"), makeEntry(t1, "a")}
	opts := Options{SortByTime: false}
	result := Run([][]Entry{src}, opts)
	if result[0].Raw != "c" {
		t.Errorf("expected original order preserved, got %q", result[0].Raw)
	}
}

func TestRun_ZeroTimestampAtEnd(t *testing.T) {
	zero := time.Time{}
	src := []Entry{makeEntry(zero, "no-ts"), makeEntry(t1, "first")}
	result := Run([][]Entry{src}, DefaultOptions())
	if result[0].Raw != "first" {
		t.Errorf("expected timestamped entry first, got %q", result[0].Raw)
	}
	if result[1].Raw != "no-ts" {
		t.Errorf("expected zero-ts entry last, got %q", result[1].Raw)
	}
}

func TestRun_DeduplicateConsecutive(t *testing.T) {
	src := []Entry{
		makeEntry(t1, "dup"),
		makeEntry(t2, "dup"),
		makeEntry(t3, "unique"),
	}
	opts := Options{SortByTime: false, DeduplicateConsecutive: true}
	result := Run([][]Entry{src}, opts)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries after dedup, got %d", len(result))
	}
	if result[0].Raw != "dup" || result[1].Raw != "unique" {
		t.Errorf("unexpected dedup result: %v", result)
	}
}

func TestRun_EmptyInput(t *testing.T) {
	result := Run([][]Entry{}, DefaultOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
