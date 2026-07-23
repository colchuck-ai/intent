package cli

import (
	"strings"
	"testing"
)

// TestPromoteCommandCascades: promoting a requirement under a fresh inline
// outcome promotes the outcome with it and reports the cascade.
func TestPromoteCommandCascades(t *testing.T) {
	path := tmpSeed(t)
	if _, err := run(t, "add", "outcome", "understand_the_rationale_behind_an_element", "inline_oc",
		"--name", "O", "--statement", "s", "-f", path); err != nil {
		t.Fatal(err)
	}
	if _, err := run(t, "add", "requirement", "inline_oc", "deep",
		"--name", "D", "--statement", "s", "--mitigates", "ambiguous_reference", "-f", path); err != nil {
		t.Fatal(err)
	}
	out, err := run(t, "promote", "deep", "-f", path)
	if err != nil {
		t.Fatalf("promote: %v\n%s", err, out)
	}
	if !strings.Contains(out, "cascade:") || !strings.Contains(out, "outcomes.inline_oc") {
		t.Errorf("expected a cascade line for the inline outcome:\n%s", out)
	}
	if !strings.Contains(out, "promoted") {
		t.Errorf("expected a promoted confirmation:\n%s", out)
	}
}

// TestSetTypeInlineGuardedThenCascade drives demotion through the CLI: refused
// with a document descendant, allowed with --cascade.
func TestSetTypeInlineGuardedThenCascade(t *testing.T) {
	path := tmpSeed(t)
	if _, err := run(t, "promote", "stable_logical_ids", "-f", path); err != nil {
		t.Fatal(err)
	}
	out, err := run(t, "set", "fast_rationale_lookup", "type", "inline", "-f", path)
	if err == nil {
		t.Fatalf("expected demotion to be refused:\n%s", out)
	}
	if !strings.Contains(out, "--cascade") {
		t.Errorf("refusal should mention --cascade:\n%s", out)
	}
	out, err = run(t, "set", "fast_rationale_lookup", "type", "inline", "--cascade", "-f", path)
	if err != nil {
		t.Fatalf("cascade demote: %v\n%s", err, out)
	}
	if !strings.Contains(out, "cascade:") {
		t.Errorf("expected a cascade report:\n%s", out)
	}
}

// TestMvRenameRoundTrips: renaming a component and renaming it back returns the
// file byte-for-byte to the seed — proof the reference cascade and canonical
// re-serialize are reversible.
func TestMvRenameRoundTrips(t *testing.T) {
	path := tmpSeed(t)
	before := readFile(t, path)

	if _, err := run(t, "mv", "generator", "--rename", "doc_generator", "-f", path); err != nil {
		t.Fatalf("rename: %v", err)
	}
	if readFile(t, path) == before {
		t.Fatal("rename should have changed the file")
	}
	if _, err := run(t, "mv", "doc_generator", "--rename", "generator", "-f", path); err != nil {
		t.Fatalf("rename back: %v", err)
	}
	if got := readFile(t, path); got != before {
		t.Errorf("rename round-trip did not restore the seed:\n--- got ---\n%s", got)
	}
}

// TestMvRenameRewritesReferences: after a rename, the old address resolves
// nowhere and the tree still validates (no dangling references).
func TestMvRenameRewritesReferences(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "mv", "generator", "--rename", "doc_generator", "-f", path)
	if err != nil {
		t.Fatalf("rename: %v\n%s", err, out)
	}
	if !strings.Contains(out, "rewrote 2 reference(s)") {
		t.Errorf("expected a reference-rewrite count:\n%s", out)
	}
	if vout, verr := run(t, "validate", "-f", path); verr != nil {
		t.Errorf("validate should pass after rename:\n%s", vout)
	}
	// The old address is gone.
	if _, err := run(t, "show", "engineering.components.generator", "-f", path); err == nil {
		t.Error("the old address should no longer resolve")
	}
}

