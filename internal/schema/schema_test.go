package schema_test

import (
	"os"
	"strings"
	"testing"

	"github.com/colchuck-ai/intent/internal/schema"
)

const seedPath = "../model/testdata/seed.intent.yaml"

// TestSeedValidates is the Phase 0 exit criterion for the schema: the seed
// fixture satisfies the structural contract.
func TestSeedValidates(t *testing.T) {
	b, err := os.ReadFile(seedPath)
	if err != nil {
		t.Fatalf("reading seed: %v", err)
	}
	if err := schema.Validate(b); err != nil {
		t.Fatalf("seed should be valid, got: %v", err)
	}
}

// TestRejectsUnknownField confirms additionalProperties:false bites.
func TestRejectsUnknownField(t *testing.T) {
	bad := `
product:
  name: X
  summary: Y
  mystery: nope
engineering:
  name: X
  summary: Y
`
	if err := schema.Validate([]byte(bad)); err == nil {
		t.Fatal("expected an unknown field to be rejected")
	}
}

// TestRejectsBadKey confirms snake_case key enforcement.
func TestRejectsBadKey(t *testing.T) {
	bad := `
product:
  name: X
  summary: Y
  jobs:
    BadKey:
      name: X
      story: Y
engineering:
  name: X
  summary: Y
`
	if err := schema.Validate([]byte(bad)); err == nil {
		t.Fatal("expected a non-snake_case key to be rejected")
	}
}

// TestRejectsRequirementWithoutMitigates confirms the minItems:1 rule that
// replaces the old "requirement without a risk" linter check.
func TestRejectsRequirementWithoutMitigates(t *testing.T) {
	bad := `
product:
  name: X
  summary: Y
  jobs:
    j:
      name: J
      story: S
      outcomes:
        o:
          name: O
          statement: S
          requirements:
            r:
              name: R
              statement: S
engineering:
  name: X
  summary: Y
`
	err := schema.Validate([]byte(bad))
	if err == nil {
		t.Fatal("expected a requirement without mitigates to be rejected")
	}
	if !strings.Contains(err.Error(), "mitigates") {
		t.Errorf("error should mention the missing field, got: %v", err)
	}
}

// TestRejectsMissingRoot confirms both roots are required.
func TestRejectsMissingRoot(t *testing.T) {
	bad := `
product:
  name: X
  summary: Y
`
	if err := schema.Validate([]byte(bad)); err == nil {
		t.Fatal("expected a missing engineering root to be rejected")
	}
}
