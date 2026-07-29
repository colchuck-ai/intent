package help

import (
	"strings"
	"testing"
)

// TestEveryTopicIsWellFormed guards the embedded content: each topic parses,
// carries index metadata and a body, and its slug matches its plane.
func TestEveryTopicIsWellFormed(t *testing.T) {
	all := All()
	if len(all) == 0 {
		t.Fatal("no topics loaded from embedded content")
	}
	for _, top := range all {
		if top.Title == "" || top.Summary == "" || top.Body == "" {
			t.Errorf("%s: title/summary/body must all be non-empty", top.Slug)
		}
		if !strings.HasPrefix(top.Slug, top.Plane+":") {
			t.Errorf("%s: slug does not match plane %q", top.Slug, top.Plane)
		}
	}
}

// TestAllIsStablyOrdered guards that index/--json output is deterministic:
// planes appear in fixed display order, slugs sorted within.
func TestAllIsStablyOrdered(t *testing.T) {
	all := All()
	rankOf := map[string]int{}
	for i, p := range PlaneOrder() {
		rankOf[p] = i
	}
	for i := 1; i < len(all); i++ {
		prev, cur := all[i-1], all[i]
		rp, rc := rankOf[prev.Plane], rankOf[cur.Plane]
		if rp > rc {
			t.Errorf("plane order violated: %s (%s) before %s (%s)", prev.Slug, prev.Plane, cur.Slug, cur.Plane)
		}
		if rp == rc && prev.Slug > cur.Slug {
			t.Errorf("slug order violated within plane %s: %s before %s", prev.Plane, prev.Slug, cur.Slug)
		}
	}
}

func TestResolveExactSlug(t *testing.T) {
	if _, ok := Resolve("judgment:coherence"); !ok {
		t.Error("expected judgment:coherence to resolve")
	}
	if _, ok := Resolve("nope:missing"); ok {
		t.Error("a made-up slug must not resolve")
	}
}

// TestResolveErrorCode covers the inline-from-a-failure path: `intent help E001`
// (any case) reaches the errors:E001 entry.
func TestResolveErrorCode(t *testing.T) {
	for _, arg := range []string{"E001", "e001"} {
		top, ok := Resolve(arg)
		if !ok {
			t.Fatalf("expected %q to resolve", arg)
		}
		if top.Slug != "errors:E001" {
			t.Errorf("%q resolved to %q, want errors:E001", arg, top.Slug)
		}
	}
	if _, ok := Resolve("E999"); ok {
		t.Error("an unknown error code must not resolve")
	}
}

// TestResolveErrorsCatalog: the bare "errors" arg stitches the whole catalog
// into one page.
func TestResolveErrorsCatalog(t *testing.T) {
	top, ok := Resolve("errors")
	if !ok {
		t.Fatal("expected the errors catalog to resolve")
	}
	for _, code := range []string{"E001", "E002", "E003", "E004", "E005"} {
		if !strings.Contains(top.Body, code) {
			t.Errorf("errors catalog missing %s", code)
		}
	}
}
