package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
)

func TestRun_JSONLines_Filtered(t *testing.T) {
	lines := []string{
		`{"time":"2024-01-01T10:00:00Z","level":"info","msg":"start"}`,
		`{"time":"2024-01-01T11:00:00Z","level":"error","msg":"fail"}`,
		`{"time":"2024-01-01T12:00:00Z","level":"info","msg":"end"}`,
	}

	var buf bytes.Buffer
	r := pipeline.New(pipeline.Options{
		Format:    "json",
		TimeField: "time",
		From:      "2024-01-01T10:30:00Z",
		OutputFmt: output.FormatJSON,
	}, &buf)

	if err := r.Run(lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "start") {
		t.Error("expected 'start' entry to be filtered out")
	}
	if !strings.Contains(out, "fail") {
		t.Error("expected 'fail' entry in output")
	}
}

func TestRun_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	r := pipeline.New(pipeline.Options{Format: "xml"}, &buf)
	if err := r.Run([]string{"<log/>"}); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestRun_TailN(t *testing.T) {
	lines := []string{
		`{"time":"2024-01-01T10:00:00Z","msg":"a"}`,
		`{"time":"2024-01-01T11:00:00Z","msg":"b"}`,
		`{"time":"2024-01-01T12:00:00Z","msg":"c"}`,
	}

	var buf bytes.Buffer
	r := pipeline.New(pipeline.Options{
		Format:    "json",
		TimeField: "time",
		TailN:     1,
		OutputFmt: output.FormatJSON,
	}, &buf)

	if err := r.Run(lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, `"a"`) || strings.Contains(out, `"b"`) {
		t.Error("expected only last entry in tail output")
	}
	if !strings.Contains(out, `"c"`) {
		t.Error("expected last entry 'c' in tail output")
	}
}

func TestRun_FieldFilter(t *testing.T) {
	lines := []string{
		`{"time":"2024-01-01T10:00:00Z","level":"info","msg":"ok"}`,
		`{"time":"2024-01-01T11:00:00Z","level":"error","msg":"bad"}`,
	}

	var buf bytes.Buffer
	r := pipeline.New(pipeline.Options{
		Format:    "json",
		TimeField: "time",
		Fields:    map[string]string{"level": "error"},
		OutputFmt: output.FormatJSON,
	}, &buf)

	if err := r.Run(lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "ok") {
		t.Error("expected info entry to be filtered out")
	}
	if !strings.Contains(out, "bad") {
		t.Error("expected error entry in output")
	}
}
