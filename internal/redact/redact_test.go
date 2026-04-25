package redact

import (
	"testing"
)

func makeEntry(fields map[string]string) Entry {
	return Entry{Fields: fields}
}

func TestRun_RedactsEntireField(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"password": "secret", "user": "alice"})}
	rules := []Rule{{Field: "password"}}
	out := Run(entries, Options{Rules: rules})
	if out[0].Fields["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", out[0].Fields["password"])
	}
	if out[0].Fields["user"] != "alice" {
		t.Error("user field should be unchanged")
	}
}

func TestRun_RedactsWithPattern(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"msg": "token=abc123 other"})}
	rule, _ := ParseRule(`msg:token=\w+`)
	out := Run(entries, Options{Rules: []Rule{rule}})
	if out[0].Fields["msg"] != "[REDACTED] other" {
		t.Errorf("unexpected value: %q", out[0].Fields["msg"])
	}
}

func TestRun_CustomMask(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"email": "user@example.com"})}
	rule, _ := ParseRule(`email:[^@]+@[^@]+=[EMAIL]`)
	out := Run(entries, Options{Rules: []Rule{rule}})
	if out[0].Fields["email"] != "[EMAIL]" {
		t.Errorf("unexpected value: %q", out[0].Fields["email"])
	}
}

func TestRun_SkipsMissingField(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"user": "bob"})}
	rules := []Rule{{Field: "password"}}
	out := Run(entries, Options{Rules: rules})
	if _, ok := out[0].Fields["password"]; ok {
		t.Error("password field should not be introduced")
	}
}

func TestRun_EmptyRules(t *testing.T) {
	entries := []Entry{makeEntry(map[string]string{"secret": "value"})}
	out := Run(entries, Options{})
	if out[0].Fields["secret"] != "value" {
		t.Error("no rules should leave entries unchanged")
	}
}

func TestRun_OriginalUnmodified(t *testing.T) {
	orig := []Entry{makeEntry(map[string]string{"token": "abc"})}
	Run(orig, Options{Rules: []Rule{{Field: "token"}}})
	if orig[0].Fields["token"] != "abc" {
		t.Error("original entries must not be mutated")
	}
}

func TestParseRule_FieldOnly(t *testing.T) {
	r, err := ParseRule("password")
	if err != nil || r.Field != "password" || r.Pattern != nil {
		t.Errorf("unexpected rule: %+v err: %v", r, err)
	}
}

func TestParseRule_WithPattern(t *testing.T) {
	r, err := ParseRule(`email:[^@]+@\S+`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Field != "email" || r.Pattern == nil {
		t.Errorf("unexpected rule: %+v", r)
	}
}

func TestParseRule_InvalidRegex(t *testing.T) {
	_, err := ParseRule("field:[invalid")
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}
