package mutate_test

import (
	"strings"
	"testing"

	"github.com/colchuck-ai/intent/internal/mutate"
	"github.com/colchuck-ai/intent/internal/tree"
	"github.com/colchuck-ai/intent/internal/validate"
)

func contains(list []string, want string) bool {
	for _, x := range list {
		if x == want {
			return true
		}
	}
	return false
}

// TestPromoteCascadesInlineAncestor: promoting a requirement under an inline
// outcome must promote the outcome too, keeping containment valid (DESIGN §5).
func TestPromoteCascadesInlineAncestor(t *testing.T) {
	r := loadSeed(t)
	ocAddr, err := mutate.Add(r, tree.KindOutcome, jobAddr, "inline_oc", mutate.Fields{Name: "O", Statement: "s"})
	if err != nil {
		t.Fatal(err)
	}
	reqA, err := mutate.Add(r, tree.KindRequirement, ocAddr, "deep", mutate.Fields{
		Name: "D", Statement: "s", Mitigates: []string{riskAddr},
	})
	if err != nil {
		t.Fatal(err)
	}

	changed, err := mutate.SetType(r, tree.Build(r), reqA, "document", false)
	if err != nil {
		t.Fatalf("promote: %v", err)
	}
	if !contains(changed, reqA) || !contains(changed, ocAddr) {
		t.Errorf("expected cascade to promote %s and %s; got %v", reqA, ocAddr, changed)
	}
	oc, _ := tree.Build(r).Get(ocAddr)
	if !oc.IsDocument {
		t.Error("inline outcome ancestor was not promoted")
	}
}

// TestPromoteAlreadyDocumentIsNoOp: the seed outcome is already a document, so
// promoting it changes nothing.
func TestPromoteAlreadyDocumentIsNoOp(t *testing.T) {
	r := loadSeed(t)
	changed, err := mutate.SetType(r, tree.Build(r), outcomeAddr, "document", false)
	if err != nil {
		t.Fatalf("promote: %v", err)
	}
	if len(changed) != 0 {
		t.Errorf("expected no changes; got %v", changed)
	}
}

// TestDemoteGuardedThenCascade: demoting an outcome with a document descendant is
// refused, then allowed (and cascades) with --cascade.
func TestDemoteGuardedThenCascade(t *testing.T) {
	r := loadSeed(t)
	stable := outcomeAddr + ".requirements.stable_logical_ids"
	if _, err := mutate.SetType(r, tree.Build(r), stable, "document", false); err != nil {
		t.Fatalf("promote descendant: %v", err)
	}

	if _, err := mutate.SetType(r, tree.Build(r), outcomeAddr, "inline", false); err == nil {
		t.Fatal("expected demotion to be blocked by the document descendant")
	}

	changed, err := mutate.SetType(r, tree.Build(r), outcomeAddr, "inline", true)
	if err != nil {
		t.Fatalf("cascade demote: %v", err)
	}
	if !contains(changed, stable) {
		t.Errorf("cascade should also demote %s; got %v", stable, changed)
	}
	ix := tree.Build(r)
	if oc, _ := ix.Get(outcomeAddr); oc.IsDocument {
		t.Error("outcome should be inline after cascade demote")
	}
	if st, _ := ix.Get(stable); st.IsDocument {
		t.Error("descendant should be inline after cascade demote")
	}
}

func TestSetTypeRejectsNonPromotable(t *testing.T) {
	r := loadSeed(t)
	if _, err := mutate.SetType(r, tree.Build(r), riskAddr, "document", false); err == nil {
		t.Fatal("expected an error promoting a risk (only outcome/requirement/component carry a type)")
	}
}

func TestMvReorderKeepsAddress(t *testing.T) {
	r := loadSeed(t)
	newAddr, n, _, err := mutate.Mv(r, tree.Build(r), reqAddr, mutate.MvOpts{Before: "stable_logical_ids"})
	if err != nil {
		t.Fatalf("mv reorder: %v", err)
	}
	if newAddr != reqAddr {
		t.Errorf("a reorder must not change the address; got %s", newAddr)
	}
	if n != 0 {
		t.Errorf("a reorder rewrites no references; got %d", n)
	}
	oc, _ := loadOutcome(t, r)
	if keys := oc.Requirements.Keys(); keys[0] != "resolvable_references" {
		t.Errorf("expected resolvable_references moved first; got %v", keys)
	}
}

// TestMvRenameCascadesReferences: renaming a component rewrites both the
// relationships-map key and the prose {{ }} address that point at it, leaving the
// tree valid.
func TestMvRenameCascadesReferences(t *testing.T) {
	r := loadSeed(t)
	gen := "engineering.components.generator"
	newAddr, n, _, err := mutate.Mv(r, tree.Build(r), gen, mutate.MvOpts{Rename: "doc_generator"})
	if err != nil {
		t.Fatalf("mv rename: %v", err)
	}
	if newAddr != "engineering.components.doc_generator" {
		t.Errorf("newAddr = %s", newAddr)
	}
	if n != 2 {
		t.Errorf("expected 2 references rewritten (relationships key + behavior prose); got %d", n)
	}
	if fs := validate.Check(tree.Build(r)); len(fs) > 0 {
		t.Errorf("tree should stay valid after rename; findings: %v", fs)
	}
}

