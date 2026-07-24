package skill

import (
	"strings"
	"testing"
)

func TestClaudeAdapterName(t *testing.T) {
	if got, want := (ClaudeAdapter{}).Name(), "claude-code"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestClaudeAdapterFiles(t *testing.T) {
	c := Render()
	files, err := (ClaudeAdapter{}).Files(".", c)
	if err != nil {
		t.Fatalf("Files() error = %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Files() returned %d files, want 1", len(files))
	}

	f := files[0]
	const wantPath = ".claude/skills/intent/SKILL.md"
	if f.Path != wantPath {
		t.Errorf("Path = %q, want %q", f.Path, wantPath)
	}

	body := string(f.Content)
	wantPrefix := "---\nname: intent\ndescription: " + c.Description + "\n---\n\n"
	if !strings.HasPrefix(body, wantPrefix) {
		t.Errorf("content does not start with the expected frontmatter:\n%s", body)
	}
	if !strings.HasSuffix(body, c.Markdown()) {
		t.Error("content does not end with the shared Content.Markdown() body")
	}
}

func TestClaudeAdapterFilesPathIgnoresRoot(t *testing.T) {
	// Path is always relative — callers join it with the real install root,
	// so Files must not fold root into the returned Path itself.
	c := Render()
	files, err := (ClaudeAdapter{}).Files("/some/other/root", c)
	if err != nil {
		t.Fatalf("Files() error = %v", err)
	}
	if got, want := files[0].Path, ".claude/skills/intent/SKILL.md"; got != want {
		t.Errorf("Path = %q, want %q (must stay relative regardless of root)", got, want)
	}
}
