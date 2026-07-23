package cli

import (
	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
)

func newSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <addr> <field> <value>",
		Short: "Set a scalar field on an element",
		Long: "set writes one scalar field (name, statement, story, responsibility,\n" +
			"summary, detail, type, ...). Edges are managed with `intent link`/`unlink`,\n" +
			"not here. The write is refused if it would make the tree invalid.",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, ix, path, err := loadModel(cmd)
			if err != nil {
				return err
			}
			addr, err := resolveAddr(ix, args[0])
			if err != nil {
				return err
			}
			if err := mutate.Set(r, addr, args[1], args[2]); err != nil {
				return err
			}
			return finishWrite(cmd, path, r, "set", addr, editFooter(addr))
		},
	}
}
