// Package cli wires up the intent command tree.
//
// Phase 0 ships only the root command and --version. Read, validate, build, and
// write commands arrive in later phases.
package cli

import "github.com/spf13/cobra"

// version is stamped at build time via -ldflags; the default marks a local dev
// build.
var version = "0.0.0-dev"

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "intent",
		Short:   "Operate a schema'd intent.yaml and generate its docs",
		Long:    "intent reads, validates, edits, and renders an intent.yaml tree.\nThe YAML is canonical; generated docs are derived.",
		Version: version,
	}
}

// Execute runs the CLI and returns a process exit code.
func Execute() int {
	if err := newRootCmd().Execute(); err != nil {
		return 1
	}
	return 0
}
