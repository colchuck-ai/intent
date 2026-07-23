package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/colchuck-ai/intent/internal/config"
	"github.com/colchuck-ai/intent/internal/gen"
)

func newBuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Generate the committed markdown docs from intent.yaml",
		Long: "build renders every element into its page under the output directory\n" +
			"(intent.config.yaml output_dir, or --out). It validates first and refuses\n" +
			"to render an invalid tree, then removes any stale generated pages so the\n" +
			"output exactly matches the tree.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			files, outDir, err := buildSite(cmd)
			if err != nil {
				return err
			}
			removed, err := writeSite(outDir, files)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "wrote %d page(s) to %s\n", len(files), outDir)
			if len(removed) > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "removed %d stale page(s)\n", len(removed))
			}
			return nil
		},
	}
	cmd.Flags().String("out", "", "output directory (overrides intent.config.yaml output_dir)")
	return cmd
}

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Fail if the committed docs have drifted from intent.yaml",
		Long: "check renders the tree in memory and compares it against the docs on\n" +
			"disk, reporting any page that is out of date, missing, or stale. It writes\n" +
			"nothing and exits nonzero on drift — the CI and pre-commit gate.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			files, outDir, err := buildSite(cmd)
			if err != nil {
				return err
			}
			return checkSite(cmd.OutOrStdout(), outDir, files)
		},
	}
	cmd.Flags().String("out", "", "output directory (overrides intent.config.yaml output_dir)")
	return cmd
}

// buildSite loads and validates the tree, then renders it to an in-memory set of
// pages. It refuses to generate anything from an invalid tree: the committed
// docs must never reflect a broken state. The resolved output directory is
// --out when set, otherwise the config's output_dir.
func buildSite(cmd *cobra.Command) ([]gen.File, string, error) {
	path, _ := cmd.Flags().GetString("file")
	ix, findings, err := checkedIndex(path)
	if err != nil {
		return nil, "", err
	}
	if len(findings) > 0 {
		w := cmd.OutOrStdout()
		for _, f := range findings {
			fmt.Fprintln(w, f.String())
		}
		return nil, "", fmt.Errorf("cannot build: %d validation problem(s)", len(findings))
	}

	cfg, err := config.LoadBeside(path)
	if err != nil {
		return nil, "", err
	}
	outDir := cfg.OutputDir
	if out, _ := cmd.Flags().GetString("out"); out != "" {
		outDir = out
	}
	return gen.Build(ix, cfg), outDir, nil
}

// writeSite writes every page to disk and removes stale generated pages (files
// bearing the banner that the tree no longer produces). It returns the removed
// paths. A page whose bytes are unchanged is rewritten identically, so a
// re-run is a no-op diff.
func writeSite(outDir string, files []gen.File) (removed []string, err error) {
	want := make(map[string]bool, len(files))
	for _, f := range files {
		want[f.Path] = true
		dest := filepath.Join(outDir, filepath.FromSlash(f.Path))
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(dest, f.Content, 0o644); err != nil {
			return nil, err
		}
	}

	onDisk, err := generatedOnDisk(outDir)
	if err != nil {
		return nil, err
	}
	for rel := range onDisk {
		if want[rel] {
			continue
		}
		if err := os.Remove(filepath.Join(outDir, filepath.FromSlash(rel))); err != nil {
			return nil, err
		}
		removed = append(removed, rel)
	}
	return removed, nil
}

// checkSite compares the rendered pages against what is on disk and reports
// drift: changed, missing, or stale pages. It returns a nonzero-exit error when
// anything differs.
func checkSite(w io.Writer, outDir string, files []gen.File) error {
	var drift []string
	want := make(map[string]bool, len(files))
	for _, f := range files {
		want[f.Path] = true
		got, err := os.ReadFile(filepath.Join(outDir, filepath.FromSlash(f.Path)))
		switch {
		case errors.Is(err, fs.ErrNotExist):
			drift = append(drift, "missing  "+f.Path)
		case err != nil:
			return err
		case !bytes.Equal(got, f.Content):
			drift = append(drift, "changed  "+f.Path)
		}
	}

	onDisk, err := generatedOnDisk(outDir)
	if err != nil {
		return err
	}
	for rel := range onDisk {
		if !want[rel] {
			drift = append(drift, "stale    "+rel)
		}
	}

	if len(drift) == 0 {
		fmt.Fprintln(w, "ok — docs match intent.yaml")
		return nil
	}
	sort.Strings(drift)
	fmt.Fprintln(w, "docs have drifted from intent.yaml — run intent build:")
	for _, d := range drift {
		fmt.Fprintln(w, "  "+d)
	}
	return fmt.Errorf("%d page(s) out of date", len(drift))
}

// generatedOnDisk returns the set of markdown files under outDir that carry the
// generated banner, keyed by their forward-slash path relative to outDir. Files
// without the banner are a project's own docs and are left untouched. A missing
// output directory yields an empty set, not an error.
func generatedOnDisk(outDir string) (map[string]bool, error) {
	set := map[string]bool{}
	banner := []byte(gen.Banner)
	err := filepath.WalkDir(outDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fs.SkipAll
			}
			return err
		}
		if d.IsDir() || !strings.HasSuffix(p, ".md") {
			return nil
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if !bytes.HasPrefix(b, banner) {
			return nil
		}
		rel, err := filepath.Rel(outDir, p)
		if err != nil {
			return err
		}
		set[filepath.ToSlash(rel)] = true
		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return set, err
	}
	return set, nil
}
