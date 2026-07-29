package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Check the tree for dangling references, bad containment, and out-of-domain records",
		Long: "validate runs the structural schema check, then the semantic linter\n" +
			"(E001 dangling reference, E002 containment, E003 duplicate reference,\n" +
			"E004 record out of domain, E005 key doesn't match key_case). It exits\n" +
			"nonzero when anything fails.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("file")
			_, findings, err := checkedIndex(path)
			if err != nil {
				return err
			}
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
