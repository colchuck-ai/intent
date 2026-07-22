package cli

import (
	"bytes"
	"strings"
	"testing"
)

const seedFlag = "../model/testdata/seed.intent.yaml"

// run executes the CLI with args and captures combined output. cobra writes both
// command output and errors to the buffers we set here.
func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := newRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestTreeListsEveryElement(t *testing.T) {
	out, err := run(t, "tree", "-f", seedFlag)
	if err != nil {
		t.Fatalf("tree: %v\n%s", err, out)
	}
	for _, want := range []string{
		"product [product] Intent",
		"validator [component] Validator (document)",
		"add_resolvable_references_requirement [cr]",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("tree output missing %q\n%s", want, out)
		}
	}
}

func TestShowResolvesSuffixAndEchoesFullAddress(t *testing.T) {
	out, err := run(t, "show", "resolvable_references", "-f", seedFlag)
	if err != nil {
		t.Fatalf("show: %v\n%s", err, out)
	}
	full := "product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.resolvable_references"
	if !strings.HasPrefix(out, full) {
		t.Errorf("show should echo the full address first:\n%s", out)
	}
	// Outgoing edges are shown, with the target's name annotated.
	if !strings.Contains(out, "→ mitigates") || !strings.Contains(out, "(Ambiguous Reference)") {
		t.Errorf("show missing annotated outgoing edge:\n%s", out)
	}
	// Incoming refs are NOT part of show (that's affects/trace).
	if strings.Contains(out, "referenced by") {
		t.Errorf("show should not list incoming refs:\n%s", out)
	}
}

func TestShowAmbiguousListsCandidates(t *testing.T) {
	out, err := run(t, "show", "shared", "-f", "testdata/ambiguous.intent.yaml")
	if err == nil {
		t.Fatalf("expected an error for an ambiguous suffix\n%s", out)
	}
	if !strings.Contains(out, "ambiguous") ||
		!strings.Contains(out, "product.jobs.job_a.outcomes.shared") ||
		!strings.Contains(out, "product.jobs.job_b.outcomes.shared") {
		t.Errorf("ambiguity error should list both candidates:\n%s", out)
	}
}

func TestShowNotFound(t *testing.T) {
	out, err := run(t, "show", "nope_nope", "-f", seedFlag)
	if err == nil {
		t.Fatalf("expected a not-found error\n%s", out)
	}
	if !strings.Contains(out, "no element matches") {
		t.Errorf("expected not-found message:\n%s", out)
	}
}

func TestFindFilters(t *testing.T) {
	// Type filter with no query lists all of a kind.
	out, err := run(t, "find", "--type", "component", "-f", seedFlag)
	if err != nil {
		t.Fatalf("find: %v\n%s", err, out)
	}
	if !strings.Contains(out, "engineering.components.validator") ||
		!strings.Contains(out, "engineering.components.generator") {
		t.Errorf("find --type component should list both components:\n%s", out)
	}
	if strings.Contains(out, "[job]") || strings.Contains(out, "[risk]") {
		t.Errorf("find --type component leaked other kinds:\n%s", out)
	}

	// No match reports plainly.
	out, err = run(t, "find", "zzz-no-such-token", "-f", seedFlag)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if strings.TrimSpace(out) != "no matches" {
		t.Errorf("expected 'no matches', got:\n%s", out)
	}
}

func TestAffectsShowsIncoming(t *testing.T) {
	out, err := run(t, "affects", "ambiguous_reference", "-f", seedFlag)
	if err != nil {
		t.Fatalf("affects: %v\n%s", err, out)
	}
	if !strings.Contains(out, "referenced by:") {
		t.Fatalf("expected a referenced-by section:\n%s", out)
	}
	// Both requirements mitigate this risk.
	if strings.Count(out, "← mitigates") != 2 {
		t.Errorf("expected two incoming mitigates refs:\n%s", out)
	}
}

func TestTraceDownstreamAndUpstream(t *testing.T) {
	// Downstream only: the requirement's own edges.
	out, err := run(t, "trace", "resolvable_references", "--down", "-f", seedFlag)
	if err != nil {
		t.Fatalf("trace: %v\n%s", err, out)
	}
	if !strings.Contains(out, "→ mitigates") || !strings.Contains(out, "→ dependsOn") {
		t.Errorf("downstream trace missing outgoing edges:\n%s", out)
	}
	if strings.Contains(out, "←") {
		t.Errorf("--down should not show upstream refs:\n%s", out)
	}

	// Edge-kind filter restricts traversal.
	out, err = run(t, "trace", "resolvable_references", "--down", "--edges", "dependsOn", "-f", seedFlag)
	if err != nil {
		t.Fatalf("trace: %v\n%s", err, out)
	}
	if strings.Contains(out, "mitigates") {
		t.Errorf("--edges dependsOn should exclude mitigates:\n%s", out)
	}
	if !strings.Contains(out, "dependsOn") {
		t.Errorf("--edges dependsOn should keep dependsOn:\n%s", out)
	}
}

func TestFindUnknownTypeErrors(t *testing.T) {
	out, err := run(t, "find", "--type", "widget", "-f", seedFlag)
	if err == nil {
		t.Fatalf("expected an error for an unknown --type\n%s", out)
	}
	if !strings.Contains(out, "unknown type") {
		t.Errorf("expected an unknown-type message:\n%s", out)
	}
}

// TestFindSnippetShowsMatchedField guards that the excerpt surfaces the matched
// term rather than an unrelated clip.
func TestFindSnippetShowsMatchedField(t *testing.T) {
	// "markdown" is in the generator component's responsibility (its prose), so
	// the snippet should include it.
	out, err := run(t, "find", "markdown", "-f", seedFlag)
	if err != nil {
		t.Fatalf("find: %v\n%s", err, out)
	}
	if !strings.Contains(strings.ToLower(out), "markdown") {
		t.Errorf("snippet should include the matched term:\n%s", out)
	}
}

func TestMissingFileErrors(t *testing.T) {
	_, err := run(t, "tree", "-f", "testdata/does-not-exist.yaml")
	if err == nil {
		t.Fatal("expected an error for a missing file")
	}
}
