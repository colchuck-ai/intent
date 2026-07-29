package tree_test

import (
	"errors"
	"testing"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
)

// The Phase 0 seed is the shared read-command fixture.
const seedPath = "../model/testdata/seed.intent.yaml"

func loadSeed(t *testing.T) *tree.Index {
	t.Helper()
	r, err := model.Load(seedPath)
	if err != nil {
		t.Fatalf("loading seed: %v", err)
	}
	return tree.Build(r)
}

func parseIndex(t *testing.T, yaml string) *tree.Index {
	t.Helper()
	r, err := model.Parse([]byte(yaml))
	if err != nil {
		t.Fatalf("parsing inline yaml: %v", err)
	}
	return tree.Build(r)
}

func TestBuildOrderAndAddresses(t *testing.T) {
	ix := loadSeed(t)

	// A spot-check of tree order: roots first, children in insertion order.
	want := []string{
		"product",
		"product.jobs.understand_the_rationale_behind_an_element",
		"product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup",
		"product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.risks.ambiguous_reference",
		"product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.stable_logical_ids",
		"product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.resolvable_references",
		"product.decision_records.prefer_explicit_declaration",
		"engineering",
		"engineering.principles.derive_only_mechanical_duals",
		"engineering.constraints.single_go_binary",
		"engineering.components.validator",
		"engineering.components.generator",
		"engineering.decision_records.single_go_binary_over_python",
		"change_records.add_resolvable_references_requirement",
	}
	all := ix.All()
	if len(all) != len(want) {
		var got []string
		for _, e := range all {
			got = append(got, e.Addr)
		}
		t.Fatalf("element count: got %d %v, want %d", len(all), got, len(want))
	}
	for i, w := range want {
		if all[i].Addr != w {
			t.Errorf("order[%d]: got %q, want %q", i, all[i].Addr, w)
		}
	}
}

func TestElementFields(t *testing.T) {
	ix := loadSeed(t)

	oc, ok := ix.Get("product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup")
	if !ok {
		t.Fatal("expected the outcome to be present")
	}
	if oc.Kind != tree.KindOutcome {
		t.Errorf("kind: got %q, want outcome", oc.Kind)
	}
	if !oc.IsDocument {
		t.Error("outcome is marked type: document in the seed")
	}
	if oc.Domain != tree.DomainProduct {
		t.Errorf("domain: got %q, want product", oc.Domain)
	}

	// A principle carries both a name and a statement.
	pr, ok := ix.Get("engineering.principles.derive_only_mechanical_duals")
	if !ok {
		t.Fatal("expected the principle to be present")
	}
	if pr.Name != "Derive Only Mechanical Duals" {
		t.Errorf("principle name = %q", pr.Name)
	}
	if pr.Prose == "" {
		t.Error("principle prose (the statement value) should be populated")
	}
}

func TestEdgesForwardAndIncoming(t *testing.T) {
	ix := loadSeed(t)

	riskAddr := "product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.risks.ambiguous_reference"

	// Two requirements mitigate the same risk; both should appear incoming.
	refs := ix.Incoming(riskAddr)
	if len(refs) != 2 {
		t.Fatalf("expected 2 incoming refs on the risk, got %d: %+v", len(refs), refs)
	}
	for _, r := range refs {
		if r.Kind != tree.EdgeMitigates {
			t.Errorf("incoming ref kind: got %q, want mitigates", r.Kind)
		}
	}

	// The validator component fulfills a requirement and relates to the generator.
	v, _ := ix.Get("engineering.components.validator")
	var kinds []tree.EdgeKind
	for _, e := range v.Edges {
		kinds = append(kinds, e.Kind)
	}
	if len(kinds) != 2 || kinds[0] != tree.EdgeFulfills || kinds[1] != tree.EdgeRelationships {
		t.Errorf("validator edges: got %v, want [fulfills relationships]", kinds)
	}
	if v.Edges[1].Note == "" {
		t.Error("relationships edge should carry its note")
	}
}

func TestResolveFullAddress(t *testing.T) {
	ix := loadSeed(t)
	addr := "engineering.components.validator"
	e, err := ix.Resolve(addr)
	if err != nil {
		t.Fatalf("resolve full address: %v", err)
	}
	if e.Addr != addr {
		t.Errorf("got %q, want %q", e.Addr, addr)
	}
}

func TestResolveUnambiguousSuffix(t *testing.T) {
	ix := loadSeed(t)
	e, err := ix.Resolve("stable_logical_ids")
	if err != nil {
		t.Fatalf("resolve suffix: %v", err)
	}
	want := "product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.stable_logical_ids"
	if e.Addr != want {
		t.Errorf("got %q, want %q", e.Addr, want)
	}
}

func TestResolveNotFound(t *testing.T) {
	ix := loadSeed(t)
	_, err := ix.Resolve("does_not_exist")
	var nf *tree.NotFoundError
	if !errors.As(err, &nf) {
		t.Fatalf("expected NotFoundError, got %v", err)
	}
}

