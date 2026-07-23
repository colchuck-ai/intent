package cli

import (
	"strings"
	"testing"
)

func TestValidatePassesCleanSeed(t *testing.T) {
	out, err := run(t, "validate", "-f", seedFlag)
	if err != nil {
		t.Fatalf("validate should pass the seed: %v\n%s", err, out)
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("expected an ok message:\n%s", out)
	}
}

func TestValidateReportsCodeWithHelpPointer(t *testing.T) {
	out, err := run(t, "validate", "-f", "../validate/testdata/e004.intent.yaml")
	if err == nil {
		t.Fatalf("expected a nonzero exit for an invalid tree\n%s", out)
	}
	if !strings.Contains(out, "E004") {
		t.Errorf("expected the E004 code in output:\n%s", out)
	}
	if !strings.Contains(out, "→ intent help E004") {
		t.Errorf("every failure must end with the help pointer:\n%s", out)
	}
}

func TestValidateReportsSchemaFailure(t *testing.T) {
	// The malformed fixture omits the engineering root, so the schema bites
	// before the linter runs.
	out, err := run(t, "validate", "-f", "testdata/malformed.intent.yaml")
	if err == nil {
		t.Fatalf("expected a schema failure\n%s", out)
	}
	if !strings.Contains(out, "schema") {
		t.Errorf("expected the error to be attributed to the schema:\n%s", out)
	}
}
