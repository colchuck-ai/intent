package cli

import (
	"fmt"
	"io"

	"github.com/colchuck-ai/intent/internal/interp"
	"github.com/colchuck-ai/intent/internal/tree"
)

// footer is the scaffolded command-output footer (DESIGN §12): a suggested next
// command and the judgment topic to consult before trusting the edit. The
// judgment slugs are placeholders — Phase 6 authors the real topics and wires
// them; keeping them in this one file means that phase edits only here.
type footer struct {
	next     string // a full command line, e.g. "intent show product.jobs.x"
	judgment string // a help slug, e.g. "judgment:requirement-vs-task"
}

func printFooter(w io.Writer, f footer) {
	if f.next != "" {
		fmt.Fprintf(w, "\nnext: %s\n", f.next)
	}
	if f.judgment != "" {
		fmt.Fprintf(w, "judgment: intent help %s\n", f.judgment)
	}
}

// kindJudgment maps a newly-added element's kind to the judgment test worth
// meeting at the moment of creation (DESIGN §12 — mutating verbs proactively
// nudge). Provisional slugs pending Phase 6.
var kindJudgment = map[tree.Kind]string{
	tree.KindJob:         "judgment:job-vs-activity",
	tree.KindOutcome:     "judgment:outcome-vs-solution",
	tree.KindRisk:        "judgment:risk-vs-feature",
	tree.KindRequirement: "judgment:requirement-vs-task",
	tree.KindComponent:   "judgment:altitude",
	tree.KindPrinciple:   "judgment:materiality",
	tree.KindConstraint:  "judgment:materiality",
}

// addFooter is the footer for a freshly-added element: show it, and read the
// judgment test for its kind.
func addFooter(kind tree.Kind, addr string) footer {
	j := kindJudgment[kind]
	if j == "" {
		j = "judgment:coherence"
	}
	return footer{next: "intent show " + addr, judgment: j}
}

// editFooter is the footer for an edit to an existing element
// (set/link/unlink/promote/mv).
func editFooter(addr string) footer {
	return footer{next: "intent show " + addr, judgment: "judgment:coherence"}
}

// recordFooter is the footer for a freshly-created record: show it, and read the
// test for whether the decision/change was worth recording.
func recordFooter(addr string) footer {
	return footer{next: "intent show " + addr, judgment: "judgment:worth-recording"}
}

// nudgeDetail warns when the just-written element's detail interpolates a
// resolvable element address that no declared edge covers (DESIGN §2 / §14). It
// is advisory only — no lint, nothing in CI — pointing at `intent link` so a
// real dependency gets declared rather than left implicit in prose.
func nudgeDetail(w io.Writer, ix *tree.Index, addr string) {
	e, ok := ix.Get(addr)
	if !ok || e.Detail == "" {
		return
	}
	edgeTargets := map[string]bool{}
	for _, ed := range e.Edges {
		edgeTargets[ed.To] = true
	}
	seen := map[string]bool{}
	for _, ref := range interp.Addresses(e.Detail, e.Addr) {
		if seen[ref.Addr] {
			continue
		}
		seen[ref.Addr] = true
		if _, ok := ix.Get(ref.Addr); !ok {
			continue // unresolvable → that's an E001, not a nudge
		}
		if edgeTargets[ref.Addr] {
			continue // already declared as an edge
		}
		fmt.Fprintf(w, "\nnudge: detail references %s but no declared edge covers it.\n", ref.Addr)
		fmt.Fprintf(w, "      if that's a real dependency, declare it: intent link <kind> %s %s\n", addr, ref.Addr)
	}
}
