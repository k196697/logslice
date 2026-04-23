package dedup

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/filter"
)

func makeAdapterEntry(msg string) filter.Entry {
	return filter.Entry{
		Timestamp: time.Now(),
		Fields:    map[string]string{"msg": msg, "level": "info"},
	}
}

func TestParseConfig_NilWhenNoFieldsAndNotGlobal(t *testing.T) {
	opts, err := ParseConfig("", false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts != nil {
		t.Errorf("expected nil opts, got %+v", opts)
	}
}

func TestParseConfig_GlobalMode(t *testing.T) {
	opts, err := ParseConfig("", false, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts == nil {
		t.Fatal("expected non-nil opts for global mode")
	}
	if opts.ConsecutiveOnly {
		t.Errorf("expected ConsecutiveOnly=false")
	}
}

func TestParseConfig_WithFields(t *testing.T) {
	opts, err := ParseConfig("msg, level", true, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts == nil {
		t.Fatal("expected non-nil opts")
	}
	if len(opts.Fields) != 2 || opts.Fields[0] != "msg" || opts.Fields[1] != "level" {
		t.Errorf("unexpected fields: %v", opts.Fields)
	}
	if !opts.ConsecutiveOnly {
		t.Errorf("expected ConsecutiveOnly=true")
	}
}

func TestParseConfig_EmptyFieldInList(t *testing.T) {
	_, err := ParseConfig("msg,,level", false, false)
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestRunFromConfig_NilOpts(t *testing.T) {
	entries := []filter.Entry{
		makeAdapterEntry("hello"),
		makeAdapterEntry("hello"),
	}
	result := RunFromConfig(entries, nil)
	if len(result) != 2 {
		t.Errorf("expected 2 entries unchanged, got %d", len(result))
	}
}

func TestRunFromConfig_RemovesDuplicates(t *testing.T) {
	opts, _ := ParseConfig("", false, true)
	entries := []filter.Entry{
		makeAdapterEntry("dup"),
		makeAdapterEntry("dup"),
		makeAdapterEntry("unique"),
	}
	result := RunFromConfig(entries, opts)
	if len(result) != 2 {
		t.Errorf("expected 2 entries after dedup, got %d", len(result))
	}
}
