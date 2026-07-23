// Package cli wires up the intent command tree.
//
// Phase 1 adds the read commands (tree/show/find/trace/affects) over a shared
// element index. Validate, build, and write commands arrive in later phases.
package cli

import (
	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
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
	)
	return root
}

// loadIndex loads the tree from --file and builds the in-memory index. Every
// read command starts here.
func loadIndex(cmd *cobra.Command) (*tree.Index, error) {
	path, _ := cmd.Flags().GetString("file")
	r, err := model.Load(path)
	if err != nil {
		return nil, err
	}
	return tree.Build(r), nil
}

// Execute runs the CLI and returns a process exit code.
func Execute() int {
	if err := newRootCmd().Execute(); err != nil {
		return 1
	}
	return 0
}
