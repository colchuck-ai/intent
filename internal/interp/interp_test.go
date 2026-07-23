package interp

import "testing"

func addrsOf(refs []Ref) []string {
	out := make([]string, len(refs))
	for i, r := range refs {
		out[i] = r.Addr
	}
	return out
}

func TestAbsoluteAddressStripsAccessor(t *testing.T) {
	prose := "See {{engineering.components.generator.link}} for the resolver."
	got := addrsOf(Addresses(prose, "engineering.components.validator"))
	want := []string{"engineering.components.generator"}
	if len(got) != 1 || got[0] != want[0] {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRelativeAddressResolvesAgainstContainer(t *testing.T) {
	// A component's prose referencing a sibling component by relative address.
	got := addrsOf(Addresses("both encode {{.reference_guides.link}}.", "engineering.components.bundle_linter"))
	want := "engineering.components.reference_guides"
	if len(got) != 1 || got[0] != want {
		t.Fatalf("got %v, want [%s]", got, want)
	}
}

func TestPathsExpressionSkipped(t *testing.T) {
	// paths.* is a config path variable, not an element address.
	got := Addresses("flow: ![x]({{paths.assets}}/linter-flow.png)", "engineering.components.validator")
	if len(got) != 0 {
		t.Fatalf("paths.* should be skipped, got %v", addrsOf(got))
	}
}

func TestBareAccessorlessAddress(t *testing.T) {
	// An address with no trailing accessor resolves to itself.
	got := addrsOf(Addresses("root is {{product}}.", "engineering"))
	if len(got) != 1 || got[0] != "product" {
		t.Fatalf("got %v, want [product]", got)
	}
}

func TestMultipleAndWhitespace(t *testing.T) {
	prose := "a {{ engineering.components.generator.name }} and b {{.stable_logical_ids.path}}"
	got := addrsOf(Addresses(prose, "product.jobs.j.outcomes.o.requirements.resolvable_references"))
	want := []string{
		"engineering.components.generator",
		"product.jobs.j.outcomes.o.requirements.stable_logical_ids",
	}
	if len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestRootOwnerRelative(t *testing.T) {
	// A relative address from a root element namespaces under the root itself.
	got := addrsOf(Addresses("{{.jobs.j.link}}", "product"))
	if len(got) != 1 || got[0] != "product.jobs.j" {
		t.Fatalf("got %v, want [product.jobs.j]", got)
	}
}

func TestExprPreserved(t *testing.T) {
	refs := Addresses("{{ .foo.link }}", "engineering.components.bar")
	if len(refs) != 1 || refs[0].Expr != ".foo.link" {
		t.Fatalf("expected the original inner expression preserved, got %+v", refs)
	}
}

func TestBareAccessorSkipped(t *testing.T) {
	// A malformed expression that is only an accessor yields no address.
	if got := Addresses("{{.link}}", "engineering.components.bar"); len(got) != 0 {
		t.Fatalf("expected no address for a bare accessor, got %v", addrsOf(got))
	}
}

func TestParseExprAccessorAndAddress(t *testing.T) {
	e := ParseExpr("engineering.components.generator.link", "engineering.components.validator")
	if e.Addr != "engineering.components.generator" || e.Accessor != "link" {
		t.Fatalf("got addr=%q acc=%q", e.Addr, e.Accessor)
	}
	if e.IsPath {
		t.Fatalf("element address parsed as path")
	}
}

func TestParseExprPathVar(t *testing.T) {
	e := ParseExpr("paths.assets", "engineering")
	if !e.IsPath || e.PathVar != "assets" {
		t.Fatalf("got IsPath=%v PathVar=%q", e.IsPath, e.PathVar)
	}
}

func TestParseExprBareAddressNoAccessor(t *testing.T) {
	e := ParseExpr("product", "engineering")
	if e.Addr != "product" || e.Accessor != "" {
		t.Fatalf("got addr=%q acc=%q", e.Addr, e.Accessor)
	}
}

func TestInterpolateRewritesEachExpression(t *testing.T) {
	prose := "see {{engineering.components.generator.link}} and {{paths.assets}}/x.png"
	got := Interpolate(prose, "engineering.components.validator", func(e Expr) string {
		if e.IsPath {
			return "ASSET"
		}
		return "<" + e.Addr + ":" + e.Accessor + ">"
	})
	want := "see <engineering.components.generator:link> and ASSET/x.png"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRewriteRemapsAbsoluteAddress(t *testing.T) {
	remap := func(a string) string {
		if a == "engineering.components.generator" {
			return "engineering.components.doc_generator"
		}
		return a
	}
	prose := "see {{engineering.components.generator.link}} and {{paths.assets}}/x.png"
	// Owner doesn't move here (same old/new owner).
	got := Rewrite(prose, "engineering.components.validator", "engineering.components.validator", remap)
	want := "see {{engineering.components.doc_generator.link}} and {{paths.assets}}/x.png"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// TestRewriteLeavesUnaffectedProseByteIdentical: an expression whose target
// doesn't move is returned untouched, preserving its original spacing (no churn).
func TestRewriteLeavesUnaffectedProseByteIdentical(t *testing.T) {
	remap := func(a string) string { return a } // nothing moves
	prose := "see {{ engineering.components.generator.link }} untouched"
	if got := Rewrite(prose, "engineering.components.validator", "engineering.components.validator", remap); got != prose {
		t.Fatalf("unaffected prose should be byte-identical:\ngot  %q\nwant %q", got, prose)
	}
}

// TestRewriteRelativeWithinMovedSubtreeStaysRelative: when both the owner and the
// target move by the same prefix (they live in the same moved subtree), a
// relative reference between them still resolves correctly and is left as-is.
func TestRewriteRelativeWithinMovedSubtreeStaysRelative(t *testing.T) {
	// A component's prose points at a sibling component relatively. Both move from
	// engineering.components.* to engineering.legacy.* (hypothetical prefix move).
	remap := func(a string) string {
		const old = "engineering.components"
		if a == old || len(a) > len(old) && a[:len(old)+1] == old+"." {
			return "engineering.legacy" + a[len(old):]
		}
		return a
	}
	prose := "resolver in {{.generator.link}}"
	oldOwner := "engineering.components.validator"
	newOwner := "engineering.legacy.validator"
	got := Rewrite(prose, oldOwner, newOwner, remap)
	if got != prose {
		t.Fatalf("relative ref within a moved subtree should be unchanged:\ngot  %q\nwant %q", got, prose)
	}
}

// TestRewriteRelativeToUnmovedTargetGoesAbsolute: an owner that moves out from
// under a target it referenced relatively must have that ref rewritten to
// absolute, or it would resolve against the wrong (new) container.
func TestRewriteRelativeToUnmovedTargetGoesAbsolute(t *testing.T) {
	// Only the owner moves; the relatively-referenced sibling stays put, so after
	// the move the relative form would dangle. Expect an absolute rewrite.
	remap := func(a string) string {
		if a == "product.jobs.j.outcomes.a" {
			return "product.jobs.k.outcomes.a"
		}
		return a
	}
	// owner outcome `a` moves from job j to job k; it referenced sibling `b`
	// relatively, but b stays under job j.
	prose := "depends on {{.b.link}}"
	oldOwner := "product.jobs.j.outcomes.a"
	newOwner := "product.jobs.k.outcomes.a"
	got := Rewrite(prose, oldOwner, newOwner, remap)
	want := "depends on {{product.jobs.j.outcomes.b.link}}"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
