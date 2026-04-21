package highlight

import (
	"testing"
)

func TestParseRule_Valid(t *testing.T) {
	r, err := ParseRule("level=error:red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Field != "level" {
		t.Errorf("expected field 'level', got %q", r.Field)
	}
	if r.Value != "error" {
		t.Errorf("expected value 'error', got %q", r.Value)
	}
	if r.Color != Red {
		t.Errorf("expected Red color, got %q", r.Color)
	}
}

func TestParseRule_Wildcard(t *testing.T) {
	r, err := ParseRule("host=*:cyan")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Value != "*" {
		t.Errorf("expected wildcard value '*', got %q", r.Value)
	}
	if r.Color != Cyan {
		t.Errorf("expected Cyan color, got %q", r.Color)
	}
}

func TestParseRule_MissingColon(t *testing.T) {
	_, err := ParseRule("level=error")
	if err == nil {
		t.Error("expected error for missing color suffix")
	}
}

func TestParseRule_MissingEquals(t *testing.T) {
	_, err := ParseRule("levelerror:red")
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRule_EmptyField(t *testing.T) {
	_, err := ParseRule("=error:red")
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestParseRule_UnknownColor(t *testing.T) {
	_, err := ParseRule("level=error:purple")
	if err == nil {
		t.Error("expected error for unknown color")
	}
}

func TestParseRules_AllValid(t *testing.T) {
	rules, err := ParseRules([]string{"level=error:red", "host=*:cyan"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_PartialError(t *testing.T) {
	rules, err := ParseRules([]string{"level=error:red", "bad-rule"})
	if err == nil {
		t.Error("expected error for partially invalid rules")
	}
	if len(rules) != 1 {
		t.Errorf("expected 1 valid rule, got %d", len(rules))
	}
}

func TestParseRules_Empty(t *testing.T) {
	rules, err := ParseRules([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}
