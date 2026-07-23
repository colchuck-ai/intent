package cli

import (
	"testing"

	"github.com/colchuck-ai/intent/internal/help"
	"github.com/colchuck-ai/intent/internal/tree"
	"github.com/colchuck-ai/intent/internal/validate"
)

// TestFooterJudgmentSlugsResolve guards the DESIGN §12 promise that every
// footer points at a real topic. It covers each kind's nudge plus the fallbacks
// used by add/edit/record footers.
func TestFooterJudgmentSlugsResolve(t *testing.T) {
	slugs := map[string]bool{}
	for _, s := range kindJudgment {
		slugs[s] = true
	}
	// Fallbacks not necessarily present in kindJudgment.
	slugs[addFooter(tree.KindCR, "x").judgment] = true // fallback path
	slugs[editFooter("x").judgment] = true
	slugs[recordFooter("x").judgment] = true

	for s := range slugs {
		if _, ok := help.Resolve(s); !ok {
			t.Errorf("footer judgment slug %q does not resolve to a topic", s)
		}
	}
}

// TestEveryErrorCodeHasATopic guards that every validator code an operator can
// hit ends its message with a `→ intent help E0NN` that actually resolves.
func TestEveryErrorCodeHasATopic(t *testing.T) {
	for _, code := range []validate.Code{validate.E001, validate.E002, validate.E003, validate.E004} {
		if _, ok := help.Resolve(string(code)); !ok {
			t.Errorf("error code %s has no help topic", code)
		}
	}
}
