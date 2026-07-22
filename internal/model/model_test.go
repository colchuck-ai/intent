package model_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/colchuck-ai/intent/internal/model"
)

const seedPath = "testdata/seed.intent.yaml"

// TestRoundTripByteIdentical is the Phase 0 exit criterion: loading the seed and
// re-serializing it reproduces the file exactly. Because the CLI is the sole
// writer and always emits canonical form, this guarantees clean diffs.
//
// Run with UPDATE_SEED=1 to rewrite the fixture into canonical form.
func TestRoundTripByteIdentical(t *testing.T) {
	original, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("reading seed: %v", err)
	}

	r, err := model.Parse(original)
	if err != nil {
		t.Fatalf("parsing seed: %v", err)
	}
	got, err := r.Serialize()
	if err != nil {
		t.Fatalf("serializing: %v", err)
	}

	if os.Getenv("UPDATE_SEED") == "1" {
		if err := os.WriteFile(seedPath, got, 0o644); err != nil {
			t.Fatalf("updating seed: %v", err)
		}
		t.Skip("seed rewritten to canonical form")
	}

	if string(got) != string(original) {
		t.Errorf("re-serialized seed is not byte-identical to the committed fixture.\n"+
			"Run UPDATE_SEED=1 go test ./internal/model to canonicalize it, then review the diff.\n"+
			"--- got ---\n%s", got)
	}
}

// TestSerializeIsIdempotent confirms a second round trip changes nothing, so
// canonical output is a fixed point.
func TestSerializeIsIdempotent(t *testing.T) {
	r, err := model.Load(seedPath)
	if err != nil {
		t.Fatalf("loading seed: %v", err)
	}
	first, err := r.Serialize()
	if err != nil {
		t.Fatalf("first serialize: %v", err)
	}
	r2, err := model.Parse(first)
	if err != nil {
		t.Fatalf("re-parsing: %v", err)
	}
	second, err := r2.Serialize()
	if err != nil {
		t.Fatalf("second serialize: %v", err)
	}
	if string(first) != string(second) {
		t.Errorf("serialize is not idempotent:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}

// TestOrderPreserved checks that map entries keep their file order after a round
// trip — the property OrderedMap exists to guarantee.
func TestOrderPreserved(t *testing.T) {
	r, err := model.Load(seedPath)
	if err != nil {
		t.Fatalf("loading seed: %v", err)
	}
	job, ok := r.Product.Jobs.Get("understand_the_rationale_behind_an_element")
	if !ok {
		t.Fatal("expected the seed job to be present")
	}
	outcome, ok := job.Outcomes.Get("fast_rationale_lookup")
	if !ok {
		t.Fatal("expected the seed outcome to be present")
	}
	got := outcome.Requirements.Keys()
	want := []string{"stable_logical_ids", "resolvable_references"}
	if len(got) != len(want) {
		t.Fatalf("requirement count: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("requirement order at %d: got %q, want %q", i, got[i], want[i])
		}
	}
}

// TestSeedFixtureExists is a guard so the fixture path stays valid.
func TestSeedFixtureExists(t *testing.T) {
	if _, err := os.Stat(filepath.FromSlash(seedPath)); err != nil {
		t.Fatalf("seed fixture missing: %v", err)
	}
}
