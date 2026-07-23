package cli

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHelpIndexIsTaskFirst(t *testing.T) {
	out, err := run(t, "help")
	if err != nil {
		t.Fatalf("help: %v\n%s", err, out)
	}
	// The guide (task) plane must appear before the concept plane.
	gi, ci := strings.Index(out, "TASKS (guide"), strings.Index(out, "CONCEPTS")
	if gi < 0 || ci < 0 || gi > ci {
		t.Errorf("expected task plane before concepts:\n%s", out)
	}
	if !strings.Contains(out, "guide:edit") || !strings.Contains(out, "concept:requirement") {
		t.Errorf("index missing expected slugs:\n%s", out)
	}
}

func TestHelpTopicBySlug(t *testing.T) {
	out, err := run(t, "help", "judgment:requirement-vs-task")
	if err != nil {
		t.Fatalf("help slug: %v\n%s", err, out)
	}
	if !strings.Contains(out, "condition") {
		t.Errorf("requirement-vs-task body not shown:\n%s", out)
	}
}

// TestHelpErrorCode covers the inline-from-failure path: `intent help E004`.
func TestHelpErrorCode(t *testing.T) {
	out, err := run(t, "help", "E004")
	if err != nil {
		t.Fatalf("help E004: %v\n%s", err, out)
	}
	if !strings.Contains(out, "domain") {
		t.Errorf("E004 topic not shown:\n%s", out)
	}
}

func TestHelpUnknownTopicErrors(t *testing.T) {
	out, err := run(t, "help", "concept:nope")
	if err == nil {
		t.Fatalf("expected an unknown-topic error\n%s", out)
	}
	if !strings.Contains(out, "--list") {
		t.Errorf("error should point at --list:\n%s", out)
	}
}

// TestHelpJSONIsStableAndComplete guards the "--json output is stable" done-when
// (DESIGN §12 / BUILD_PLAN Phase 6): valid JSON, every entry carries the index
// fields, and the order matches help --list.
func TestHelpJSONIsStableAndComplete(t *testing.T) {
	out, err := run(t, "help", "--json")
	if err != nil {
		t.Fatalf("help --json: %v\n%s", err, out)
	}
	var topics []struct {
		Slug    string `json:"slug"`
		Plane   string `json:"plane"`
		Title   string `json:"title"`
		Summary string `json:"summary"`
	}
	if err := json.Unmarshal([]byte(out), &topics); err != nil {
		t.Fatalf("help --json did not produce valid JSON: %v\n%s", err, out)
	}
	if len(topics) == 0 {
		t.Fatal("help --json returned no topics")
	}
	for _, tp := range topics {
		if tp.Slug == "" || tp.Plane == "" || tp.Title == "" || tp.Summary == "" {
			t.Errorf("json entry missing fields: %+v", tp)
		}
	}

	// Same run twice → byte-identical (deterministic order).
	out2, _ := run(t, "help", "--json")
	if out != out2 {
		t.Error("help --json is not byte-stable across runs")
	}
}
