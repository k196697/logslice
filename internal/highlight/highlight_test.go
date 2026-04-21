package highlight

import (
	"strings"
	"testing"
)

func TestColorizeField_MatchingRule(t *testing.T) {
	h := New(true, []Rule{
		{Field: "level", Value: "error", Color: Red},
	})
	result := h.ColorizeField("level", "error")
	if !strings.Contains(result, Red) {
		t.Errorf("expected Red color code in result, got: %q", result)
	}
	if !strings.Contains(result, "error") {
		t.Errorf("expected value 'error' in result, got: %q", result)
	}
	if !strings.Contains(result, Reset) {
		t.Errorf("expected Reset code in result, got: %q", result)
	}
}

func TestColorizeField_NoMatchingRule(t *testing.T) {
	h := New(true, []Rule{
		{Field: "level", Value: "error", Color: Red},
	})
	result := h.ColorizeField("level", "info")
	if result != "info" {
		t.Errorf("expected plain 'info', got: %q", result)
	}
}

func TestColorizeField_WildcardValue(t *testing.T) {
	h := New(true, []Rule{
		{Field: "host", Value: "*", Color: Cyan},
	})
	result := h.ColorizeField("host", "myhost")
	if !strings.Contains(result, Cyan) {
		t.Errorf("expected Cyan in wildcard match, got: %q", result)
	}
}

func TestColorizeField_DisabledHighlighter(t *testing.T) {
	h := New(false, []Rule{
		{Field: "level", Value: "error", Color: Red},
	})
	result := h.ColorizeField("level", "error")
	if result != "error" {
		t.Errorf("expected plain text when disabled, got: %q", result)
	}
}

func TestColorizeLine_ContainsKeyword(t *testing.T) {
	h := New(true, nil)
	result := h.ColorizeLine("this is an error line", "error", Red)
	if !strings.Contains(result, Red) {
		t.Errorf("expected Red in line highlight, got: %q", result)
	}
}

func TestColorizeLine_NoKeyword(t *testing.T) {
	h := New(true, nil)
	result := h.ColorizeLine("some normal line", "error", Red)
	if result != "some normal line" {
		t.Errorf("expected unchanged line, got: %q", result)
	}
}

func TestParseColor_KnownColors(t *testing.T) {
	cases := []struct {
		name     string
		expected string
	}{
		{"red", Red},
		{"green", Green},
		{"yellow", Yellow},
		{"blue", Blue},
		{"cyan", Cyan},
		{"bold", Bold},
	}
	for _, tc := range cases {
		got, ok := ParseColor(tc.name)
		if !ok {
			t.Errorf("ParseColor(%q): expected ok=true", tc.name)
		}
		if got != tc.expected {
			t.Errorf("ParseColor(%q): expected %q, got %q", tc.name, tc.expected, got)
		}
	}
}

func TestParseColor_Unknown(t *testing.T) {
	_, ok := ParseColor("purple")
	if ok {
		t.Error("expected ok=false for unknown color")
	}
}