// TestMvReparentRewritesIncomingEdges: moving a requirement under a new outcome
// rewrites the edges that referenced it (a component's fulfills, a CR's affects),
// so nothing dangles.
func TestMvReparentRewritesIncomingEdges(t *testing.T) {
	r := loadSeed(t)
	oc2, err := mutate.Add(r, tree.KindOutcome, jobAddr, "second_outcome", mutate.Fields{Name: "S", Statement: "s"})
	if err != nil {
		t.Fatal(err)
	}
	newAddr, n, _, err := mutate.Mv(r, tree.Build(r), reqAddr, mutate.MvOpts{NewParent: oc2})
	if err != nil {
		t.Fatalf("mv reparent: %v", err)
	}
	want := oc2 + ".requirements.resolvable_references"
	if newAddr != want {
		t.Errorf("newAddr = %s, want %s", newAddr, want)
	}
	if n != 2 {
		t.Errorf("expected 2 incoming edges rewritten; got %d", n)
	}
	if fs := validate.Check(tree.Build(r)); len(fs) > 0 {
		t.Errorf("tree should stay valid after reparent; findings: %v", fs)
	}
}

func TestMvReparentRejectsWrongParent(t *testing.T) {
	r := loadSeed(t)
	// A requirement cannot be re-parented under a job (only under an outcome).
	_, _, _, err := mutate.Mv(r, tree.Build(r), reqAddr, mutate.MvOpts{NewParent: jobAddr})
	if err == nil || !strings.Contains(err.Error(), "cannot be re-parented") {
		t.Fatalf("expected a wrong-parent error, got %v", err)
	}
}

func TestMvReparentRejectsFixedLocationKind(t *testing.T) {
	r := loadSeed(t)
	// A component has a fixed home; re-parenting it under engineering's... it can
	// only differ by domain, which canNest refuses.
	comp := "engineering.components.validator"
	_, _, _, err := mutate.Mv(r, tree.Build(r), comp, mutate.MvOpts{NewParent: "product"})
	if err == nil {
		t.Fatal("expected re-parenting a component to be refused")
	}
}

// TestMvReparentPromotesDestinationForDocument: moving a document requirement
// under an inline outcome must promote that outcome (containment by
// construction, DESIGN §5) rather than leave an E002 for the linter to catch.
func TestMvReparentPromotesDestinationForDocument(t *testing.T) {
	r := loadSeed(t)
	// A fresh inline outcome to receive the move.
	oc2, err := mutate.Add(r, tree.KindOutcome, jobAddr, "inline_target", mutate.Fields{Name: "T", Statement: "s"})
	if err != nil {
		t.Fatal(err)
	}
	// stable_logical_ids is a document after promotion; move it under the inline outcome.
	if _, err := mutate.SetType(r, tree.Build(r), outcomeAddr+".requirements.stable_logical_ids", "document", false); err != nil {
		t.Fatal(err)
	}
	_, _, promoted, err := mutate.Mv(r, tree.Build(r), outcomeAddr+".requirements.stable_logical_ids", mutate.MvOpts{NewParent: oc2})
	if err != nil {
		t.Fatalf("mv reparent of a document: %v", err)
	}
	if !contains(promoted, oc2) {
		t.Errorf("expected the inline destination outcome %s to be promoted; got %v", oc2, promoted)
	}
	// The tree must be valid with no E002 — the linter never had to catch it.
	if fs := validate.Check(tree.Build(r)); len(fs) > 0 {
		t.Errorf("containment must hold by construction; findings: %v", fs)
	}
	if oc, _ := tree.Build(r).Get(oc2); !oc.IsDocument {
		t.Error("destination outcome should now render as a document")
	}
}

func TestAddRecordPDRAndCR(t *testing.T) {
	r := loadSeed(t)
	addr, err := mutate.AddRecord(r, tree.KindPDR, "new_pdr", mutate.RecordFields{
		Name: "New PDR", Affects: []string{"product"}, Summary: "A product decision.",
	})
	if err != nil {
		t.Fatalf("record pdr: %v", err)
	}
	if addr != "product.decision_records.new_pdr" {
		t.Errorf("pdr addr = %s", addr)
	}

	craddr, err := mutate.AddRecord(r, tree.KindCR, "new_cr", mutate.RecordFields{
		Name: "New CR", Affects: []string{"product", "engineering"},
		Change: "Changed X.", Rationale: "Because Y.",
	})
	if err != nil {
		t.Fatalf("record cr: %v", err)
	}
	if craddr != "change_records.new_cr" {
		t.Errorf("cr addr = %s", craddr)
	}
	if fs := validate.Check(tree.Build(r)); len(fs) > 0 {
		t.Errorf("records should be valid; findings: %v", fs)
	}
}

func TestAddRecordDomainPurity(t *testing.T) {
	r := loadSeed(t)
	_, err := mutate.AddRecord(r, tree.KindPDR, "bad", mutate.RecordFields{
		Name: "Bad", Affects: []string{"engineering"}, Summary: "x",
	})
	if err == nil || !strings.Contains(err.Error(), "product") {
		t.Fatalf("expected a domain-purity error, got %v", err)
	}
}

func TestAddRecordRequiresFields(t *testing.T) {
	r := loadSeed(t)
	// A cr with no change/rationale.
	if _, err := mutate.AddRecord(r, tree.KindCR, "cr", mutate.RecordFields{
		Name: "C", Affects: []string{"product"},
	}); err == nil {
		t.Error("expected a cr to require --change and --rationale")
	}
	// A decision record with no affects.
	if _, err := mutate.AddRecord(r, tree.KindPDR, "pdr", mutate.RecordFields{
		Name: "P", Summary: "s",
	}); err == nil {
		t.Error("expected a record to require at least one --affects")
	}
}
