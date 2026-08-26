package gen

import "testing"

func TestFilePathFor(t *testing.T) {
	cases := []struct {
		addr  string
		isDir bool
		want  string
	}{
		{"product", true, "product/README.md"},
		{"product.jobs.j", true, "product/j/README.md"},
		{"product.jobs.j.outcomes.o", false, "product/j/o.md"},
		{"product.decision_records.d", false, "product/decision-records/d.md"},
		{"engineering.components.c", false, "engineering/c.md"},
		{"engineering.decision_records.a", false, "engineering/decision-records/a.md"},
		{"change_records.cr", false, "change-records/cr.md"},
	}
	for _, c := range cases {
		if got := filePathFor(c.addr, c.isDir); got != c.want {
			t.Errorf("filePathFor(%q, %v) = %q, want %q", c.addr, c.isDir, got, c.want)
		}
	}
}

func TestSlashRel(t *testing.T) {
	cases := []struct{ from, to, want string }{
		{"engineering/README.md", "engineering/validator.md", "validator.md"},
		{"engineering/validator.md", "engineering/README.md", "README.md"},
		{"product/j/o.md", "product/decision-records/d.md", "../decision-records/d.md"},
		{"engineering/validator.md", "product/j/o.md", "../product/j/o.md"},
		{"product/README.md", "product/j/README.md", "j/README.md"},
	}
	for _, c := range cases {
		if got := slashRel(c.from, c.to); got != c.want {
			t.Errorf("slashRel(%q, %q) = %q, want %q", c.from, c.to, got, c.want)
		}
	}
}

func TestSlug(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Stable Logical IDs", "stable-logical-ids"},
		{"Ambiguous Reference", "ambiguous-reference"},
		{"single_go_binary", "single_go_binary"},
		{"Foo: Bar (baz)", "foo-bar-baz"},
	}
	for _, c := range cases {
		if got := slug(c.in); got != c.want {
			t.Errorf("slug(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}
