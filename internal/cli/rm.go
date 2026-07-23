package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
	"github.com/colchuck-ai/intent/internal/tree"
)

func newRmCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rm <addr>",
		Short: "Remove an element (refused if it would dangle a reference)",
		Long: "rm deletes the element at <addr> and its subtree. It is refused when any\n" +
			"element outside that subtree still declares an edge into it, so a removal\n" +
			"never leaves a dangling reference behind.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			addr, err := resolveAddr(ix, args[0])
			if err != nil {
				return err
			}

			if blockers := incomingFromOutside(ix, addr); len(blockers) > 0 {
				w := cmd.OutOrStdout()
				fmt.Fprintf(w, "refusing to remove %s — it is still referenced by:\n", addr)
				for _, b := range blockers {
					fmt.Fprintln(w, "  "+b)
				}
				return fmt.Errorf("remove would dangle %d reference(s)", len(blockers))
			}

			if err := mutate.Remove(r, addr); err != nil {
				return err
			}
			// The element is gone, so there is no detail to nudge on; commit
			// directly rather than through finishWrite.
			if _, err := commitWrite(cmd, path, r); err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			fmt.Fprintf(w, "removed %s\n", addr)
			printFooter(w, footer{next: "intent tree", judgment: "judgment:coherence"})
			return nil
		},
	}
}

// incomingFromOutside returns the declared edges that point into the subtree
// rooted at addr from an element outside it — the references a removal would
// dangle. Prose {{ }} references are left to the post-write validate backstop
// (E001); this targets the common declared-edge case with a precise message.
func incomingFromOutside(ix *tree.Index, addr string) []string {
	removed := map[string]bool{addr: true}
	for _, e := range ix.All() {
		if strings.HasPrefix(e.Addr, addr+".") {
			removed[e.Addr] = true
		}
	}

	var blockers []string
	for target := range removed {
		for _, ref := range ix.Incoming(target) {
			if removed[ref.From] {
				continue
			}
			blockers = append(blockers, fmt.Sprintf("%s (%s → %s)", ref.From, ref.Kind, target))
		}
	}
	sort.Strings(blockers)
	return blockers
}
