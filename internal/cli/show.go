package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/tree"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <addr>",
		Short: "Show one element and its outgoing edges",
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
			renderShow(cmd.OutOrStdout(), ix, e)
			return nil
		},
	}
}

func renderShow(w io.Writer, ix *tree.Index, e *tree.Element) {
	fmt.Fprintln(w, e.Addr) // always echo the full address (DESIGN §11)
	fmt.Fprintln(w, label(e))

	if e.Prose != "" {
		field := e.ProseField
		if field == "" {
			field = "prose"
		}
		fmt.Fprintf(w, "\n%s: %s\n", field, e.Prose)
	}

	for _, f := range e.Fields {
		if f.List != nil {
			fmt.Fprintf(w, "\n%s:\n", f.Label)
			for _, item := range f.List {
				fmt.Fprintf(w, "  - %s\n", item)
			}
		} else if f.Text != "" {
			fmt.Fprintf(w, "\n%s: %s\n", f.Label, f.Text)
		}
	}

	if e.Detail != "" {
		fmt.Fprintf(w, "\ndetail: %s\n", e.Detail)
	}

	if len(e.Edges) > 0 {
		fmt.Fprintln(w, "\nedges:")
		for _, edge := range e.Edges {
			fmt.Fprintln(w, refLine(ix, "  ", "→", edge.Kind, edge.To, edge.Note))
		}
	}
}
