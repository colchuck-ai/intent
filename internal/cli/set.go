package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
)

func newSetCmd() *cobra.Command {
	var cascade bool
	cmd := &cobra.Command{
		Use:   "set <addr> <field> <value>",
		Short: "Set a scalar field on an element",
		Long: "set writes one scalar field (name, statement, story, responsibility,\n" +
			"summary, detail, ...). Edges are managed with `intent link`/`unlink`, not\n" +
			"here. Setting type (inline|document) keeps containment valid by\n" +
			"construction: `set <addr> type document` promotes ancestors as needed, and\n" +
			"`set <addr> type inline` is refused if the element has document descendants\n" +
			"unless --cascade demotes the subtree. The write is refused if it would make\n" +
			"the tree invalid.",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			addr, err := resolveAddr(ix, args[0])
			if err != nil {
				return err
			}
			// Type changes route through the cascade-aware path so a document never
			// ends up under an inline parent (DESIGN §5).
			if args[1] == "type" {
				changed, err := mutate.SetType(r, ix, addr, args[2], cascade)
				if err != nil {
					return err
				}
				if len(changed) == 0 {
					fmt.Fprintf(cmd.OutOrStdout(), "%s is already %s; nothing changed\n", addr, args[2])
					return nil
				}
				reportCascade(cmd, addr, changed)
				return finishWrite(cmd, path, r, "set type "+args[2]+" on", addr, editFooter(addr))
			}
			if err := mutate.Set(r, addr, args[1], args[2]); err != nil {
				return err
			}
			return finishWrite(cmd, path, r, "set", addr, editFooter(addr))
		},
	}
	cmd.Flags().BoolVar(&cascade, "cascade", false, "when demoting to inline, demote the whole subtree")
	return cmd
}

// reportCascade prints the ancestors/descendants a type change swept along, so
// the containment cascade is never silent.
func reportCascade(cmd *cobra.Command, addr string, changed []string) {
	w := cmd.OutOrStdout()
	for _, a := range changed {
		if a != addr {
			fmt.Fprintf(w, "cascade: %s\n", a)
		}
	}
}
