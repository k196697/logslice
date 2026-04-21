package parser_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/parser"
)

func TestDetectFormatFromString_JSON(t *testing.T) {
	line := `{"level":"info","msg":"started","ts":"2024-01-01T00:00:00Z"}`
	got := parser.DetectFormatFromString(line)
	if got != "json" {
		t.Errorf("expected json, got %q", got)
	}
}

func TestDetectFormatFromString_CSV(t *testing.T) {
	line := "2024-01-01T00:00:00Z,info,started"
	got := parser.DetectFormatFromString(line)
	if got != "csv" {
		t.Errorf("expected csv, got %q", got)
	}
}

func TestDetectFormatFromString_Logfmt(t *testing.T) {
	line := `ts=2024-01-01T00:00:00Z level=info msg=started`
	got := parser.DetectFormatFromString(line)
	if got != "logfmt" {
		t.Errorf("expected logfmt, got %q", got)
	}
}

func TestDetectFormatFromString_Empty(t *testing.T) {
	got := parser.DetectFormatFromString("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestDetectFormatFromString_Unknown(t *testing.T) {
	got := parser.DetectFormatFromString("just a plain log line with no structure")
	if got != "" {
		t.Errorf("expected empty string for unknown format, got %q", got)
	}
}

func TestDetectFormat_ReadsFirstNonEmptyLine(t *testing.T) {
	input := "\n\n" + `{"level":"info","msg":"hello"}` + "\n"
	r := strings.NewReader(input)
	got, err := parser.DetectFormat(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "json" {
		t.Errorf("expected json, got %q", got)
	}
}

func TestDetectFormat_EmptyReader(t *testing.T) {
	r := strings.NewReader("")
	got, err := parser.DetectFormat(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("expected empty string for empty reader, got %q", got)
	}
}
