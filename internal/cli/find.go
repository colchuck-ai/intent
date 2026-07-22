package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/tree"
)

func newFindCmd() *cobra.Command {
	var typ, domain string
	cmd := &cobra.Command{
		Use:   "find [query]",
		Short: "Search elements by name/prose/detail; filter by --type/--domain",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind, err := tree.ParseKind(typ)
			if err != nil {
				return err
			}
			dom, err := tree.ParseDomain(domain)
			if err != nil {
				return err
			}
			ix, err := loadIndex(cmd)
			if err != nil {
				return err
			}
			var query string
			if len(args) == 1 {
				query = args[0]
			}
			matches := ix.Find(query, kind, dom)

			w := cmd.OutOrStdout()
			if len(matches) == 0 {
				fmt.Fprintln(w, "no matches")
				return nil
			}
			for _, e := range matches {
				fmt.Fprintf(w, "%s  %s\n", e.Addr, label(e))
				if snip := matchSnippet(e, query, 100); snip != "" {
					fmt.Fprintf(w, "    %s\n", snip)
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&typ, "type", "", "filter by element type (job/outcome/risk/requirement/component/principle/constraint/pdr/adr/cr)")
	cmd.Flags().StringVar(&domain, "domain", "", "filter by domain (product/engineering/change)")
	return cmd
}
