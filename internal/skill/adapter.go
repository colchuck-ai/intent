// Package skill assembles the agent skill content from embedded help
// (DESIGN §12) and hands it to per-agent Adapters that emit a concrete file
// tree. The shared content is generated, never hand-edited; Adapters control
// only placement and wrapping (frontmatter, section markers, file paths).
package skill

// File is one file an Adapter emits, with Path relative to the install root.
type File struct {
	Path    string
	Content []byte
}

// Adapter turns shared Content into a concrete per-agent file tree. Concrete
// adapters (Claude Code, AGENTS.md, ...) implement this to control where and
// how the shared content lands on disk; they do not author content
// themselves, and they must not write to disk — callers install the
// returned Files.
type Adapter interface {
	// Name identifies this adapter for the --agent flag (e.g. "claude-code").
	Name() string

	// Files returns the file tree to install under root for the given
	// Content. Adapters that update an existing file in place (for example
	// an AGENTS.md fallback merging into a shared file) may read from root
	// to compute the merged result, but must return it rather than writing
	// it themselves.
	Files(root string, c Content) ([]File, error)
}
