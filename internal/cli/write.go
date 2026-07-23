package cli

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/schema"
	"github.com/colchuck-ai/intent/internal/tree"
	"github.com/colchuck-ai/intent/internal/validate"
)

// sortedKeys returns a map's keys sorted, joined with ", ". It keeps a
// user-facing list of accepted values in sync with the map that defines them —
// the same idiom tree.keyList uses for the --type filter — so adding a value in
// one place never leaves the help or error message stale.
func sortedKeys[V any](m map[string]V) string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return strings.Join(out, ", ")
}

// loadModel loads the raw model at --file for a write and builds a read index
// for address resolution. Writes resolve suffix arguments against this index,
// mutate the model, then re-validate the result before committing.
func loadModel(cmd *cobra.Command) (*model.Root, *tree.Index, string, error) {
	path, _ := cmd.Flags().GetString("file")
	r, err := model.Load(path)
	if err != nil {
		return nil, nil, "", err
	}
	return r, tree.Build(r), path, nil
}

// resolveAddr turns a user query (a full address or the shortest unambiguous
// suffix) into a full address via the read index (DESIGN §11).
func resolveAddr(ix *tree.Index, query string) (string, error) {
	e, err := ix.Resolve(query)
	if err != nil {
		return "", err
	}
	return e.Addr, nil
}

// resolveAll resolves a list of queries to full addresses, failing on the first
// that doesn't resolve.
func resolveAll(ix *tree.Index, queries []string) ([]string, error) {
	out := make([]string, 0, len(queries))
	for _, q := range queries {
		a, err := resolveAddr(ix, q)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}

// commitWrite is the single gate every write passes through: serialize to
// canonical YAML, re-check the schema, then run the linter, and only write the
// file if both pass (DESIGN §11 — validate auto-runs before every write). It
// returns the index built from the written tree so callers can compute the
// detail nudge without a re-parse. On any validation failure it prints the
// findings and writes nothing.
func commitWrite(cmd *cobra.Command, path string, r *model.Root) (*tree.Index, error) {
	b, err := r.Serialize()
	if err != nil {
		return nil, err
	}
	if err := schema.Validate(b); err != nil {
		return nil, fmt.Errorf("refusing to write — schema check failed: %w", err)
	}
	ix := tree.Build(r)
	if findings := validate.Check(ix); len(findings) > 0 {
		w := cmd.OutOrStdout()
		fmt.Fprintln(w, "refusing to write — validation failed:")
		for _, f := range findings {
			fmt.Fprintln(w, "  "+f.String())
		}
		return nil, fmt.Errorf("%d validation problem(s)", len(findings))
	}
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return nil, err
	}
	return ix, nil
}

// finishWrite commits the mutated model, prints the outcome line, the detail
// nudge for the affected element, and the scaffolded footer. addr is the element
// the write centered on (used for the nudge and the "next" suggestion).
func finishWrite(cmd *cobra.Command, path string, r *model.Root, verb, addr string, f footer) error {
	ix, err := commitWrite(cmd, path, r)
	if err != nil {
		return err
	}
	w := cmd.OutOrStdout()
	fmt.Fprintf(w, "%s %s\n", verb, addr)
	nudgeDetail(w, ix, addr)
	printFooter(w, f)
	return nil
}
