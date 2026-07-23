package cli

import (
	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/mutate"
)

func newPromoteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promote <addr>",
		Short: "Promote an element to a document (cascading ancestors)",
		Long: "promote makes an outcome, requirement, or component render as its own\n" +
			"document. Every inline ancestor is promoted with it (like `mkdir -p`), so\n" +
			"a document never ends up under an inline parent. It is the same as\n" +
			"`intent set <addr> type document`; to go the other way use\n" +
			"`intent set <addr> type inline`.",
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
			changed, err := mutate.SetType(r, ix, addr, "document", false)
			if err != nil {
				return err
			}
			if len(changed) == 0 {
				cmd.Printf("%s is already a document; nothing changed\n", addr)
				return nil
			}
			reportCascade(cmd, addr, changed)
			return finishWrite(cmd, path, r, "promoted", addr, editFooter(addr))
		},
	}
}
