package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/skill"
)

// skillAdapters are the agent targets --agent accepts, keyed by Adapter.Name().
var skillAdapters = map[string]skill.Adapter{
	"claude-code": skill.ClaudeAdapter{},
	"agents-md":   skill.AGENTSAdapter{},
}

func newInstallSkillCmd() *cobra.Command {
	var agent, dir string

	cmd := &cobra.Command{
		Use:   "install-skill --agent <target>",
		Short: "Render the agent skill from embedded help and install it",
		Long: "install-skill renders the shared skill content (DESIGN §12) — a trigger\n" +
			"description, an Intent paragraph, the entry-point instruction, gotchas, and\n" +
			"judgment one-liners — and writes it in the target agent's format under --dir.\n" +
			"Installed files are generated; reinstall after changing help content rather\n" +
			"than hand-editing them. Targets: " + sortedKeys(skillAdapters) + ".",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			a, ok := skillAdapters[agent]
			if !ok {
				return fmt.Errorf("unknown agent %q; supported agents: %s", agent, sortedKeys(skillAdapters))
			}

			files, err := a.Files(dir, skill.Render())
			if err != nil {
				return err
			}
			w := cmd.OutOrStdout()
			for _, f := range files {
				dest := filepath.Join(dir, filepath.FromSlash(f.Path))
				if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
					return err
				}
				if err := os.WriteFile(dest, f.Content, 0o644); err != nil {
					return err
				}
				fmt.Fprintln(w, dest)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&agent, "agent", "", "target agent: "+sortedKeys(skillAdapters)+" (required)")
	cmd.Flags().StringVar(&dir, "dir", ".", "install root")
	return cmd
}