// TestResolvePartialSegmentIsNotASuffix guards the segment-boundary rule: a
// substring of a segment must not match.
func TestResolvePartialSegmentIsNotASuffix(t *testing.T) {
	ix := loadSeed(t)
	if _, err := ix.Resolve("logical_ids"); err == nil {
		t.Fatal("expected 'logical_ids' not to match segment 'stable_logical_ids'")
	}
}

// TestResolveAmbiguous uses a tree with an outcome key repeated under two jobs,
// so the bare suffix is ambiguous but a longer suffix disambiguates.
func TestResolveAmbiguous(t *testing.T) {
	ix := parseIndex(t, ambiguousTree)

	_, err := ix.Resolve("shared")
	var amb *tree.AmbiguousError
	if !errors.As(err, &amb) {
		t.Fatalf("expected AmbiguousError, got %v", err)
	}
	if len(amb.Candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %v", amb.Candidates)
	}
	// Candidates are sorted for stable output.
	if amb.Candidates[0] >= amb.Candidates[1] {
		t.Errorf("candidates should be sorted: %v", amb.Candidates)
	}

	// A longer suffix disambiguates.
	e, err := ix.Resolve("job_a.outcomes.shared")
	if err != nil {
		t.Fatalf("longer suffix should resolve: %v", err)
	}
	if e.Addr != "product.jobs.job_a.outcomes.shared" {
		t.Errorf("got %q", e.Addr)
	}
}

func TestFindByQueryTypeDomain(t *testing.T) {
	ix := loadSeed(t)

	// Query matches prose/name.
	got := ix.Find("reference", "", "")
	if len(got) == 0 {
		t.Fatal("expected matches for 'reference'")
	}

	// Type filter alone lists all of a kind.
	comps := ix.Find("", tree.KindComponent, "")
	if len(comps) != 2 {
		t.Fatalf("expected 2 components, got %d", len(comps))
	}

	// Domain filter narrows to engineering.
	eng := ix.Find("", "", tree.DomainEngineering)
	for _, e := range eng {
		if e.Domain != tree.DomainEngineering {
			t.Errorf("domain filter leaked %q", e.Addr)
		}
	}

	// Combined filters intersect.
	none := ix.Find("nonexistent-token-xyz", tree.KindComponent, "")
	if len(none) != 0 {
		t.Errorf("expected no matches, got %d", len(none))
	}
}

func TestTraceDepthAndDirection(t *testing.T) {
	ix := loadSeed(t)
	reqAddr := "product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.resolvable_references"

	// Downstream depth 1: the requirement's own outgoing edges (mitigates a risk,
	// depends on another requirement).
	down := ix.Trace(reqAddr, 1, false, true, nil)
	if down[0].Addr != reqAddr || down[0].Depth != 0 {
		t.Fatalf("first node should be the start at depth 0, got %+v", down[0])
	}
	if len(down) != 3 {
		t.Fatalf("expected start + 2 downstream, got %d: %+v", len(down), down)
	}

	// Upstream: who references the risk. The risk is mitigated by two requirements.
	riskAddr := "product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.risks.ambiguous_reference"
	up := ix.Trace(riskAddr, 1, true, false, nil)
	if len(up) != 3 {
		t.Fatalf("expected start + 2 upstream refs, got %d: %+v", len(up), up)
	}
	for _, n := range up[1:] {
		if !n.Up {
			t.Errorf("upstream node %q should be marked Up", n.Addr)
		}
	}
}

func TestTraceEdgeKindFilter(t *testing.T) {
	ix := loadSeed(t)
	reqAddr := "product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.resolvable_references"

	// Restrict to dependsOn only: start + one dependency.
	only := ix.Trace(reqAddr, 1, false, true, map[tree.EdgeKind]bool{tree.EdgeDependsOn: true})
	if len(only) != 2 {
		t.Fatalf("expected start + 1 dependsOn node, got %d: %+v", len(only), only)
	}
	if only[1].Via != tree.EdgeDependsOn {
		t.Errorf("expected dependsOn edge, got %q", only[1].Via)
	}
}

func TestParseKindAndDomain(t *testing.T) {
	if k, err := tree.ParseKind(""); err != nil || k != "" {
		t.Errorf("empty type should be no-filter: %q %v", k, err)
	}
	if k, err := tree.ParseKind("component"); err != nil || k != tree.KindComponent {
		t.Errorf("valid type: %q %v", k, err)
	}
	if _, err := tree.ParseKind("widget"); err == nil {
		t.Error("unknown type should error")
	}
	if d, err := tree.ParseDomain("engineering"); err != nil || d != tree.DomainEngineering {
		t.Errorf("valid domain: %q %v", d, err)
	}
	if _, err := tree.ParseDomain("marketing"); err == nil {
		t.Error("unknown domain should error")
	}
}

const ambiguousTree = `
product:
  name: P
  summary: s
  jobs:
    job_a:
      name: A
      story: story a
      outcomes:
        shared:
          name: Shared A
          statement: stmt
    job_b:
      name: B
      story: story b
      outcomes:
        shared:
          name: Shared B
          statement: stmt
engineering:
  name: E
  summary: s
`
