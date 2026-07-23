// Package cli wires up the intent command tree.
//
// Phase 1 adds the read commands (tree/show/find/trace/affects) over a shared
// element index. Validate, build, and write commands arrive in later phases.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/schema"
	"github.com/colchuck-ai/intent/internal/tree"
	"github.com/colchuck-ai/intent/internal/validate"
)

// version is stamped at build time via -ldflags; the default marks a local dev
// build.
var version = "0.0.0-dev"

// defaultFile is where every command looks for the tree unless --file overrides.
const defaultFile = "intent.yaml"

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:     "intent",
		Short:   "Operate a schema'd intent.yaml and generate its docs",
		Long:    "intent reads, validates, edits, and renders an intent.yaml tree.\nThe YAML is canonical; generated docs are derived.",
		Version: version,
		// Runtime errors (bad path, unknown address) shouldn't dump usage.
		SilenceUsage: true,
	}
	root.PersistentFlags().StringP("file", "f", defaultFile, "path to intent.yaml")
	root.AddCommand(
		newTreeCmd(),
		newShowCmd(),
		newFindCmd(),
		newTraceCmd(),
		newAffectsCmd(),
		newValidateCmd(),
		newBuildCmd(),
		newCheckCmd(),
		newAddCmd(),
		newSetCmd(),
		newRmCmd(),
		newLinkCmd(),
		newUnlinkCmd(),
		newPromoteCmd(),
		newMvCmd(),
		newRecordCmd(),
	)
	// `intent help` is our embedded topic system (DESIGN §12), not cobra's
	// subcommand-usage dumper. SetHelpCommand replaces the auto-generated one;
	// per-command usage is still reachable via `intent <cmd> --help`.
	root.SetHelpCommand(newHelpCmd())
	return root
}

// loadIndex loads the tree from --file and builds the in-memory index. Every
// read command starts here; it trusts the file's shape (no validation), because
// reads only surface what's there.
func loadIndex(cmd *cobra.Command) (*tree.Index, error) {
	path, _ := cmd.Flags().GetString("file")
	r, err := model.Load(path)
	if err != nil {
		return nil, err
	}
	return tree.Build(r), nil
}

// checkedIndex reads the tree at path, runs the structural schema check then the
// semantic linter, and returns the built index with any findings. It errors only
// for IO/shape failures that stop analysis; a well-formed tree with rule
// violations returns its findings and a nil error, leaving the caller (validate
// reports them; build refuses to render) to decide how to react. This is the one
// place the read → schema → parse → check pipeline lives.
func checkedIndex(path string) (*tree.Index, []validate.Finding, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}
	if err := schema.Validate(b); err != nil {
		return nil, nil, fmt.Errorf("schema: %w", err)
	}
	r, err := model.Parse(b)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	ix := tree.Build(r)
	return ix, validate.Check(ix), nil
}

// Execute runs the CLI and returns a process exit code.
func Execute() int {
	if err := newRootCmd().Execute(); err != nil {
		return 1
	}
	return 0
}
