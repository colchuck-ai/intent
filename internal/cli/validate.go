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

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Check the tree for dangling references, bad containment, and out-of-domain records",
		Long: "validate runs the structural schema check, then the semantic linter\n" +
			"(E001 dangling reference, E002 containment, E003 duplicate reference,\n" +
			"E004 record out of domain). It exits nonzero when anything fails.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("file")
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Shape first: the semantic rules assume a well-formed tree.
			if err := schema.Validate(b); err != nil {
				return fmt.Errorf("schema: %w", err)
			}
			r, err := model.Parse(b)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}

			findings := validate.Check(tree.Build(r))
			w := cmd.OutOrStdout()
			if len(findings) == 0 {
				fmt.Fprintln(w, "ok — no problems found")
				return nil
			}
			for _, f := range findings {
				fmt.Fprintln(w, f.String())
			}
			return fmt.Errorf("%d problem(s) found", len(findings))
		},
	}
}
