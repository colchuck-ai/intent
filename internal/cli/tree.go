package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/tree"
)

func newTreeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tree",
		Short: "Print the whole tree as an indented outline",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ix, err := loadIndex(cmd)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			for _, e := range ix.All() {
				indent := strings.Repeat("  ", e.Level)
				fmt.Fprintf(w, "%s%s %s\n", indent, tree.LastSegment(e.Addr), tagged(e))
			}
			return nil
		},
	}
}

// tagged renders the "[kind] Name (document)" suffix for an outline line. The
// key is printed separately, so the name is omitted when absent
// (principles/constraints) rather than falling back to the key.
func tagged(e *tree.Element) string { return kindLabel(e, e.Name) }
