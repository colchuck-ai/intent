package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
	"github.com/colchuck-ai/intent/internal/tree"
)

// recordKinds maps the record subcommand word to its kind. Location implies the
// record type, so choosing the word chooses where it lands (DESIGN §2).
var recordKinds = map[string]tree.Kind{
	"pdr": tree.KindPDR,
	"adr": tree.KindADR,
	"cr":  tree.KindCR,
}

func newRecordCmd() *cobra.Command {
	var f mutate.RecordFields
	var affects []string
	kinds := sortedKeys(recordKinds)

	cmd := &cobra.Command{
		Use:   "record <pdr|adr|cr> <key>",
		Short: "Create a decision or change record",
		Long: "record creates a record of the chosen type and puts it where that type\n" +
			"lives: a pdr under product, an adr under engineering, a cr at top level.\n" +
			"Every record needs at least one --affects. A pdr/adr needs a --summary and\n" +
			"its affects must stay in its own domain; a cr needs --change and\n" +
			"--rationale and may affect any domain. The write is refused if it would\n" +
			"make the tree invalid.",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind, ok := recordKinds[args[0]]
			if !ok {
				return fmt.Errorf("cannot record %q; types: %s", args[0], kinds)
			}
			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			if f.Affects, err = resolveAll(ix, affects); err != nil {
				return err
			}

			addr, err := mutate.AddRecord(r, kind, args[1], f)
			if err != nil {
				return err
			}
			return finishWrite(cmd, path, r, "recorded", addr, recordFooter(addr))
		},
	}

	cmd.Flags().StringVar(&f.Name, "name", "", "display name (heading text)")
	cmd.Flags().StringSliceVar(&affects, "affects", nil, "an address (or root literal) the record affects (repeatable, required)")
	cmd.Flags().StringVar(&f.Summary, "summary", "", "decision-record core statement (pdr/adr)")
	cmd.Flags().StringVar(&f.Context, "context", "", "decision-record context (pdr/adr)")
	cmd.Flags().StringSliceVar(&f.Options, "option", nil, "a decision-record option (repeatable, pdr/adr)")
	cmd.Flags().StringVar(&f.Decision, "decision", "", "the decision reached (pdr/adr)")
	cmd.Flags().StringSliceVar(&f.Consequences, "consequence", nil, "a decision-record consequence (repeatable, pdr/adr)")
	cmd.Flags().StringVar(&f.Change, "change", "", "what changed (cr)")
	cmd.Flags().StringVar(&f.Rationale, "rationale", "", "why it changed (cr)")
	return cmd
}
