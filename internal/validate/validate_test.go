package validate_test

import (
	"strings"
	"testing"

	"github.com/colchuck-ai/intent/internal/config"
	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
	"github.com/colchuck-ai/intent/internal/validate"
)

const seedPath = "../model/testdata/seed.intent.yaml"

// check loads path and runs the linter against it. Every fixture under
// testdata/ (and the seed) is snake_case, so E001-E004 tests check under
// config.Snake to keep their single expected finding isolated from E005.
func check(t *testing.T, path string) []validate.Finding {
	t.Helper()
	return checkAs(t, path, config.Snake)
}

func checkAs(t *testing.T, path string, casing config.Casing) []validate.Finding {
	t.Helper()
	r, err := model.Load(path)
	if err != nil {
		t.Fatalf("loading %s: %v", path, err)
	}
	return validate.Check(tree.Build(r), casing)
}

// TestSeedIsClean is the phase exit criterion: the canonical seed passes.
func TestSeedIsClean(t *testing.T) {
	if fs := check(t, seedPath); len(fs) != 0 {
		t.Fatalf("seed should be clean, got: %v", fs)
	}
}

// codeAt asserts exactly one finding, with the given code, at the given address.
func codeAt(t *testing.T, fs []validate.Finding, code validate.Code, addrSuffix string) {
	t.Helper()
	if len(fs) != 1 {
		t.Fatalf("expected one finding, got %d: %v", len(fs), fs)
	}
	if fs[0].Code != code {
		t.Errorf("got code %s, want %s: %v", fs[0].Code, code, fs[0])
	}
	if !strings.HasSuffix(fs[0].Addr, addrSuffix) {
		t.Errorf("finding at %q, want suffix %q", fs[0].Addr, addrSuffix)
	}
}

func TestE001DanglingEdge(t *testing.T) {
	fs := check(t, "testdata/e001_edge.intent.yaml")
	codeAt(t, fs, validate.E001, "requirements.resolvable_references")
	if !strings.Contains(fs[0].Detail, "dependsOn") || !strings.Contains(fs[0].Detail, "does_not_exist") {
		t.Errorf("detail should name the dangling edge: %s", fs[0].Detail)
	}
}

func TestE001DanglingProse(t *testing.T) {
	fs := check(t, "testdata/e001_prose.intent.yaml")
	codeAt(t, fs, validate.E001, "components.validator")
	if !strings.Contains(fs[0].Detail, "ghost") {
		t.Errorf("detail should name the dangling prose target: %s", fs[0].Detail)
	}
}

func TestE002Containment(t *testing.T) {
	fs := check(t, "testdata/e002.intent.yaml")
	codeAt(t, fs, validate.E002, "requirements.resolvable_references")
}

func TestE003Duplicate(t *testing.T) {
	fs := check(t, "testdata/e003.intent.yaml")
	codeAt(t, fs, validate.E003, "requirements.resolvable_references")
	if !strings.Contains(fs[0].Detail, "mitigates") {
		t.Errorf("detail should name the duplicated edge kind: %s", fs[0].Detail)
	}
}

func TestE004DomainScope(t *testing.T) {
	fs := check(t, "testdata/e004.intent.yaml")
	codeAt(t, fs, validate.E004, "decision_records.prefer_explicit_declaration")
	if !strings.Contains(fs[0].Detail, "product") {
		t.Errorf("detail should name the required domain: %s", fs[0].Detail)
	}
}

func TestE005KeyCase(t *testing.T) {
	fs := checkAs(t, "testdata/e005.intent.yaml", config.Snake)
	codeAt(t, fs, validate.E005, "outcomes.some-outcome")
	if !strings.Contains(fs[0].Detail, "snake_case") {
		t.Errorf("detail should name the configured convention: %s", fs[0].Detail)
	}

	// The same tree is clean under kebab-case, and the sole snake_case key
	// (some_job) trips it instead — the rule enforces whichever single
	// convention is configured, not "kebab" or "snake" specifically.
	fs = checkAs(t, "testdata/e005.intent.yaml", config.Kebab)
	codeAt(t, fs, validate.E005, "jobs.some_job")
}

// TestFindingStringHasHelpPointer guards the DESIGN §9 rule that every failure
// message ends with the help pointer.
func TestFindingStringHasHelpPointer(t *testing.T) {
	f := validate.Finding{Code: validate.E001, Addr: "product", Detail: "x"}
	if !strings.HasSuffix(f.String(), "→ intent help E001") {
		t.Errorf("finding line should end with the help pointer: %s", f.String())
	}
}
