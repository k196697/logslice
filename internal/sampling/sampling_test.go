package sampling

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeEntry(msg string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    map[string]string{"msg": msg},
		Raw:       msg,
	}
}

func makeEntries(n int) []filter.Entry {
	entries := make([]filter.Entry, n)
	for i := 0; i < n; i++ {
		entries[i] = makeEntry("entry")
	}
	return entries
}

func TestRun_KeepsAllWhenRateOne(t *testing.T) {
	entries := makeEntries(20)
	opts := DefaultOptions()
	result := Run(entries, opts)
	if len(result) != 20 {
		t.Errorf("expected 20 entries, got %d", len(result))
	}
}

func TestRun_KeepsNoneWhenRateZero(t *testing.T) {
	entries := makeEntries(20)
	opts := Options{Rate: 0.0}
	result := Run(entries, opts)
	if len(result) != 0 {
		t.Errorf("expected 0 entries, got %d", len(result))
	}
}

func TestRun_EveryN(t *testing.T) {
	entries := makeEntries(10)
	opts := Options{Every: 2, Rate: 1.0}
	result := Run(entries, opts)
	if len(result) != 5 {
		t.Errorf("expected 5 entries, got %d", len(result))
	}
}

func TestRun_EveryNPrioritisedOverRate(t *testing.T) {
	entries := makeEntries(9)
	opts := Options{Every: 3, Rate: 0.1}
	result := Run(entries, opts)
	if len(result) != 3 {
		t.Errorf("expected 3 entries (every-N wins), got %d", len(result))
	}
}

func TestRun_EmptyInput(t *testing.T) {
	result := Run([]filter.Entry{}, DefaultOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d", len(result))
	}
}

func TestRun_ReproducibleWithSeed(t *testing.T) {
	entries := makeEntries(100)
	opts := Options{Rate: 0.5, Seed: 42}
	first := Run(entries, opts)
	second := Run(entries, opts)
	if len(first) != len(second) {
		t.Errorf("expected same count with same seed: %d vs %d", len(first), len(second))
	}
}
