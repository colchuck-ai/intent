package mutate_test

import (
	"os"
	"strings"
	"testing"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/mutate"
	"github.com/colchuck-ai/intent/internal/tree"
)

const seedPath = "../model/testdata/seed.intent.yaml"

func loadSeed(t *testing.T) *model.Root {
	t.Helper()
	b, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("reading seed: %v", err)
	}
	r, err := model.Parse(b)
	if err != nil {
		t.Fatalf("parsing seed: %v", err)
	}
	return r
}

const (
	jobAddr     = "product.jobs.understand_the_rationale_behind_an_element"
	outcomeAddr = jobAddr + ".outcomes.fast_rationale_lookup"
	riskAddr    = outcomeAddr + ".risks.ambiguous_reference"
	reqAddr     = outcomeAddr + ".requirements.resolvable_references"
)

func TestAddJobUnderProduct(t *testing.T) {
	r := loadSeed(t)
	addr, err := mutate.Add(r, tree.KindJob, "product", "new_job", mutate.Fields{
		Name: "New Job", Story: "When I do X, I want Y, so Z.",
	})
	if err != nil {
		t.Fatalf("add job: %v", err)
	}
	if addr != "product.jobs.new_job" {
		t.Errorf("addr = %q", addr)
	}
	got, ok := r.Product.Jobs.Get("new_job")
	if !ok || got.Name != "New Job" {
		t.Errorf("job not inserted: %+v", got)
	}
	// Insertion order: the new job comes last.
	keys := r.Product.Jobs.Keys()
	if keys[len(keys)-1] != "new_job" {
		t.Errorf("new job should be appended last; keys = %v", keys)
	}
}

func TestAddOutcomeUnderJob(t *testing.T) {
	r := loadSeed(t)
	addr, err := mutate.Add(r, tree.KindOutcome, jobAddr, "new_outcome", mutate.Fields{
		Name: "New Outcome", Statement: "Minimize the time to X.",
	})
	if err != nil {
		t.Fatalf("add outcome: %v", err)
	}
	if addr != jobAddr+".outcomes.new_outcome" {
		t.Errorf("addr = %q", addr)
	}
	job, _ := r.Product.Jobs.Get("understand_the_rationale_behind_an_element")
	if _, ok := job.Outcomes.Get("new_outcome"); !ok {
		t.Error("outcome not inserted into the job")
	}
}

func TestAddRequirementRequiresMitigates(t *testing.T) {
	r := loadSeed(t)
	_, err := mutate.Add(r, tree.KindRequirement, outcomeAddr, "new_req", mutate.Fields{
		Name: "New Req", Statement: "Intent must X.",
	})
	if err == nil || !strings.Contains(err.Error(), "mitigates") {
		t.Fatalf("expected a mitigates-required error, got %v", err)
	}
}

func TestAddRequirementWithEdges(t *testing.T) {
	r := loadSeed(t)
	addr, err := mutate.Add(r, tree.KindRequirement, outcomeAddr, "new_req", mutate.Fields{
		Name: "New Req", Statement: "Intent must X.",
		Mitigates:  []string{riskAddr},
		Acceptance: []string{"It works."},
	})
	if err != nil {
		t.Fatalf("add requirement: %v", err)
	}
	oc, _ := loadOutcome(t, r)
	q, ok := oc.Requirements.Get("new_req")
	if !ok {
		t.Fatal("requirement not inserted")
	}
	if len(q.Mitigates) != 1 || q.Mitigates[0] != riskAddr {
		t.Errorf("mitigates = %v", q.Mitigates)
	}
	if len(q.AcceptanceCriteria) != 1 {
		t.Errorf("acceptance_criteria = %v", q.AcceptanceCriteria)
	}
	_ = addr
}

func TestAddRejectsWrongParent(t *testing.T) {
	r := loadSeed(t)
	// A component cannot live under a job.
	_, err := mutate.Add(r, tree.KindComponent, jobAddr, "nope", mutate.Fields{
		Responsibility: "X", Fulfills: []string{reqAddr},
	})
	if err == nil {
		t.Fatal("expected a parent-mismatch error")
	}
}

func TestAddRejectsDuplicateKey(t *testing.T) {
	r := loadSeed(t)
	_, err := mutate.Add(r, tree.KindOutcome, jobAddr, "fast_rationale_lookup", mutate.Fields{
		Name: "Dup", Statement: "X",
	})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected an already-exists error, got %v", err)
	}
}

func TestAddRejectsBadKey(t *testing.T) {
	r := loadSeed(t)
	_, err := mutate.Add(r, tree.KindJob, "product", "Not A Key", mutate.Fields{Name: "X", Story: "Y"})
	if err == nil || !strings.Contains(err.Error(), "lowercase letters") {
		t.Fatalf("expected a bad-key-shape error, got %v", err)
	}
}

func TestSetScalarField(t *testing.T) {
	r := loadSeed(t)
	if err := mutate.Set(r, riskAddr, "statement", "A new statement."); err != nil {
		t.Fatalf("set: %v", err)
	}
	oc, _ := loadOutcome(t, r)
	risk, _ := oc.Risks.Get("ambiguous_reference")
	if risk.Statement != "A new statement." {
		t.Errorf("statement = %q", risk.Statement)
	}
}

func TestSetRejectsEdgeField(t *testing.T) {
	r := loadSeed(t)
	err := mutate.Set(r, reqAddr, "mitigates", riskAddr)
	if err == nil || !strings.Contains(err.Error(), "link") {
		t.Fatalf("expected a use-link error, got %v", err)
	}
}

func TestSetRejectsUnknownField(t *testing.T) {
	r := loadSeed(t)
	err := mutate.Set(r, riskAddr, "bogus", "x")
	if err == nil || !strings.Contains(err.Error(), "settable") {
		t.Fatalf("expected an unsettable-field error, got %v", err)
	}
}

