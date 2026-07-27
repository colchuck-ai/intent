package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// tmpSeed copies the seed fixture into a temp file the test can mutate, and
// returns its path.
func tmpSeed(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("../model/testdata/seed.intent.yaml")
	if err != nil {
		t.Fatalf("reading seed: %v", err)
	}
	path := filepath.Join(t.TempDir(), "intent.yaml")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatalf("writing temp seed: %v", err)
	}
	return path
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v", path, err)
	}
	return string(b)
}

func TestAddWritesCanonicalValidElement(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "add", "requirement", "fast_rationale_lookup", "speedy",
		"--name", "Speedy", "--statement", "Must be fast.",
		"--mitigates", "ambiguous_reference", "-f", path)
	if err != nil {
		t.Fatalf("add: %v\n%s", err, out)
	}
	if !strings.Contains(out, "added product.jobs.") || !strings.Contains(out, "requirements.speedy") {
		t.Errorf("missing added confirmation:\n%s", out)
	}
	// The footer carries the proactive judgment nudge for the kind (DESIGN §12).
	if !strings.Contains(out, "judgment: intent help judgment:requirement-vs-task") {
		t.Errorf("missing requirement judgment footer:\n%s", out)
	}
	// The element is persisted and resolvable.
	show, err := run(t, "show", "speedy", "-f", path)
	if err != nil {
		t.Fatalf("show after add: %v\n%s", err, show)
	}
	if !strings.Contains(show, "[requirement] Speedy") {
		t.Errorf("added element not shown:\n%s", show)
	}
}

// TestRepeatableFlagsDontSplitOnComma guards against a pflag footgun:
// StringSliceVar treats a comma inside a single value as another list item, so
// a natural-English sentence passed to a repeatable flag like --acceptance
// would silently shatter into fragments. These flags must use StringArrayVar
// instead, which only splits on repeated flag occurrences.
func TestRepeatableFlagsDontSplitOnComma(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "add", "requirement", "fast_rationale_lookup", "speedy",
		"--name", "Speedy", "--statement", "Must be fast.",
		"--mitigates", "ambiguous_reference",
		"--acceptance", "Covers add, set, and link in one pass.",
		"-f", path)
	if err != nil {
		t.Fatalf("add: %v\n%s", err, out)
	}
	show, err := run(t, "show", "speedy", "-f", path)
	if err != nil {
		t.Fatalf("show after add: %v\n%s", err, show)
	}
	if !strings.Contains(show, "Covers add, set, and link in one pass.") {
		t.Errorf("acceptance criterion was split on commas instead of kept whole:\n%s", show)
	}
}

// TestWritesAreCanonical adds then removes an element, expecting the file to
// return byte-identical to the seed — proof every write re-serializes canonically.
func TestWritesAreCanonical(t *testing.T) {
	path := tmpSeed(t)
	before := readFile(t, path)

	if _, err := run(t, "add", "risk", "fast_rationale_lookup", "temp_risk",
		"--name", "Temp", "--statement", "x", "-f", path); err != nil {
		t.Fatalf("add: %v", err)
	}
	if _, err := run(t, "rm", "temp_risk", "-f", path); err != nil {
		t.Fatalf("rm: %v", err)
	}
	if got := readFile(t, path); got != before {
		t.Errorf("add+rm did not round-trip to the original seed:\n--- got ---\n%s", got)
	}
}

func TestAddRefusesInvalidWrite(t *testing.T) {
	path := tmpSeed(t)
	before := readFile(t, path)
	// A requirement with no --mitigates is rejected before any write.
	out, err := run(t, "add", "requirement", "fast_rationale_lookup", "bad",
		"--name", "Bad", "--statement", "x", "-f", path)
	if err == nil {
		t.Fatalf("expected an error for a requirement without mitigates\n%s", out)
	}
	if readFile(t, path) != before {
		t.Error("a rejected add must not modify the file")
	}
}

func TestSetUpdatesField(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "set", "ambiguous_reference", "statement", "A brand new statement.", "-f", path)
	if err != nil {
		t.Fatalf("set: %v\n%s", err, out)
	}
	show, _ := run(t, "show", "ambiguous_reference", "-f", path)
	if !strings.Contains(show, "A brand new statement.") {
		t.Errorf("set field not persisted:\n%s", show)
	}
}

func TestSetRejectsEdgeField(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "set", "resolvable_references", "mitigates", "ambiguous_reference", "-f", path)
	if err == nil {
		t.Fatalf("expected an error setting an edge field via set\n%s", out)
	}
	if !strings.Contains(out, "link") {
		t.Errorf("error should point to `link`:\n%s", out)
	}
}

