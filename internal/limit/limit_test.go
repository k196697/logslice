package limit

import (
	"testing"
)

func makeEntry(nano int64, fields map[string]string) Entry {
	return Entry{Fields: fields, TimestampNano: nano}
}

func makeEntries(n int) []Entry {
	entries := make([]Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = makeEntry(int64(i+1)*1000, map[string]string{"idx": string(rune('0'+i))})
	}
	return entries
}

func TestRun_NoLimit(t *testing.T) {
	entries := makeEntries(5)
	result := Run(entries, DefaultOptions())
	if len(result) != 5 {
		t.Fatalf("expected 5 entries, got %d", len(result))
	}
}

func TestRun_LimitFromStart(t *testing.T) {
	entries := makeEntries(10)
	opts := Options{MaxEntries: 3, FromEnd: false}
	result := Run(entries, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].TimestampNano != 1000 {
		t.Errorf("expected first entry nano=1000, got %d", result[0].TimestampNano)
	}
}

func TestRun_LimitFromEnd(t *testing.T) {
	entries := makeEntries(10)
	opts := Options{MaxEntries: 3, FromEnd: true}
	result := Run(entries, opts)
	if len(result) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(result))
	}
	if result[0].TimestampNano != 8000 {
		t.Errorf("expected first entry nano=8000, got %d", result[0].TimestampNano)
	}
}

func TestRun_LimitLargerThanInput(t *testing.T) {
	entries := makeEntries(4)
	opts := Options{MaxEntries: 100}
	result := Run(entries, opts)
	if len(result) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(result))
	}
}

func TestRun_EmptyInput(t *testing.T) {
	result := Run([]Entry{}, Options{MaxEntries: 5})
	if len(result) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(result))
	}
}

func TestParseConfig_NilWhenZero(t *testing.T) {
	if ParseConfig(0, false) != nil {
		t.Error("expected nil for maxEntries=0")
	}
}

func TestParseConfig_ValidLimit(t *testing.T) {
	opts := ParseConfig(10, true)
	if opts == nil {
		t.Fatal("expected non-nil options")
	}
	if opts.MaxEntries != 10 {
		t.Errorf("expected MaxEntries=10, got %d", opts.MaxEntries)
	}
	if !opts.FromEnd {
		t.Error("expected FromEnd=true")
	}
}
