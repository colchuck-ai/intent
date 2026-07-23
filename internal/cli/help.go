package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/help"
)

// planeTitles labels each plane in the task-first index.
var planeTitles = map[string]string{
	"guide":    "TASKS (guide:*) — goal → commands",
	"concept":  "CONCEPTS (concept:*) — what things are & why",
	"ref":      "REFERENCE (ref:*) — fields, edges, syntax, commands",
	"judgment": "JUDGMENT (judgment:*) — good, not just valid",
	"errors":   "ERRORS — cause + fix (also: intent help E0NN)",
}

func newHelpCmd() *cobra.Command {
	var (
		list   bool
		asJSON bool
	)
	cmd := &cobra.Command{
		Use:   "help [slug]",
		Short: "Read the embedded guidance: concepts, reference, judgment, guides, errors",
		Long: "help serves the instruction layer that ships inside the binary.\n" +
			"With no argument it prints a task-first index. Pass a slug\n" +
			"(e.g. concept:requirement, ref:edges, judgment:coherence, guide:edit)\n" +
			"or an error code (E001) to read one topic. --list is a flat, greppable\n" +
			"index; --json is the same index machine-readable.",
		Args: cobra.MaximumNArgs(1),
		// help never touches intent.yaml, so a missing file must not error.
		RunE: func(cmd *cobra.Command, args []string) error {
			w := cmd.OutOrStdout()
			switch {
			case asJSON:
				return writeHelpJSON(w)
			case len(args) == 1:
				return writeHelpTopic(w, args[0])
			case list:
				writeHelpList(w)
				return nil
			default:
				writeHelpIndex(w)
				return nil
			}
		},
	}
	cmd.Flags().BoolVar(&list, "list", false, "flat, greppable list of every topic")
	cmd.Flags().BoolVar(&asJSON, "json", false, "machine-readable index (stable order)")
	return cmd
}

// writeHelpTopic prints one resolved topic, or an actionable not-found error
// that points back at the index.
func writeHelpTopic(w io.Writer, arg string) error {
	t, ok := help.Resolve(arg)
	if !ok {
		return fmt.Errorf("no help topic %q; try `intent help --list`", arg)
	}
	fmt.Fprintln(w, t.Body)
	return nil
}

// writeHelpIndex prints the task-first index: planes in fixed order, each with
// its topics and one-line summaries.
func writeHelpIndex(w io.Writer) {
	fmt.Fprintln(w, "intent — declare product & engineering intent; the tool derives the docs.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Drive the CLI; intent.yaml is canonical; never hand-edit generated docs.")
	fmt.Fprintln(w, "Read a topic with `intent help <slug>` (e.g. `intent help guide:edit`).")
	for _, plane := range help.PlaneOrder() {
		topics := help.InPlane(plane)
		if len(topics) == 0 {
			continue
		}
		fmt.Fprintf(w, "\n%s\n", planeTitles[plane])
		for _, t := range topics {
			fmt.Fprintf(w, "  %-40s %s\n", t.Slug, t.Summary)
		}
	}
	fmt.Fprintln(w, "\n  intent help --list    every topic, one per line")
	fmt.Fprintln(w, "  intent help --json    machine-readable index")
}

// writeHelpList prints a flat, greppable "slug — summary" per line.
func writeHelpList(w io.Writer) {
	for _, t := range help.All() {
		fmt.Fprintf(w, "%s — %s\n", t.Slug, t.Summary)
	}
}

// writeHelpJSON prints the index as a stable JSON array (order from help.All).
func writeHelpJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(help.All())
}
