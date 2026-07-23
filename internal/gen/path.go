package gen

import (
	"path"
	"strings"
	"unicode"
)

// pathSegments turns an element address into the directory segments of its page,
// relative to the output dir. The collection segments in the address (jobs,
// outcomes, requirements, components) are dropped — directory nesting comes from
// document ancestry, and the containment invariant (DESIGN §5) guarantees every
// ancestor of a document is itself a document, so every element key in the
// address is a real directory level. Two names are remapped for the reader:
// change_records → change-records (its own top directory) and decision_records
// → drs (DESIGN §8).
func pathSegments(addr string) []string {
	segs := strings.Split(addr, ".")
	if segs[0] == "change_records" {
		// change records never nest: change_records.<key>.
		return []string{"change-records", segs[1]}
	}
	out := []string{segs[0]} // product | engineering
	for i := 1; i+1 < len(segs); i += 2 {
		coll, key := segs[i], segs[i+1]
		if coll == "decision_records" {
			out = append(out, "drs", key)
		} else {
			out = append(out, key)
		}
	}
	return out
}

// filePathFor is the markdown file for a document element, relative to the
// output dir (forward slashes). A directory element (one with document children)
// lands on a README.md inside its own directory so GitHub renders it as the
// folder's landing page; a leaf document is a single <key>.md file.
func filePathFor(addr string, isDir bool) string {
	segs := pathSegments(addr)
	if isDir {
		return path.Join(append(segs, "README.md")...)
	}
	last := len(segs) - 1
	return path.Join(append(segs[:last:last], segs[last]+".md")...)
}

// slashRel is the relative href from fromFile to toFile, both forward-slash
// paths relative to the same output dir. It is filepath.Rel restricted to slash
// paths so output never depends on the host OS separator.
func slashRel(fromFile, toFile string) string {
	fromDir := splitDir(fromFile)
	toDir := splitDir(toFile)
	i := 0
	for i < len(fromDir) && i < len(toDir) && fromDir[i] == toDir[i] {
		i++
	}
	var out []string
	for j := i; j < len(fromDir); j++ {
		out = append(out, "..")
	}
	out = append(out, toDir[i:]...)
	out = append(out, path.Base(toFile))
	return strings.Join(out, "/")
}

func splitDir(file string) []string {
	d := path.Dir(file)
	if d == "." || d == "" {
		return nil
	}
	return strings.Split(d, "/")
}

// slug renders a heading's GitHub anchor: lowercased, keeping letters, numbers,
// and underscores, with spaces and hyphens collapsed to hyphens and everything
// else dropped. It matches how GitHub derives fragment ids from heading text so
// in-page links resolve.
func slug(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		switch {
		case unicode.IsLetter(r), unicode.IsNumber(r), r == '_':
			b.WriteRune(r)
		case r == ' ', r == '-':
			b.WriteByte('-')
		}
	}
	return b.String()
}
