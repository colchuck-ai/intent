package gen

import (
	"bytes"
	"sort"
	"strings"
	"testing"

	"github.com/colchuck-ai/intent/internal/config"
	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
)

func buildSeed(t *testing.T) []File {
	t.Helper()
	r, err := model.Load("../model/testdata/seed.intent.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return Build(tree.Build(r), config.Default())
}

func find(files []File, path string) (File, bool) {
	for _, f := range files {
		if f.Path == path {
			return f, true
		}
	}
	return File{}, false
}

func TestBuildFileSet(t *testing.T) {
	files := buildSeed(t)
	var got []string
	for _, f := range files {
		got = append(got, f.Path)
	}
	want := []string{
		"change-records/add_resolvable_references_requirement.md",
		"engineering/README.md",
		"engineering/decision-records/single_go_binary_over_python.md",
		"engineering/validator.md",
		"product/README.md",
		"product/decision-records/prefer_explicit_declaration.md",
		"product/understand_the_rationale_behind_an_element/README.md",
		"product/understand_the_rationale_behind_an_element/fast_rationale_lookup.md",
	}
	if !sort.StringsAreSorted(got) {
		t.Errorf("files are not path-sorted: %v", got)
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("file set mismatch\n got: %v\nwant: %v", got, want)
	}
}

func TestBuildIsDeterministic(t *testing.T) {
	a, b := buildSeed(t), buildSeed(t)
	if len(a) != len(b) {
		t.Fatalf("file counts differ: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i].Path != b[i].Path || !bytes.Equal(a[i].Content, b[i].Content) {
			t.Errorf("build not byte-identical at %s", a[i].Path)
		}
	}
}

func TestEveryPageHasBanner(t *testing.T) {
	for _, f := range buildSeed(t) {
		if !bytes.HasPrefix(f.Content, []byte(Banner+"\n")) {
			t.Errorf("%s missing banner", f.Path)
		}
		if !bytes.HasSuffix(f.Content, []byte("\n")) || bytes.HasSuffix(f.Content, []byte("\n\n")) {
			t.Errorf("%s should end with exactly one newline", f.Path)
		}
	}
}

func TestInlineTargetLinksUseAnchors(t *testing.T) {
	// A requirement mitigates an inline risk on the same page -> bare fragment.
	f, ok := find(buildSeed(t), "product/understand_the_rationale_behind_an_element/fast_rationale_lookup.md")
	if !ok {
		t.Fatal("outcome page missing")
	}
	s := string(f.Content)
	if !strings.Contains(s, "[Ambiguous Reference](#ambiguous-reference)") {
		t.Errorf("same-page inline link should be a bare anchor:\n%s", s)
	}
	// A cross-page link into an inline element carries file + anchor.
	eng, _ := find(buildSeed(t), "engineering/validator.md")
	if !strings.Contains(string(eng.Content),
		"fast_rationale_lookup.md#resolvable-references") {
		t.Errorf("cross-page inline link should carry file#anchor:\n%s", eng.Content)
	}
}

func TestDerivedFooterBacklinksAndContents(t *testing.T) {
	files := buildSeed(t)

	// A leaf page's backlinks are a flat list under a single heading.
	val, _ := find(files, "engineering/validator.md")
	if !strings.Contains(string(val.Content), "## Referenced by") ||
		!strings.Contains(string(val.Content), "(affects)") {
		t.Errorf("validator should list its incoming affects backlink:\n%s", val.Content)
	}

	// A multi-element page groups backlinks under each target's name.
	oc, _ := find(files, "product/understand_the_rationale_behind_an_element/fast_rationale_lookup.md")
	body := string(oc.Content)
	for _, want := range []string{"**Ambiguous Reference**", "**Stable Logical IDs**", "**Resolvable References**"} {
		if !strings.Contains(body, want) {
			t.Errorf("grouped backlinks missing %q:\n%s", want, body)
		}
	}

	// A directory page lists its children, address-sorted.
	prod, _ := find(files, "product/README.md")
	pc := string(prod.Content)
	drAt := strings.Index(pc, "decision-records/prefer_explicit_declaration.md")
	jobAt := strings.Index(pc, "understand_the_rationale_behind_an_element/README.md")
	if drAt < 0 || jobAt < 0 || !strings.Contains(pc, "## Contents") {
		t.Fatalf("product index missing contents map:\n%s", pc)
	}
	if drAt > jobAt {
		t.Errorf("contents should be address-sorted (decision-records before understand_*):\n%s", pc)
	}
}

func TestSeeAlsoLinksSharedTargetSiblings(t *testing.T) {
	files := buildSeed(t)

	// Generator and Resolvable References both point at Stable Logical IDs, so
	// each is the other's shared-target sibling — cross-linked under "See also",
	// and in footer order after "Referenced by".
	oc, _ := find(files, "product/understand_the_rationale_behind_an_element/fast_rationale_lookup.md")
	body := string(oc.Content)
	seeAt := strings.Index(body, "## See also")
	refAt := strings.Index(body, "## Referenced by")
	if seeAt < 0 || refAt < 0 || seeAt < refAt {
		t.Fatalf("See also should follow Referenced by:\n%s", body)
	}
	if !strings.Contains(body[seeAt:], "**Resolvable References**") ||
		!strings.Contains(body[seeAt:], "engineering/README.md#generator") {
		t.Errorf("outcome See also should link Resolvable References -> Generator:\n%s", body[seeAt:])
	}

	// The relation is symmetric: Generator's page points back.
	eng, _ := find(files, "engineering/README.md")
	es := string(eng.Content)
	esSee := strings.Index(es, "## See also")
	if esSee < 0 || !strings.Contains(es[esSee:], "fast_rationale_lookup.md#resolvable-references") {
		t.Errorf("engineering See also should link Generator -> Resolvable References:\n%s", es)
	}

	// A directly-connected neighbor is not repeated as a see-also: Stable Logical
	// IDs' only sibling (Resolvable References) already references it, so it gets
	// no See also entry of its own.
	if strings.Contains(body[seeAt:], "**Stable Logical IDs**") {
		t.Errorf("directly-connected neighbor should be excluded from See also:\n%s", body[seeAt:])
	}
}

func TestKindLabelDisambiguatesSiblingTypes(t *testing.T) {
	files := buildSeed(t)

	// The engineering root flattens principles, constraints, and components as
	// sibling headings with no other visible distinction -- each must carry an
	// explicit kind label.
	eng, ok := find(files, "engineering/README.md")
	if !ok {
		t.Fatal("engineering README missing")
	}
	es := string(eng.Content)
	for _, want := range []string{
		"## Derive Only Mechanical Duals\n\n**Kind:** Principle",
		"## Single Go Binary\n\n**Kind:** Constraint",
		"## Generator\n\n**Kind:** Component",
	} {
		if !strings.Contains(es, want) {
			t.Errorf("engineering README missing kind label %q:\n%s", want, es)
		}
	}

	// A promoted component still carries its kind label even though it's alone
	// on its page -- a reader can land on it directly via a link.
	val, ok := find(files, "engineering/validator.md")
	if !ok {
		t.Fatal("validator page missing")
	}
	if !strings.Contains(string(val.Content), "# Validator\n\n**Kind:** Component") {
		t.Errorf("validator page missing kind label:\n%s", val.Content)
	}

	// An outcome flattens its risks and requirements as sibling headings --
	// each must carry an explicit kind label.
	oc, ok := find(files, "product/understand_the_rationale_behind_an_element/fast_rationale_lookup.md")
	if !ok {
		t.Fatal("outcome page missing")
	}
	body := string(oc.Content)
	for _, want := range []string{
		"## Ambiguous Reference\n\n**Kind:** Risk",
		"## Stable Logical IDs\n\n**Kind:** Requirement",
		"## Resolvable References\n\n**Kind:** Requirement",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("outcome page missing kind label %q:\n%s", want, body)
		}
	}

	// The outcome itself, and the job hosting it, are unambiguous from
	// structure alone (a job's only children are outcomes) and get no label.
	if strings.Contains(body, "**Kind:** Outcome") || strings.Contains(body, "**Kind:** Job") {
		t.Errorf("outcome and job should not carry a kind label:\n%s", body)
	}
}

func TestPathsAndRelativeInterpolation(t *testing.T) {
	src := `
product:
  name: P
  summary: s
engineering:
  name: E
  summary: s
  components:
    alpha:
      name: Alpha
      responsibility: See ![d]({{paths.assets}}/d.png) and {{.beta.link}} and {{.beta.name}}.
      fulfills: []
    beta:
      name: Beta
      responsibility: r
      fulfills: []
`
	r, err := model.Parse([]byte(src))
	if err != nil {
		t.Fatal(err)
	}
	files := Build(tree.Build(r), config.Default())
	// No document children, so the engineering root is a single leaf file.
	eng, ok := find(files, "engineering.md")
	if !ok {
		t.Fatal("engineering page missing")
	}
	s := string(eng.Content)
	if !strings.Contains(s, "![d](docs/assets/d.png)") {
		t.Errorf("paths.assets should substitute literally:\n%s", s)
	}
	if !strings.Contains(s, "[Beta](#beta)") {
		t.Errorf(".link should render a same-page anchor link:\n%s", s)
	}
	if !strings.Contains(s, "and Beta.") {
		t.Errorf(".name should render the bare name:\n%s", s)
	}
}
