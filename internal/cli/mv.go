package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
)

func newMvCmd() *cobra.Command {
	var opts mutate.MvOpts
	cmd := &cobra.Command{
		Use:   "mv <addr> [new-parent]",
		Short: "Reorder, rename, or re-parent an element",
		Long: "mv moves the element at <addr>. Give a [new-parent] to re-parent it,\n" +
			"--rename to change its key, and --before/--after a sibling to position it.\n" +
			"Only outcomes (between jobs) and risks/requirements (between outcomes) can\n" +
			"be re-parented; everything else has a fixed home. Every reference to the\n" +
			"element — declared edges, relationship keys, and prose {{ }} addresses — is\n" +
			"rewritten across the whole tree, so nothing is left dangling. The write is\n" +
			"refused if it would make the tree invalid.",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			addr, err := resolveAddr(ix, args[0])
			if err != nil {
				return err
			}
			if len(args) == 2 {
				if opts.NewParent, err = resolveAddr(ix, args[1]); err != nil {
					return err
				}
			}

			newAddr, refs, promoted, err := mutate.Mv(r, ix, addr, opts)
			if err != nil {
				return err
			}
			writtenIx, err := commitWrite(cmd, path, r)
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			switch {
			case newAddr == addr:
				fmt.Fprintf(w, "reordered %s\n", addr)
			case refs > 0:
				fmt.Fprintf(w, "moved %s → %s (rewrote %d reference(s))\n", addr, newAddr, refs)
			default:
				fmt.Fprintf(w, "moved %s → %s\n", addr, newAddr)
			}
			// A re-parent may promote the destination chain to keep containment valid;
			// report it so the type change is never silent.
			for _, a := range promoted {
				fmt.Fprintf(w, "promoted %s (to host the moved document)\n", a)
			}
			nudgeDetail(w, writtenIx, newAddr)
			printFooter(w, editFooter(newAddr))
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.Rename, "rename", "", "new snake_case key for the element")
	cmd.Flags().StringVar(&opts.Before, "before", "", "position the element before this sibling")
	cmd.Flags().StringVar(&opts.After, "after", "", "position the element after this sibling")
	return cmd
}
