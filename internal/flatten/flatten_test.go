package flatten

import (
	"testing"
)

func makeEntry(ts int64, fields map[string]interface{}) Entry {
	return Entry{Timestamp: ts, Fields: fields}
}

func TestRun_FlattensNestedMap(t *testing.T) {
	input := []Entry{
		makeEntry(1, map[string]interface{}{
			"level": "info",
			"http": map[string]interface{}{
				"method": "GET",
				"status": 200,
			},
		}),
	}
	out := Run(input, DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["http.method"] != "GET" {
		t.Errorf("expected http.method=GET, got %v", out[0].Fields["http.method"])
	}
	if out[0].Fields["http.status"] != 200 {
		t.Errorf("expected http.status=200, got %v", out[0].Fields["http.status"])
	}
	if _, ok := out[0].Fields["http"]; ok {
		t.Error("expected nested 'http' key to be removed")
	}
}

func TestRun_CustomSeparator(t *testing.T) {
	input := []Entry{
		makeEntry(2, map[string]interface{}{
			"db": map[string]interface{}{"host": "localhost"},
		}),
	}
	opts := Options{Separator: "_", MaxDepth: 0}
	out := Run(input, opts)
	if out[0].Fields["db_host"] != "localhost" {
		t.Errorf("expected db_host=localhost, got %v", out[0].Fields["db_host"])
	}
}

func TestRun_MaxDepthLimitsRecursion(t *testing.T) {
	input := []Entry{
		makeEntry(3, map[string]interface{}{
			"a": map[string]interface{}{
				"b": map[string]interface{}{"c": "deep"},
			},
		}),
	}
	opts := Options{Separator: ".", MaxDepth: 1}
	out := Run(input, opts)
	// At depth 1 the nested map should be stringified, not further expanded.
	if _, ok := out[0].Fields["a.b.c"]; ok {
		t.Error("expected a.b.c to NOT be present when MaxDepth=1")
	}
	if _, ok := out[0].Fields["a.b"]; !ok {
		t.Error("expected a.b to be present when MaxDepth=1")
	}
}

func TestRun_FlatEntryUnchanged(t *testing.T) {
	input := []Entry{
		makeEntry(4, map[string]interface{}{"msg": "hello", "level": "warn"}),
	}
	out := Run(input, DefaultOptions())
	if out[0].Fields["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out[0].Fields["msg"])
	}
}

func TestRun_EmptyInput(t *testing.T) {
	out := Run([]Entry{}, DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}

func TestRun_PreservesTimestamp(t *testing.T) {
	input := []Entry{makeEntry(9999, map[string]interface{}{"x": "y"})}
	out := Run(input, DefaultOptions())
	if out[0].Timestamp != 9999 {
		t.Errorf("expected timestamp 9999, got %d", out[0].Timestamp)
	}
}
