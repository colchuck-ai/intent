// Package config loads the optional intent.config.yaml that tunes generation
// (DESIGN §7). Everything has a sane default, so a project with no config file
// still builds; the file only overrides.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Casing is the element-key convention a project requires, enforced by the
// linter (validate.E005) and the write path (mutate).
type Casing string

const (
	Kebab Casing = "kebab-case"
	Snake Casing = "snake_case"
)

// Config is the resolved generator configuration.
//
//   - OutputDir is where generated markdown is written (default "docs").
//   - Paths is the map of named path variables `{{paths.*}}` resolves against,
//     by pure literal substitution (DESIGN §7 — deliberately not a template
//     language). `assets` defaults to <output_dir>/assets.
//   - KeyCase is the element-key convention this project requires (default
//     Kebab). Gen never translates a key's casing — it renders whatever the
//     key already is — so this only governs what the linter and write path
//     accept.
type Config struct {
	OutputDir string
	Paths     map[string]string
	KeyCase   Casing
}

// Filename is the config file the loader looks for beside intent.yaml.
const Filename = "intent.config.yaml"

// raw mirrors the on-disk shape so absent keys stay distinguishable from
// explicit zero values (an empty output_dir falls back to the default).
type raw struct {
	OutputDir string            `yaml:"output_dir"`
	Paths     map[string]string `yaml:"paths"`
	KeyCase   string            `yaml:"key_case"`
}

// Default returns the configuration used when no file is present.
func Default() *Config {
	return withDefaults(&Config{})
}

// Load reads the config file if it exists, filling any unset field with its
// default. A missing file is not an error — defaults apply. Errors only surface
// for a present-but-unreadable or malformed file.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Default(), nil
		}
		return nil, err
	}
	var r raw
	if err := yaml.Unmarshal(b, &r); err != nil {
		return nil, err
	}
	if r.KeyCase != "" && r.KeyCase != string(Kebab) && r.KeyCase != string(Snake) {
		return nil, fmt.Errorf("key_case must be %q or %q, got %q", Kebab, Snake, r.KeyCase)
	}
	return withDefaults(&Config{OutputDir: r.OutputDir, Paths: r.Paths, KeyCase: Casing(r.KeyCase)}), nil
}

// LoadBeside loads the config that sits in the same directory as intentFile.
func LoadBeside(intentFile string) (*Config, error) {
	return Load(filepath.Join(filepath.Dir(intentFile), Filename))
}

// withDefaults fills unset fields. It mutates and returns c.
func withDefaults(c *Config) *Config {
	if c.OutputDir == "" {
		c.OutputDir = "docs"
	}
	if c.Paths == nil {
		c.Paths = map[string]string{}
	}
	if _, ok := c.Paths["assets"]; !ok {
		c.Paths["assets"] = filepath.ToSlash(filepath.Join(c.OutputDir, "assets"))
	}
	if c.KeyCase == "" {
		c.KeyCase = Kebab
	}
	return c
}
