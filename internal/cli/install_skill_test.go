package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/colchuck-ai/intent/internal/skill"
)

func TestInstallSkillUnknownAgentListsSupported(t *testing.T) {
	out, err := run(t, "install-skill", "--agent", "bogus")
	if err == nil {
		t.Fatalf("expected an error for an unknown agent\n%s", out)
	}
	if !strings.Contains(err.Error(), "agents-md") || !strings.Contains(err.Error(), "claude-code") {
		t.Errorf("unknown-agent error should list supported agents: %v", err)
	}
}

func TestInstallSkillClaudeCodeEndToEnd(t *testing.T) {
	dir := t.TempDir()
	out, err := run(t, "install-skill", "--agent", "claude-code", "--dir", dir)
	if err != nil {
		t.Fatalf("install-skill: %v\n%s", err, out)
	}

	wantPath := filepath.Join(dir, ".claude", "skills", "intent", "SKILL.md")
	if !strings.Contains(out, wantPath) {
		t.Errorf("output should print the installed path %q, got:\n%s", wantPath, out)
	}

	b, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	c := skill.Render()
	files, err := (skill.ClaudeAdapter{}).Files(dir, c)
	if err != nil {
		t.Fatalf("Files() error = %v", err)
	}
	if string(b) != string(files[0].Content) {
		t.Error("installed SKILL.md content does not match the rendered adapter output")
	}
}

func TestInstallSkillAgentsMdEndToEnd(t *testing.T) {
	dir := t.TempDir()
	out, err := run(t, "install-skill", "--agent", "agents-md", "--dir", dir)
	if err != nil {
		t.Fatalf("install-skill: %v\n%s", err, out)
	}

	wantPath := filepath.Join(dir, "AGENTS.md")
	if !strings.Contains(out, wantPath) {
		t.Errorf("output should print the installed path %q, got:\n%s", wantPath, out)
	}

	b, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(b), skill.Render().Markdown()) {
		t.Error("installed AGENTS.md missing the shared skill body")
	}
}

func TestInstallSkillMissingAgentFlagErrors(t *testing.T) {
	_, err := run(t, "install-skill", "--dir", t.TempDir())
	if err == nil {
		t.Fatal("expected an error when --agent is not given")
	}
}
