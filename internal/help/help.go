// Package help is the embedded, on-demand instruction layer (DESIGN §12).
//
// Content ships inside the binary as markdown files under content/, one file
// per topic, grouped into planes by directory:
//
//	concept/   what things are & why
//	ref/       factual / structural, tabular
//	judgment/  good, not just valid
//	guide/     task-routed recipes (goal → commands)
//	errors/    the E0NN recovery catalog
//
// A topic's slug is "<plane>:<name>" (the directory and the filename without
// ".md"), e.g. concept/requirement.md → "concept:requirement". Each file opens
// with a small YAML frontmatter block carrying the title and one-line summary
// used by the index; the rest is the markdown body shown on lookup.
//
// The registry is the single source of truth co-versioned with the binary, so
// help can never drift from the tool the way a separate skill file would.
package help

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed content
var contentFS embed.FS

// Topic is one help entry: its address, the plane it belongs to, the index
// metadata, and the markdown body shown on lookup.
type Topic struct {
	Slug    string `json:"slug"`
	Plane   string `json:"plane"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"-"`
}

// planeOrder fixes the display order of planes in the index (task-first: guides
// lead). A plane not listed here sorts last, alphabetically.
var planeOrder = map[string]int{
	"guide":    0,
	"concept":  1,
	"ref":      2,
	"judgment": 3,
	"errors":   4,
}

// PlaneOrder returns planes in their fixed display order, derived from
// planeOrder so that map is the single source of truth.
func PlaneOrder() []string {
	planes := make([]string, 0, len(planeOrder))
	for p := range planeOrder {
		planes = append(planes, p)
	}
	sort.Slice(planes, func(i, j int) bool { return planeOrder[planes[i]] < planeOrder[planes[j]] })
	return planes
}

// registry is built once at init from the embedded files. A parse failure in an
// embedded file is a build-time authoring bug, so we panic rather than ship a
// binary whose help is silently broken.
var registry = mustLoad()

func mustLoad() map[string]Topic {
	m, err := load()
	if err != nil {
		panic("help: loading embedded content: " + err.Error())
	}
	return m
}

func load() (map[string]Topic, error) {
	m := map[string]Topic{}
	err := fs.WalkDir(contentFS, "content", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}
		// content/<plane>/<name>.md
		rel := strings.TrimPrefix(path, "content/")
		plane, file, ok := strings.Cut(rel, "/")
		if !ok {
			return fmt.Errorf("%s: content files must live under a plane directory", path)
		}
		name := strings.TrimSuffix(file, ".md")
		slug := plane + ":" + name

		raw, err := contentFS.ReadFile(path)
		if err != nil {
			return err
		}
		t, err := parseTopic(raw)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		t.Slug = slug
		t.Plane = plane
		m[slug] = t
		return nil
	})
	if err != nil {
		return nil, err
	}
	return m, nil
}

// parseTopic splits a "---\n<yaml>\n---\n<body>" file into its frontmatter
// metadata and markdown body.
func parseTopic(raw []byte) (Topic, error) {
	var t Topic
	rest, ok := bytes.CutPrefix(raw, []byte("---\n"))
	if !ok {
		return t, fmt.Errorf("missing frontmatter (must start with '---')")
	}
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return t, fmt.Errorf("unterminated frontmatter (no closing '---')")
	}
	var meta struct {
		Title   string `yaml:"title"`
		Summary string `yaml:"summary"`
	}
	if err := yaml.Unmarshal(rest[:end], &meta); err != nil {
		return t, fmt.Errorf("parsing frontmatter: %w", err)
	}
	if meta.Title == "" {
		return t, fmt.Errorf("frontmatter missing title")
	}
	if meta.Summary == "" {
		return t, fmt.Errorf("frontmatter missing summary")
	}
	t.Title = meta.Title
	t.Summary = meta.Summary
	t.Body = strings.TrimSpace(string(rest[end+len("\n---\n"):]))
	return t, nil
}

// All returns every topic sorted by plane display-order then slug, so index and
// --json output are stable across runs.
func All() []Topic {
	ts := make([]Topic, 0, len(registry))
	for _, t := range registry {
		ts = append(ts, t)
	}
	sort.Slice(ts, func(i, j int) bool {
		pi, pj := planeRank(ts[i].Plane), planeRank(ts[j].Plane)
		if pi != pj {
			return pi < pj
		}
		return ts[i].Slug < ts[j].Slug
	})
	return ts
}

func planeRank(p string) int {
	if r, ok := planeOrder[p]; ok {
		return r
	}
	return len(planeOrder) // unknown planes sort last
}

// InPlane returns the plane's topics in slug order.
func InPlane(plane string) []Topic {
	var ts []Topic
	for _, t := range All() {
		if t.Plane == plane {
			ts = append(ts, t)
		}
	}
	return ts
}

// Resolve maps an operator-typed argument to a topic. It accepts the exact slug
// (concept:requirement), a bare error code (E001 → errors:E001, case-folded),
// and the synthetic "errors" catalog that stitches every errors:* entry into
// one page. Anything else is not found.
func Resolve(arg string) (Topic, bool) {
	if t, ok := registry[arg]; ok {
		return t, true
	}
	if isErrorCode(arg) {
		t, ok := registry["errors:"+strings.ToUpper(arg)]
		return t, ok
	}
	if arg == "errors" {
		return errorsCatalog()
	}
	return Topic{}, false
}

func isErrorCode(s string) bool {
	if len(s) < 2 || (s[0] != 'E' && s[0] != 'e') {
		return false
	}
	for _, r := range s[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// errorsCatalog synthesizes the full E0NN catalog as one topic (DESIGN §12 —
// "the full E0NN catalog"), concatenating every errors:* body in code order.
func errorsCatalog() (Topic, bool) {
	entries := InPlane("errors")
	if len(entries) == 0 {
		return Topic{}, false
	}
	var b strings.Builder
	for i, e := range entries {
		if i > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(e.Body)
	}
	return Topic{
		Slug:    "errors",
		Plane:   "errors",
		Title:   "Error catalog",
		Summary: "Every E0NN code: cause and exact fix.",
		Body:    b.String(),
	}, true
}
