package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newAffectsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "affects <addr>",
		Short: "Show what references an element (its impact set)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ix, err := loadIndex(cmd)
			if err != nil {
				return err
			}
			e, err := ix.Resolve(args[0])
			if err != nil {
				return err
			}

			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "%s  %s\n", e.Addr, label(e))

			refs := ix.Incoming(e.Addr)
			if len(refs) == 0 {
				fmt.Fprintln(w, "  (nothing references this)")
				return nil
			}
			fmt.Fprintln(w, "referenced by:")
			for _, r := range refs {
				fmt.Fprintln(w, refLine(ix, "  ", "←", r.Kind, r.From, r.Note))
			}
			return nil
		},
	}
}
