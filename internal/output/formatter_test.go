package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(ts time.Time, fields map[string]interface{}) parser.LogEntry {
	return parser.LogEntry{Timestamp: ts, Fields: fields}
}

var sampleTime = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatJSON, "timestamp")
	entry := makeEntry(sampleTime, map[string]interface{}{"level": "info", "msg": "hello"})

	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"level\":\"info\"") {
		t.Errorf("expected level field in output, got: %s", out)
	}
	if !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected valid JSON object, got: %s", out)
	}
}

func TestFormatter_Text(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatText, "timestamp")
	entry := makeEntry(sampleTime, map[string]interface{}{"level": "warn", "timestamp": "ignored"})

	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "2024-01-15T10:30:00Z") {
		t.Errorf("expected formatted timestamp in output, got: %s", out)
	}
	if strings.Contains(out, "timestamp=ignored") {
		t.Errorf("time field should be excluded from text output, got: %s", out)
	}
}

func TestFormatter_Compact(t *testing.T) {
	var buf bytes.Buffer
	f := output.NewFormatter(&buf, output.FormatCompact, "")
	entry := makeEntry(sampleTime, map[string]interface{}{"level": "error", "msg": "something failed"})

	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected uppercased level, got: %s", out)
	}
	if !strings.Contains(out, "something failed") {
		t.Errorf("expected message in output, got: %s", out)
	}
}

func TestFormatter_DefaultTimeField(t *testing.T) {
	f := output.NewFormatter(&bytes.Buffer{}, output.FormatText, "")
	if f.TimeField != "timestamp" {
		t.Errorf("expected default time field 'timestamp', got %q", f.TimeField)
	}
}
