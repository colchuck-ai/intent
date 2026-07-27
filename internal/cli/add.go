package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
	"github.com/colchuck-ai/intent/internal/tree"
)

// addableKinds are the element types `add` creates. Records (pdr/adr/cr) are
// created with `record` (Phase 5), so they are excluded here.
var addableKinds = map[string]tree.Kind{
	"job":         tree.KindJob,
	"outcome":     tree.KindOutcome,
	"risk":        tree.KindRisk,
	"requirement": tree.KindRequirement,
	"component":   tree.KindComponent,
	"principle":   tree.KindPrinciple,
	"constraint":  tree.KindConstraint,
}

// promotable are the kinds that may carry a --type (inline | document).
var promotable = map[tree.Kind]bool{
	tree.KindOutcome: true, tree.KindRequirement: true, tree.KindComponent: true,
}

func newAddCmd() *cobra.Command {
	var f mutate.Fields
	var mitigates, dependsOn, fulfills, acceptance []string
	types := sortedKeys(addableKinds)

	cmd := &cobra.Command{
		Use:   "add <type> <parent> <key>",
		Short: "Create a new element under a parent",
		Long: "add creates one element of <type> under <parent> with key <key>.\n" +
			"Types: " + types + "\n" +
			"(records are created with `intent record`). Required edges must be given\n" +
			"here — a requirement needs --mitigates, a component needs --fulfills — because\n" +
			"the tree is validated before it is written and can't hold an element that\n" +
			"lacks them. The write is refused if it would make the tree invalid.",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			kind, ok := addableKinds[args[0]]
			if !ok {
				return fmt.Errorf("cannot add %q; types: %s (use `intent record` for pdr/adr/cr)", args[0], types)
			}
			// --type value validation lives in mutate.Add; the CLI only guards
			// that the kind can carry a type at all (jobs/risks/records can't).
			if f.Type != "" && !promotable[kind] {
				return fmt.Errorf("--type applies only to outcome, requirement, and component, not %s", kind)
			}

			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			parent, err := resolveAddr(ix, args[1])
			if err != nil {
				return err
			}
			if f.Mitigates, err = resolveAll(ix, mitigates); err != nil {
				return err
			}
			if f.DependsOn, err = resolveAll(ix, dependsOn); err != nil {
				return err
			}
			if f.Fulfills, err = resolveAll(ix, fulfills); err != nil {
				return err
			}
			f.Acceptance = acceptance

			addr, err := mutate.Add(r, kind, parent, args[2], f)
			if err != nil {
				return err
			}
			return finishWrite(cmd, path, r, "added", addr, addFooter(kind, addr))
		},
	}

	cmd.Flags().StringVar(&f.Name, "name", "", "display name (heading text)")
	cmd.Flags().StringVar(&f.Story, "story", "", "job story (the JTBD narrative)")
	cmd.Flags().StringVar(&f.Statement, "statement", "", "statement text (outcome/risk/requirement, or the principle/constraint value)")
	cmd.Flags().StringVar(&f.Responsibility, "responsibility", "", "component responsibility (one-line charter)")
	cmd.Flags().StringVar(&f.Detail, "detail", "", "freeform detail markdown")
	cmd.Flags().StringVar(&f.Type, "type", "", "inline | document (outcome/requirement/component only)")
	cmd.Flags().StringArrayVar(&mitigates, "mitigates", nil, "risk address a requirement mitigates (repeatable, required for a requirement)")
	cmd.Flags().StringArrayVar(&dependsOn, "depends-on", nil, "requirement address a requirement depends on (repeatable)")
	cmd.Flags().StringArrayVar(&fulfills, "fulfills", nil, "requirement address a component fulfills (repeatable, required for a component)")
	cmd.Flags().StringArrayVar(&acceptance, "acceptance", nil, "an acceptance criterion (repeatable)")
	return cmd
}
