package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
	"github.com/colchuck-ai/intent/internal/tree"
)

// edgeKinds is the closed set of edge kinds link/unlink accept.
var edgeKinds = map[string]tree.EdgeKind{
	"mitigates":     tree.EdgeMitigates,
	"dependsOn":     tree.EdgeDependsOn,
	"fulfills":      tree.EdgeFulfills,
	"affects":       tree.EdgeAffects,
	"relationships": tree.EdgeRelationships,
}

func parseEdgeKind(s string) (tree.EdgeKind, error) {
	if k, ok := edgeKinds[s]; ok {
		return k, nil
	}
	return "", fmt.Errorf("unknown edge kind %q; valid kinds: %s", s, sortedKeys(edgeKinds))
}

func newLinkCmd() *cobra.Command {
	var note string
	cmd := &cobra.Command{
		Use:   "link <kind> <from> <to>",
		Short: "Declare an edge from one element to another",
		Long: "link writes the storage shape for a declared edge so no one hand-builds\n" +
			"an edge list. <kind> is mitigates, dependsOn, fulfills, affects, or\n" +
			"relationships; <to> may be a root literal (product/engineering) for affects.\n" +
			"Use --note to attach the note a relationships edge carries.",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind, err := parseEdgeKind(args[0])
			if err != nil {
				return err
			}
			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			from, err := resolveAddr(ix, args[1])
			if err != nil {
				return err
			}
			to, err := resolveAddr(ix, args[2])
			if err != nil {
				return err
			}
			if err := mutate.Link(r, kind, from, to, note); err != nil {
				return err
			}
			return finishWrite(cmd, path, r, "linked", from, editFooter(from))
		},
	}
	cmd.Flags().StringVar(&note, "note", "", "note for a relationships edge")
	return cmd
}

// resolveEdgeTarget resolves an unlink target. Normally it resolves through the
// index like any address. But a dangling edge — one whose target no longer
// exists, left by a hand-edit — can't resolve there, and unlink must still be
// able to remove it. So on failure it matches the query against the addresses
// `from` actually declares an edge of this kind to (by whole-segment suffix, the
// same input mode as everything else), and only falls back to the literal query
// if nothing matches — leaving mutate.Unlink to report "no such edge".
func resolveEdgeTarget(ix *tree.Index, from string, kind tree.EdgeKind, query string) string {
	if addr, err := resolveAddr(ix, query); err == nil {
		return addr
	}
	if e, ok := ix.Get(from); ok {
		for _, ed := range e.Edges {
			if ed.Kind == kind && (ed.To == query || strings.HasSuffix(ed.To, "."+query)) {
				return ed.To
			}
		}
	}
	return query
}

func newUnlinkCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unlink <kind> <from> <to>",
		Short: "Remove a declared edge",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind, err := parseEdgeKind(args[0])
			if err != nil {
				return err
			}
			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			from, err := resolveAddr(ix, args[1])
			if err != nil {
				return err
			}
			to := resolveEdgeTarget(ix, from, kind, args[2])
			if err := mutate.Unlink(r, kind, from, to); err != nil {
				return err
			}
			return finishWrite(cmd, path, r, "unlinked", from, editFooter(from))
		},
	}
}