// TestMvReparentValidates: re-parenting a requirement keeps the tree valid.
func TestMvReparentValidates(t *testing.T) {
	path := tmpSeed(t)
	if _, err := run(t, "add", "outcome", "understand_the_rationale_behind_an_element", "second_outcome",
		"--name", "S", "--statement", "s", "-f", path); err != nil {
		t.Fatal(err)
	}
	out, err := run(t, "mv", "resolvable_references", "second_outcome", "-f", path)
	if err != nil {
		t.Fatalf("reparent: %v\n%s", err, out)
	}
	if !strings.Contains(out, "moved") || !strings.Contains(out, "second_outcome.requirements.resolvable_references") {
		t.Errorf("expected a moved confirmation to the new parent:\n%s", out)
	}
	if vout, verr := run(t, "validate", "-f", path); verr != nil {
		t.Errorf("validate should pass after reparent:\n%s", vout)
	}
}

// TestMvReorderResolvesSiblingSuffix: --before/--after accept the shortest
// unambiguous suffix, like every other address argument (DESIGN §11).
func TestMvReorderResolvesSiblingSuffix(t *testing.T) {
	path := tmpSeed(t)
	// Position resolvable_references before its sibling, named by a dotted suffix.
	out, err := run(t, "mv", "resolvable_references",
		"--before", "fast_rationale_lookup.requirements.stable_logical_ids", "-f", path)
	if err != nil {
		t.Fatalf("mv reorder by suffix: %v\n%s", err, out)
	}
	if !strings.Contains(out, "reordered") {
		t.Errorf("expected a reordered confirmation:\n%s", out)
	}
	// The order actually changed: resolvable_references now leads.
	tree, _ := run(t, "tree", "-f", path)
	iRes := strings.Index(tree, "resolvable_references")
	iStable := strings.Index(tree, "stable_logical_ids")
	if iRes < 0 || iStable < 0 || iRes > iStable {
		t.Errorf("resolvable_references should now precede stable_logical_ids:\n%s", tree)
	}
}

func TestMvReorderUnknownSiblingErrors(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "mv", "resolvable_references", "--before", "no_such_sibling", "-f", path)
	if err == nil {
		t.Fatalf("expected an error for an unknown sibling:\n%s", out)
	}
	if !strings.Contains(out, "no element matches") {
		t.Errorf("expected a not-found resolution error:\n%s", out)
	}
}

func TestRecordCreatesValidRecords(t *testing.T) {
	path := tmpSeed(t)
	out, err := run(t, "record", "adr", "use_cobra",
		"--name", "Use Cobra", "--affects", "engineering",
		"--summary", "Use cobra for the command tree.", "-f", path)
	if err != nil {
		t.Fatalf("record adr: %v\n%s", err, out)
	}
	if !strings.Contains(out, "recorded engineering.decision_records.use_cobra") {
		t.Errorf("missing recorded confirmation:\n%s", out)
	}
	if !strings.Contains(out, "judgment:worth-recording") {
		t.Errorf("missing worth-recording judgment footer:\n%s", out)
	}
	show, _ := run(t, "show", "use_cobra", "-f", path)
	if !strings.Contains(show, "[adr] Use Cobra") {
		t.Errorf("recorded adr not shown:\n%s", show)
	}
}

func TestRecordCrossDomainRefused(t *testing.T) {
	path := tmpSeed(t)
	before := readFile(t, path)
	out, err := run(t, "record", "pdr", "bad",
		"--name", "Bad", "--affects", "engineering", "--summary", "x", "-f", path)
	if err == nil {
		t.Fatalf("expected a cross-domain refusal:\n%s", out)
	}
	if readFile(t, path) != before {
		t.Error("a refused record must not modify the file")
	}
}

func TestSetTypeStillRoutesThroughCascade(t *testing.T) {
	path := tmpSeed(t)
	// `set <addr> type document` behaves like promote (DESIGN §5).
	out, err := run(t, "set", "generator", "type", "document", "-f", path)
	if err != nil {
		t.Fatalf("set type document: %v\n%s", err, out)
	}
	show, _ := run(t, "show", "generator", "-f", path)
	if !strings.Contains(show, "(document)") {
		t.Errorf("component should now render as a document:\n%s", show)
	}
}
