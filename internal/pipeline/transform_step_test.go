package pipeline

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeTransformEntry(fields map[string]string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
	}
}

func TestApplyTransform_NoRules(t *testing.T) {
	entries := []filter.Entry{
		makeTransformEntry(map[string]string{"level": "info", "msg": "hello"}),
	}
	cfg := &Config{TransformRules: nil}

	result, err := applyTransform(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
}

func TestApplyTransform_RenameField(t *testing.T) {
	entries := []filter.Entry{
		makeTransformEntry(map[string]string{"level": "info", "msg": "hello"}),
	}
	cfg := &Config{TransformRules: []string{"rename:msg=message"}}

	result, err := applyTransform(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if _, ok := result[0].Fields["message"]; !ok {
		t.Error("expected field 'message' after rename, not found")
	}
	if _, ok := result[0].Fields["msg"]; ok {
		t.Error("expected field 'msg' to be removed after rename")
	}
}

func TestApplyTransform_DropField(t *testing.T) {
	entries := []filter.Entry{
		makeTransformEntry(map[string]string{"level": "debug", "secret": "abc123"}),
	}
	cfg := &Config{TransformRules: []string{"drop:secret"}}

	result, err := applyTransform(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result[0].Fields["secret"]; ok {
		t.Error("expected field 'secret' to be dropped")
	}
}

func TestApplyTransform_InvalidRule(t *testing.T) {
	entries := []filter.Entry{
		makeTransformEntry(map[string]string{"level": "info"}),
	}
	cfg := &Config{TransformRules: []string{"bogus:rule:format"}}

	_, err := applyTransform(entries, cfg)
	if err == nil {
		t.Fatal("expected error for invalid transform rule, got nil")
	}
}

func TestApplyTransform_SetField(t *testing.T) {
	entries := []filter.Entry{
		makeTransformEntry(map[string]string{"level": "info"}),
	}
	cfg := &Config{TransformRules: []string{"set:env=production"}}

	result, err := applyTransform(entries, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := result[0].Fields["env"]; !ok || v != "production" {
		t.Errorf("expected field 'env'='production', got %q", v)
	}
}
