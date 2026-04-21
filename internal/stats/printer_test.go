package stats

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestPrint_ContainsExpectedLines(t *testing.T) {
	s := Summary{
		Total:    10,
		Matched:  7,
		Skipped:  3,
		Earliest: time.Unix(0, 0).UTC(),
		Latest:   time.Unix(3600, 0).UTC(),
		Fields:   map[string]int{"level": 7, "host": 4},
	}

	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	for _, want := range []string{
		"Total lines:", "10",
		"Matched:", "7",
		"Skipped:", "3",
		"Earliest:",
		"Latest:",
		"level", "host",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\nfull output:\n%s", want, out)
		}
	}
}

func TestPrint_NoTimestampsWhenEmpty(t *testing.T) {
	s := Summary{
		Total:   5,
		Matched: 0,
		Skipped: 5,
		Fields:  map[string]int{},
	}

	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	if strings.Contains(out, "Earliest") {
		t.Error("should not print Earliest when no entries matched")
	}
}

func TestPrint_FieldsSorted(t *testing.T) {
	s := Summary{
		Total:    3,
		Matched:  3,
		Earliest: time.Unix(1, 0),
		Latest:   time.Unix(2, 0),
		Fields:   map[string]int{"zebra": 1, "alpha": 2, "mango": 1},
	}

	var buf bytes.Buffer
	Print(&buf, s)
	out := buf.String()

	alphaIdx := strings.Index(out, "alpha")
	mangoIdx := strings.Index(out, "mango")
	zebraIdx := strings.Index(out, "zebra")

	if !(alphaIdx < mangoIdx && mangoIdx < zebraIdx) {
		t.Errorf("fields not sorted alphabetically in output:\n%s", out)
	}
}
