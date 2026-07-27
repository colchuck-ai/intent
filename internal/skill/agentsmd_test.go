package skill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAGENTSAdapterName(t *testing.T) {
	if got, want := (AGENTSAdapter{}).Name(), "agents-md"; got != want {
		t.Errorf("Name() = %q, want %q", got, want)
	}
}

func TestAGENTSAdapterFilesNoExistingFile(t *testing.T) {
	c := Render()
	dir := t.TempDir()

	files, err := (AGENTSAdapter{}).Files(dir, c)
	if err != nil {
		t.Fatalf("Files() error = %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Files() returned %d files, want 1", len(files))
	}

	f := files[0]
	if f.Path != "AGENTS.md" {
		t.Errorf("Path = %q, want %q", f.Path, "AGENTS.md")
	}
	body := string(f.Content)
	if !strings.Contains(body, agentsBeginMarker) || !strings.Contains(body, agentsEndMarker) {
		t.Error("content missing begin/end markers")
	}
	if !strings.Contains(body, c.Markdown()) {
		t.Error("content missing shared Content.Markdown() body")
	}
	if strings.Contains(body, "---\nname: intent") {
		t.Error("content must not carry Claude-style frontmatter")
	}
}

func TestAGENTSAdapterFilesMergesIntoExisting(t *testing.T) {
	c := Render()
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")
	original := "# Project Notes\n\nSome hand-written project context.\n"
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	files, err := (AGENTSAdapter{}).Files(dir, c)
	if err != nil {
		t.Fatalf("Files() error = %v", err)
	}
	body := string(files[0].Content)

	if !strings.Contains(body, "Some hand-written project context.") {
		t.Error("merge dropped pre-existing AGENTS.md content")
	}
	if !strings.Contains(body, c.Markdown()) {
		t.Error("merge missing shared Content.Markdown() body")
	}
	if strings.Index(body, "# Project Notes") > strings.Index(body, agentsBeginMarker) {
		t.Error("expected pre-existing content to remain before the appended generated section")
	}
}

func TestAGENTSAdapterFilesReplacesPriorGeneratedSection(t *testing.T) {
	c := Render()
	dir := t.TempDir()
	path := filepath.Join(dir, "AGENTS.md")

	first, err := (AGENTSAdapter{}).Files(dir, c)
	if err != nil {
		t.Fatalf("Files() (first) error = %v", err)
	}
	original := "# Project Notes\n\nHand-written context.\n\n" + string(first[0].Content)
	if err := os.WriteFile(path, []byte(original), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	second, err := (AGENTSAdapter{}).Files(dir, c)
	if err != nil {
		t.Fatalf("Files() (second) error = %v", err)
	}
	body := string(second[0].Content)

	if strings.Count(body, agentsBeginMarker) != 1 {
		t.Fatalf("reinstall must not duplicate the generated section; got body:\n%s", body)
	}
	if !strings.Contains(body, "Hand-written context.") {
		t.Error("reinstall dropped pre-existing hand-written content")
	}
}

func TestAGENTSAdapterFilesPathIgnoresRoot(t *testing.T) {
	c := Render()
	dir := t.TempDir()
	files, err := (AGENTSAdapter{}).Files(dir, c)
	if err != nil {
		t.Fatalf("Files() error = %v", err)
	}
	if got, want := files[0].Path, "AGENTS.md"; got != want {
		t.Errorf("Path = %q, want %q (must stay relative regardless of root)", got, want)
	}
}

func TestAGENTSAdapterFilesPropagatesReadError(t *testing.T) {
	c := Render()
	dir := t.TempDir()
	// AGENTS.md as a directory makes os.ReadFile fail with something other
	// than ErrNotExist.
	if err := os.Mkdir(filepath.Join(dir, "AGENTS.md"), 0o755); err != nil {
		t.Fatalf("Mkdir: %v", err)
	}

	if _, err := (AGENTSAdapter{}).Files(dir, c); err == nil {
		t.Error("Files() error = nil, want non-nil for unreadable AGENTS.md")
	}
}
