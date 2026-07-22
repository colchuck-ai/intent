package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/tree"
)

func newTraceCmd() *cobra.Command {
	var depth int
	var up, down bool
	var edges []string
	cmd := &cobra.Command{
		Use:   "trace <addr>",
		Short: "Walk an element's neighborhood across declared edges",
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

			// Default (neither flag) walks both directions.
			wantUp, wantDown := up, down
			if !up && !down {
				wantUp, wantDown = true, true
			}

			kinds := map[tree.EdgeKind]bool{}
			for _, k := range edges {
				kinds[tree.EdgeKind(k)] = true
			}

			nodes := ix.Trace(e.Addr, depth, wantUp, wantDown, kinds)
			renderTrace(cmd.OutOrStdout(), ix, nodes)
			return nil
		},
	}
	cmd.Flags().IntVar(&depth, "depth", 1, "how many hops to walk")
	cmd.Flags().BoolVar(&up, "up", false, "walk only upstream (incoming refs)")
	cmd.Flags().BoolVar(&down, "down", false, "walk only downstream (outgoing edges)")
	cmd.Flags().StringSliceVar(&edges, "edges", nil, "restrict to these edge kinds (comma-separated)")
	return cmd
}

// renderTrace prints the neighborhood as a nested outline. The traversal visits
// each node once and records the node it came from, so the nodes form a tree
// rooted at the start; we walk that parent→child structure (rather than banding
// purely by depth) so the graph shape stays legible.
func renderTrace(w io.Writer, ix *tree.Index, nodes []tree.TraceNode) {
	if len(nodes) == 0 {
		return
	}
	children := map[string][]tree.TraceNode{}
	for _, n := range nodes[1:] {
		children[n.From] = append(children[n.From], n)
	}

	root := nodes[0]
	fmt.Fprintf(w, "%s  %s\n", root.Addr, label(getOrStub(ix, root.Addr)))

	var walk func(addr string, depth int)
	walk = func(addr string, depth int) {
		for _, c := range children[addr] {
			arrow := "→" // downstream: we point at it
			if c.Up {
				arrow = "←" // upstream: it points at us
			}
			fmt.Fprintln(w, refLine(ix, strings.Repeat("  ", depth), arrow, c.Via, c.Addr, ""))
			walk(c.Addr, depth+1)
		}
	}
	walk(root.Addr, 1)
}

// getOrStub returns the element at addr, or a bare stub carrying just the address
// if it isn't in the index. The trace start is always resolved; a stub only
// arises for a dangling edge target (flagged as E001 in Phase 2).
func getOrStub(ix *tree.Index, addr string) *tree.Element {
	if e, ok := ix.Get(addr); ok {
		return e
	}
	return &tree.Element{Addr: addr}
}
