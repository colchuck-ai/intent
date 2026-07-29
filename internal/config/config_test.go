package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultsWhenFileMissing(t *testing.T) {
	c, err := LoadBeside(filepath.Join(t.TempDir(), "intent.yaml"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.OutputDir != "docs" {
		t.Errorf("default output_dir = %q, want docs", c.OutputDir)
	}
	if c.Paths["assets"] != "docs/assets" {
		t.Errorf("default assets = %q, want docs/assets", c.Paths["assets"])
	}
	if c.KeyCase != Kebab {
		t.Errorf("default key_case = %q, want %q", c.KeyCase, Kebab)
	}
}

func TestKeyCaseOverride(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, Filename), "key_case: snake_case\n")

	c, err := LoadBeside(filepath.Join(dir, "intent.yaml"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.KeyCase != Snake {
		t.Errorf("key_case = %q, want %q", c.KeyCase, Snake)
	}
}

func TestKeyCaseRejectsUnknownValue(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, Filename), "key_case: camelCase\n")

	if _, err := LoadBeside(filepath.Join(dir, "intent.yaml")); err == nil {
		t.Fatal("expected an unrecognized key_case value to be rejected")
	}
}

func TestOverridesAndAssetsFollowsOutputDir(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, Filename), "output_dir: site\npaths:\n  diagrams: design/diagrams\n")

	c, err := LoadBeside(filepath.Join(dir, "intent.yaml"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.OutputDir != "site" {
		t.Errorf("output_dir = %q, want site", c.OutputDir)
	}
	if c.Paths["diagrams"] != "design/diagrams" {
		t.Errorf("diagrams = %q, want design/diagrams", c.Paths["diagrams"])
	}
	// assets is defaulted relative to the overridden output_dir.
	if c.Paths["assets"] != "site/assets" {
		t.Errorf("assets = %q, want site/assets", c.Paths["assets"])
	}
}

func TestExplicitAssetsWins(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, Filename), "paths:\n  assets: static/img\n")

	c, err := LoadBeside(filepath.Join(dir, "intent.yaml"))
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if c.Paths["assets"] != "static/img" {
		t.Errorf("assets = %q, want static/img", c.Paths["assets"])
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