func TestSetPrincipleStatement(t *testing.T) {
	r := loadSeed(t)
	addr := "engineering.principles.derive_only_mechanical_duals"
	if err := mutate.Set(r, addr, "statement", "New axiom."); err != nil {
		t.Fatalf("set principle: %v", err)
	}
	v, _ := r.Engineering.Principles.Get("derive_only_mechanical_duals")
	if v.Statement != "New axiom." {
		t.Errorf("principle statement = %q", v.Statement)
	}
}

func TestSetPrincipleName(t *testing.T) {
	r := loadSeed(t)
	addr := "engineering.principles.derive_only_mechanical_duals"
	if err := mutate.Set(r, addr, "name", "Renamed Axiom"); err != nil {
		t.Fatalf("set principle name: %v", err)
	}
	v, _ := r.Engineering.Principles.Get("derive_only_mechanical_duals")
	if v.Name != "Renamed Axiom" {
		t.Errorf("principle name = %q", v.Name)
	}
}

func TestAddPrincipleWithName(t *testing.T) {
	r := loadSeed(t)
	addr, err := mutate.Add(r, tree.KindPrinciple, "engineering", "new_rule",
		mutate.Fields{Name: "New Rule", Statement: "A fresh axiom."})
	if err != nil {
		t.Fatalf("add principle: %v", err)
	}
	v, _ := r.Engineering.Principles.Get("new_rule")
	if v.Name != "New Rule" || v.Statement != "A fresh axiom." {
		t.Errorf("principle = %+v", v)
	}
	if addr != "engineering.principles.new_rule" {
		t.Errorf("addr = %q", addr)
	}
}

func TestLinkAndUnlink(t *testing.T) {
	r := loadSeed(t)
	// stable_logical_ids initially has no dependsOn; link one.
	stable := outcomeAddr + ".requirements.stable_logical_ids"
	if err := mutate.Link(r, tree.EdgeDependsOn, stable, reqAddr, ""); err != nil {
		t.Fatalf("link: %v", err)
	}
	oc, _ := loadOutcome(t, r)
	q, _ := oc.Requirements.Get("stable_logical_ids")
	if len(q.DependsOn) != 1 || q.DependsOn[0] != reqAddr {
		t.Fatalf("dependsOn after link = %v", q.DependsOn)
	}

	// Duplicate link is refused.
	if err := mutate.Link(r, tree.EdgeDependsOn, stable, reqAddr, ""); err == nil {
		t.Error("expected a duplicate-edge error")
	}

	// Unlink removes it.
	if err := mutate.Unlink(r, tree.EdgeDependsOn, stable, reqAddr); err != nil {
		t.Fatalf("unlink: %v", err)
	}
	oc, _ = loadOutcome(t, r)
	q, _ = oc.Requirements.Get("stable_logical_ids")
	if len(q.DependsOn) != 0 {
		t.Errorf("dependsOn after unlink = %v", q.DependsOn)
	}

	// Unlinking a missing edge errors.
	if err := mutate.Unlink(r, tree.EdgeDependsOn, stable, reqAddr); err == nil {
		t.Error("expected a no-such-edge error")
	}
}

func TestLinkRejectsWrongEdgeKind(t *testing.T) {
	r := loadSeed(t)
	// A requirement can't fulfill anything (that's a component edge).
	err := mutate.Link(r, tree.EdgeFulfills, reqAddr, riskAddr, "")
	if err == nil || !strings.Contains(err.Error(), "cannot declare") {
		t.Fatalf("expected a cannot-declare error, got %v", err)
	}
}

func TestLinkRelationshipWithNote(t *testing.T) {
	r := loadSeed(t)
	from := "engineering.components.generator"
	to := "engineering.components.validator"
	if err := mutate.Link(r, tree.EdgeRelationships, from, to, "shares the resolver"); err != nil {
		t.Fatalf("link relationship: %v", err)
	}
	c, _ := r.Engineering.Components.Get("generator")
	note, ok := c.Relationships.Get(to)
	if !ok || note != "shares the resolver" {
		t.Errorf("relationship note = %q (ok=%v)", note, ok)
	}
}

func TestRemoveLeaf(t *testing.T) {
	r := loadSeed(t)
	// stable_logical_ids has no incoming refs from outside — a plain removal.
	addr := outcomeAddr + ".requirements.stable_logical_ids"
	if err := mutate.Remove(r, addr); err != nil {
		t.Fatalf("remove: %v", err)
	}
	oc, _ := loadOutcome(t, r)
	if _, ok := oc.Requirements.Get("stable_logical_ids"); ok {
		t.Error("requirement should be gone")
	}
}

func TestRemoveChangeRecord(t *testing.T) {
	r := loadSeed(t)
	addr := "change_records.add_resolvable_references_requirement"
	if err := mutate.Remove(r, addr); err != nil {
		t.Fatalf("remove CR: %v", err)
	}
	if _, ok := r.ChangeRecords.Get("add_resolvable_references_requirement"); ok {
		t.Error("change record should be gone")
	}
}

func TestRemoveMissing(t *testing.T) {
	r := loadSeed(t)
	err := mutate.Remove(r, outcomeAddr+".requirements.does_not_exist")
	if err == nil {
		t.Fatal("expected a not-found error")
	}
}

// loadOutcome is a helper to fetch the seed's single outcome after a mutation.
func loadOutcome(t *testing.T, r *model.Root) (model.Outcome, bool) {
	t.Helper()
	job, ok := r.Product.Jobs.Get("understand_the_rationale_behind_an_element")
	if !ok {
		t.Fatal("seed job missing")
	}
	return job.Outcomes.Get("fast_rationale_lookup")
}
