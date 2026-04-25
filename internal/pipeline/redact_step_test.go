package pipeline

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeRedactEntry(fields map[string]string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

func TestApplyRedact_NoConfig(t *testing.T) {
	entries := []filter.Entry{
		makeRedactEntry(map[string]string{"email": "user@example.com", "msg": "hello"}),
	}

	result, err := applyRedact(entries, nil, nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Fields["email"] != "user@example.com" {
		t.Errorf("expected email unchanged, got %q", result[0].Fields["email"])
	}
}

func TestApplyRedact_RedactsField(t *testing.T) {
	entries := []filter.Entry{
		makeRedactEntry(map[string]string{"email": "user@example.com", "msg": "hello"}),
	}

	result, err := applyRedact(entries, []string{"email"}, nil, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Fields["email"] != "***" {
		t.Errorf("expected email redacted to ***, got %q", result[0].Fields["email"])
	}
	if result[0].Fields["msg"] != "hello" {
		t.Errorf("expected msg unchanged, got %q", result[0].Fields["msg"])
	}
}

func TestApplyRedact_RedactsWithPattern(t *testing.T) {
	entries := []filter.Entry{
		makeRedactEntry(map[string]string{"msg": "token=abc123 and other stuff"}),
	}

	result, err := applyRedact(entries, nil, []string{"msg=token=\\S+"}, "[REDACTED]")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Fields["msg"] == "token=abc123 and other stuff" {
		t.Errorf("expected msg to be redacted, got original value")
	}
}

func TestApplyRedact_EmptyEntries(t *testing.T) {
	result, err := applyRedact([]filter.Entry{}, []string{"email"}, nil, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
