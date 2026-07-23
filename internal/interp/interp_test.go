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
