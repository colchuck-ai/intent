package skill

import "strings"

// ClaudeAdapter emits Claude Code's generated skill file: YAML frontmatter
// carrying the trigger Description, followed by the shared Markdown body.
type ClaudeAdapter struct{}

var _ Adapter = ClaudeAdapter{}

// Name identifies this adapter for the --agent claude-code flag value.
func (ClaudeAdapter) Name() string { return "claude-code" }

// Files returns the single .claude/skills/intent/SKILL.md file, with Path
// relative to root (installed skill files are generated, never committed —
// see .gitignore). Claude Code always gets exactly this one file, so unlike
// an adapter that merges into an existing shared file, it never reads from
// root.
func (ClaudeAdapter) Files(root string, c Content) ([]File, error) {
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("name: intent\n")
	b.WriteString("description: ")
	b.WriteString(c.Description)
	b.WriteString("\n---\n\n")
	b.WriteString(c.Markdown())

	return []File{{
		Path:    ".claude/skills/intent/SKILL.md",
		Content: []byte(b.String()),
	}}, nil
}
