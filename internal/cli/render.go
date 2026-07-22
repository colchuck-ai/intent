package cli

import (
	"fmt"
	"strings"

	"github.com/colchuck-ai/intent/internal/tree"
)

// displayName is the element's name, falling back to its key when there is none
// (principles/constraints have no name — the key is identity).
func displayName(e *tree.Element) string {
	if e.Name != "" {
		return e.Name
	}
	return tree.LastSegment(e.Addr)
}

// kindLabel renders "[kind] <name> (document)". The name is passed in so callers
// choose whether to fall back to the key: `show`/`find` want a name always,
// `tree` prints the key separately and omits a redundant one.
func kindLabel(e *tree.Element, name string) string {
	s := "[" + string(e.Kind) + "]"
	if name != "" {
		s += " " + name
	}
	if e.IsDocument {
		s += " (document)"
	}
	return s
}

// label is the head-of-element tag used by show/find/affects/trace, with the
// name always present.
func label(e *tree.Element) string { return kindLabel(e, displayName(e)) }

// nameOf returns the display name of an address if it resolves in the index,
// else "" — used to annotate edge targets with their name.
func nameOf(ix *tree.Index, addr string) string {
	if e, ok := ix.Get(addr); ok {
		return displayName(e)
	}
	return ""
}

// refLine renders one edge/ref line shared by show, affects, and trace:
//
//	<indent><arrow> <kind>  <addr>  (Name)  — note
//
// arrow is → for an outgoing edge and ← for an incoming ref. The name and note
// clauses are added only when present.
func refLine(ix *tree.Index, indent, arrow string, kind tree.EdgeKind, addr, note string) string {
	line := fmt.Sprintf("%s%s %s  %s", indent, arrow, kind, addr)
	if name := nameOf(ix, addr); name != "" {
		line += "  (" + name + ")"
	}
	if note != "" {
		line += "  — " + note
	}
	return line
}

// matchSnippet returns a one-line excerpt of whichever field actually contains
// the query (prose, then detail, then name), so `find` shows the matched
// context rather than an unrelated prose clip. When the query matched only the
// address (or is empty), it falls back to the prose as a summary.
func matchSnippet(e *tree.Element, query string, width int) string {
	for _, f := range []string{e.Prose, e.Detail, e.Name} {
		if f == "" {
			continue
		}
		if query == "" || strings.Contains(strings.ToLower(f), strings.ToLower(query)) {
			return snippet(f, query, width)
		}
	}
	return snippet(e.Prose, query, width)
}

// snippet returns a one-line excerpt of text centered on the first occurrence of
// query (case-insensitive), clipped to width runes with ellipses. With an empty
// query it clips from the start. Newlines collapse to spaces so the excerpt
// stays on one line.
func snippet(text, query string, width int) string {
	text = strings.Join(strings.Fields(text), " ")
	if text == "" {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}

	start := 0
	if query != "" {
		if i := strings.Index(strings.ToLower(text), strings.ToLower(query)); i > 0 {
			// Convert byte offset to rune offset, then center the window.
			ri := len([]rune(text[:i]))
			start = ri - width/3
			if start < 0 {
				start = 0
			}
		}
	}
	end := start + width
	if end > len(runes) {
		end = len(runes)
		start = end - width
	}

	out := string(runes[start:end])
	if start > 0 {
		out = "…" + out
	}
	if end < len(runes) {
		out = out + "…"
	}
	return out
}
