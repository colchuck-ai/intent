package skill

import (
	"sort"
	"strings"
	"testing"

	"github.com/colchuck-ai/intent/internal/help"
)

// TestRenderJudgmentsMatchHelp guards the one-liner extraction: every
// judgment:* topic in help appears exactly once, carrying its real slug and
// summary rather than a duplicated or hand-copied body.
func TestRenderJudgmentsMatchHelp(t *testing.T) {
	want := help.InPlane("judgment")
	if len(want) == 0 {
		t.Fatal("no judgment topics loaded from help; test fixture assumption broken")
	}

	got := Render().Judgments
	if len(got) != len(want) {
		t.Fatalf("got %d judgment pointers, want %d", len(got), len(want))
	}
	for i, w := range want {
		if got[i].Slug != w.Slug {
			t.Errorf("index %d: slug = %q, want %q", i, got[i].Slug, w.Slug)
		}
		if got[i].Summary != w.Summary {
			t.Errorf("index %d (%s): summary = %q, want %q", i, w.Slug, got[i].Summary, w.Summary)
		}
		if !strings.HasPrefix(got[i].Slug, "judgment:") {
			t.Errorf("index %d: slug %q does not carry the judgment: plane prefix", i, got[i].Slug)
		}
	}
}

// TestRenderJudgmentsStableOrder guards determinism: repeated Render() calls
// return the same order, and that order is sorted by slug (help.InPlane
// already sorts within a plane; this pins the contract from skill's side too
// so a future change to how help orders topics fails loudly here).
func TestRenderJudgmentsStableOrder(t *testing.T) {
	first := Render().Judgments
	second := Render().Judgments

	if len(first) != len(second) {
		t.Fatalf("Render() returned %d judgments then %d", len(first), len(second))
	}
	if !sort.SliceIsSorted(first, func(i, j int) bool { return first[i].Slug < first[j].Slug }) {
		t.Error("judgment pointers are not sorted by slug")
	}
	for i := range first {
		if first[i] != second[i] {
			t.Errorf("index %d differs across calls: %+v vs %+v", i, first[i], second[i])
		}
	}
}

// TestContentMarkdownStructure guards the assembled body shape: the required
// sections appear in order, every judgment pointer is represented by its
// summary and slug, and the full topic body is never duplicated in.
func TestContentMarkdownStructure(t *testing.T) {
	c := Render()
	md := c.Markdown()

	sections := []string{"# Intent", "## Entry point", "## Gotchas", "## Judgment"}
	last := -1
	for _, s := range sections {
		idx := strings.Index(md, s)
		if idx < 0 {
			t.Fatalf("markdown missing section %q", s)
		}
		if idx <= last {
			t.Errorf("section %q out of order", s)
		}
		last = idx
	}

	for _, g := range c.Gotchas {
		if !strings.Contains(md, g) {
			t.Errorf("markdown missing gotcha %q", g)
		}
	}

	for _, j := range c.Judgments {
		if !strings.Contains(md, j.Summary) {
			t.Errorf("markdown missing judgment summary for %s", j.Slug)
		}
		if !strings.Contains(md, "`intent help "+j.Slug+"`") {
			t.Errorf("markdown missing pointer to %s", j.Slug)
		}
	}
}

// TestContentMarkdownExcludesFullTopicBodies guards the "point at the slug,
// don't duplicate the body" rule: prose that only exists in a judgment
// topic's full body (its good/bad examples) must not leak into the rendered
// skill body.
func TestContentMarkdownExcludesFullTopicBodies(t *testing.T) {
	md := Render().Markdown()

	topic, ok := help.Resolve("judgment:risk-vs-feature")
	if !ok {
		t.Fatal("judgment:risk-vs-feature not found in help; test fixture assumption broken")
	}
	const bodyOnlyLine = "we don't have a search command"
	if !strings.Contains(strings.ToLower(topic.Body), bodyOnlyLine) {
		t.Fatalf("test fixture assumption broken: %q not found in topic body", bodyOnlyLine)
	}
	if strings.Contains(strings.ToLower(md), bodyOnlyLine) {
		t.Error("rendered markdown duplicates full judgment topic body instead of pointing at the slug")
	}
}

// TestContentMarkdownExcludesDescription guards that Description (Claude-style
// frontmatter trigger text) is left for Adapters to place, not baked into the
// shared body every adapter emits verbatim.
func TestContentMarkdownExcludesDescription(t *testing.T) {
	c := Render()
	if c.Description == "" {
		t.Fatal("Description must not be empty")
	}
	if strings.Contains(c.Markdown(), c.Description) {
		t.Error("Markdown() must not embed Description; adapters place it themselves")
	}
}
