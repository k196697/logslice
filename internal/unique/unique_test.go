package unique

import (
	"testing"
)

func makeEntry(fields map[string]string) Entry {
	return Entry{Fields: fields}
}

func TestRun_ExtractsUniqueValues(t *testing.T) {
	entries := []Entry{
		makeEntry(map[string]string{"level": "info"}),
		makeEntry(map[string]string{"level": "error"}),
		makeEntry(map[string]string{"level": "info"}),
	}
	opts := Options{Fields: []string{"level"}}
	res, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["level"]) != 2 {
		t.Errorf("expected 2 unique values, got %d", len(res["level"]))
	}
}

func TestRun_EmptyFields(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"level": "info"})}
	res, err := Run(entries, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Errorf("expected empty result, got %v", res)
	}
}

func TestRun_LimitCapsValues(t *testing.T) {
	entries := []Entry{
		makeEntry(map[string]string{"svc": "a"}),
		makeEntry(map[string]string{"svc": "b"}),
		makeEntry(map[string]string{"svc": "c"}),
	}
	opts := Options{Fields: []string{"svc"}, Limit: 2}
	res, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["svc"]) != 2 {
		t.Errorf("expected 2 values due to limit, got %d", len(res["svc"]))
	}
}

func TestRun_MissingFieldSkipped(t *testing.T) {
	entries := []Entry{
		makeEntry(map[string]string{"level": "info"}),
		makeEntry(map[string]string{"other": "x"}),
	}
	opts := Options{Fields: []string{"level"}}
	res, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["level"]) != 1 {
		t.Errorf("expected 1 unique value, got %d", len(res["level"]))
	}
}

func TestRun_MultipleFields(t *testing.T) {
	entries := []Entry{
		makeEntry(map[string]string{"level": "info", "svc": "api"}),
		makeEntry(map[string]string{"level": "warn", "svc": "api"}),
		makeEntry(map[string]string{"level": "info", "svc": "db"}),
	}
	opts := Options{Fields: []string{"level", "svc"}}
	res, err := Run(entries, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["level"]) != 2 {
		t.Errorf("expected 2 unique levels, got %d", len(res["level"]))
	}
	if len(res["svc"]) != 2 {
		t.Errorf("expected 2 unique services, got %d", len(res["svc"]))
	}
}

func TestRun_EmptyInput(t *testing.T) {
	res, err := Run([]Entry{}, Options{Fields: []string{"level"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res["level"]) != 0 {
		t.Errorf("expected 0 values for empty input")
	}
}