func TestDetailNudgeFires(t *testing.T) {
	path := tmpSeed(t)
	// resolvable_references has no outgoing edge to the validator component, so a
	// detail that interpolates it should nudge to declare one.
	ref := "engineering.components.validator"
	out, err := run(t, "set", "resolvable_references", "detail",
		"See {{"+ref+".link}} for context.", "-f", path)
	if err != nil {
		t.Fatalf("set detail: %v\n%s", err, out)
	}
	if !strings.Contains(out, "nudge:") || !strings.Contains(out, "no declared edge covers it") {
		t.Errorf("expected a detail nudge:\n%s", out)
	}
	if !strings.Contains(out, "intent link <kind>") {
		t.Errorf("nudge should suggest declaring the edge:\n%s", out)
	}
}

// TestDetailNudgeSilentWhenEdgeExists: resolvable_references already depends on
// stable_logical_ids, so a detail that references it should NOT nudge.
func TestDetailNudgeSilentWhenEdgeExists(t *testing.T) {
	path := tmpSeed(t)
	ref := "product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.stable_logical_ids"
	out, err := run(t, "set", "resolvable_references", "detail",
		"Builds on {{"+ref+".link}}.", "-f", path)
	if err != nil {
		t.Fatalf("set detail: %v\n%s", err, out)
	}
	if strings.Contains(out, "nudge:") {
		t.Errorf("no nudge expected when an edge already covers the reference:\n%s", out)
	}
}

func TestRmGuardedByIncomingRefs(t *testing.T) {
	path := tmpSeed(t)
	before := readFile(t, path)
	// Two requirements mitigate this risk, so removing it would dangle them.
	out, err := run(t, "rm", "ambiguous_reference", "-f", path)
	if err == nil {
		t.Fatalf("expected rm to be refused\n%s", out)
	}
	if !strings.Contains(out, "still referenced by") {
		t.Errorf("expected a referenced-by refusal:\n%s", out)
	}
	if readFile(t, path) != before {
		t.Error("a refused rm must not modify the file")
	}
}

func TestLinkAndUnlink(t *testing.T) {
	path := tmpSeed(t)
	// stable_logical_ids has no dependsOn initially; link one.
	out, err := run(t, "link", "dependsOn", "stable_logical_ids", "resolvable_references", "-f", path)
	if err != nil {
		t.Fatalf("link: %v\n%s", err, out)
	}
	show, _ := run(t, "show", "stable_logical_ids", "-f", path)
	if !strings.Contains(show, "→ dependsOn") {
		t.Errorf("link not persisted:\n%s", show)
	}

	// Unlink removes it.
	if _, err := run(t, "unlink", "dependsOn", "stable_logical_ids", "resolvable_references", "-f", path); err != nil {
		t.Fatalf("unlink: %v", err)
	}
	show, _ = run(t, "show", "stable_logical_ids", "-f", path)
	if strings.Contains(show, "dependsOn") {
		t.Errorf("unlink did not remove the edge:\n%s", show)
	}
}

func TestLinkUnknownEdgeKind(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "link", "bogus", "stable_logical_ids", "ambiguous_reference", "-f", path)
	if err == nil {
		t.Fatalf("expected an unknown-edge-kind error\n%s", out)
	}
	if !strings.Contains(out, "unknown edge kind") {
		t.Errorf("expected an unknown-edge-kind message:\n%s", out)
	}
}

// TestUnlinkDanglingEdgeBySuffix covers the hand-edit recovery path: a dangling
// edge (its target no longer exists) can still be removed with a suffix, after
// which the tree is clean and the write goes through.
func TestUnlinkDanglingEdgeBySuffix(t *testing.T) {
	b, err := os.ReadFile("testdata/dangling_edge.intent.yaml")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}
	path := filepath.Join(t.TempDir(), "intent.yaml")
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatalf("writing fixture: %v", err)
	}

	// The dangling target "ghost" does not resolve, but unlink must still remove
	// the edge given the suffix.
	out, err := run(t, "unlink", "dependsOn", "q", "ghost", "-f", path)
	if err != nil {
		t.Fatalf("unlink dangling: %v\n%s", err, out)
	}
	// The tree is now clean, so validate passes.
	if vout, verr := run(t, "validate", "-f", path); verr != nil {
		t.Errorf("validate after unlink should pass:\n%s", vout)
	}
}

func TestAddRejectsInvalidType(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "add", "outcome", "understand_the_rationale_behind_an_element", "oc",
		"--name", "O", "--statement", "s", "--type", "folder", "-f", path)
	if err == nil {
		t.Fatalf("expected an invalid --type error\n%s", out)
	}
	if !strings.Contains(out, "inline or document") {
		t.Errorf("expected an inline-or-document message:\n%s", out)
	}
}

func TestAddUnknownTypeErrors(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "add", "widget", "product", "x", "-f", path)
	if err == nil {
		t.Fatalf("expected an unknown-type error\n%s", out)
	}
	if !strings.Contains(out, "cannot add") {
		t.Errorf("expected a cannot-add message:\n%s", out)
	}
}
