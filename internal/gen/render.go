package gen

import (
	"strings"

	"github.com/colchuck-ai/intent/internal/interp"
	"github.com/colchuck-ai/intent/internal/tree"
)

// edgeOrder fixes how outgoing declared edges are grouped on a page, so the
// section never reorders (DESIGN §8).
var edgeOrder = []tree.EdgeKind{
	tree.EdgeFulfills, tree.EdgeDependsOn, tree.EdgeMitigates,
	tree.EdgeAffects, tree.EdgeRelationships,
}

// edgeLabel is the human heading for each edge kind.
var edgeLabel = map[tree.EdgeKind]string{
	tree.EdgeFulfills:      "Fulfills",
	tree.EdgeDependsOn:     "Depends on",
	tree.EdgeMitigates:     "Mitigates",
	tree.EdgeAffects:       "Affects",
	tree.EdgeRelationships: "Relationships",
}

// siblingKindLabel names the element kinds that render as undifferentiated
// sibling headings with no other structural cue to tell them apart: a
// principle, a constraint, and a component all sit directly under the
// engineering root, and a risk and a requirement both sit directly under an
// outcome. It is deliberately partial, unlike edgeLabel: every other kind is
// unambiguous from where it renders (a job's only children are outcomes; a
// decision record has its own unmistakable field shape), so it gets no
// label.
var siblingKindLabel = map[tree.Kind]string{
	tree.KindPrinciple:   "Principle",
	tree.KindConstraint:  "Constraint",
	tree.KindComponent:   "Component",
	tree.KindRisk:        "Risk",
	tree.KindRequirement: "Requirement",
}

// renderPage assembles one page: banner, then the authored core + declared edges
// of the host document and each inline element it hosts, then a single derived
// footer (DESIGN §8 page anatomy).
func (s *site) renderPage(host string) []byte {
	file := s.file[host]
	var b strings.Builder
	b.WriteString(Banner)
	b.WriteByte('\n')

	els := s.pageEls[host]
	hostLevel := els[0].Level
	for _, e := range els {
		s.renderElement(&b, e, file, e.Level-hostLevel)
	}

	s.renderFooter(&b, host, file)

	out := strings.TrimRight(b.String(), "\n") + "\n"
	return []byte(out)
}

// renderElement writes one element's authored core (heading, prose, structured
// fields, detail) and its outgoing declared-edge section. depth is the element's
// nesting below the page host (0 = the host itself), which sets the heading
// level.
func (s *site) renderElement(b *strings.Builder, e *tree.Element, file string, depth int) {
	level := depth + 1
	if level > 6 {
		level = 6
	}
	b.WriteByte('\n')
	b.WriteString(strings.Repeat("#", level) + " " + displayName(e) + "\n")

	if label, ok := siblingKindLabel[e.Kind]; ok {
		b.WriteString("\n**Kind:** " + label + "\n")
	}

	if e.Prose != "" {
		b.WriteByte('\n')
		b.WriteString(s.block(e.Prose, e.Addr, file))
		b.WriteByte('\n')
	}

	for _, f := range e.Fields {
		s.renderField(b, e.Addr, file, f)
	}

	if e.Detail != "" {
		b.WriteString("\n**Detail:**\n\n")
		b.WriteString(s.block(e.Detail, e.Addr, file))
		b.WriteByte('\n')
	}

	s.renderEdges(b, e, file)
}

// renderField writes one type-specific structured field: a labeled paragraph for
// text, a labeled bullet list for a list.
func (s *site) renderField(b *strings.Builder, owner, file string, f tree.Field) {
	if f.List != nil {
		b.WriteString("\n**" + fieldLabel(f.Label) + ":**\n\n")
		for _, item := range f.List {
			b.WriteString("- " + s.interpolate(item, owner, file) + "\n")
		}
		return
	}
	if f.Text != "" {
		b.WriteString("\n**" + fieldLabel(f.Label) + ":**\n\n")
		b.WriteString(s.block(f.Text, owner, file))
		b.WriteByte('\n')
	}
}

// renderEdges writes the outgoing declared references, grouped by kind in fixed
// order and rendered as links (DESIGN §8 declared-edge section).
func (s *site) renderEdges(b *strings.Builder, e *tree.Element, file string) {
	byKind := map[tree.EdgeKind][]tree.Edge{}
	for _, edge := range e.Edges {
		byKind[edge.Kind] = append(byKind[edge.Kind], edge)
	}
	for _, kind := range edgeOrder {
		edges := byKind[kind]
		if len(edges) == 0 {
			continue
		}
		b.WriteString("\n**" + edgeLabel[kind] + ":**\n\n")
		for _, edge := range edges {
			line := "- " + s.link(file, edge.To)
			if edge.Note != "" {
				line += " — " + s.interpolate(edge.Note, e.Addr, file)
			}
			b.WriteString(line + "\n")
		}
	}
}

// renderFooter writes the derived zone (DESIGN §8): a `---` rule then, in fixed
// order, "Referenced by" (incoming edges), "See also" (shared-target siblings),
// and — for a directory page — a "Contents" map of child pages. Every section is
// address-sorted, and the whole zone is omitted when all three are empty.
func (s *site) renderFooter(b *strings.Builder, host, file string) {
	refs := s.groupedSection(host, file, "Referenced by", s.referencedByLines)
	see := s.groupedSection(host, file, "See also", s.seeAlsoLines)
	contents := s.renderContents(host, file)
	if refs == "" && see == "" && contents == "" {
		return
	}
	b.WriteString("\n---\n")
	b.WriteString(refs)
	b.WriteString(see)
	b.WriteString(contents)
}

// groupedSection renders one per-element footer section. It calls lines for each
// element the page hosts; an element contributing no lines is skipped. On a
// single-element page the lines are listed flat (the owner is the page itself);
// on a multi-element page each element's lines are grouped under its name.
func (s *site) groupedSection(host, file, heading string, lines func(*tree.Element, string) []string) string {
	els := s.pageEls[host]
	var b strings.Builder
	for _, e := range els {
		ls := lines(e, file)
		if len(ls) == 0 {
			continue
		}
		if b.Len() == 0 {
			b.WriteString("\n## " + heading + "\n\n")
		}
		if len(els) > 1 {
			b.WriteString("**" + displayName(e) + "**\n\n")
		}
		for _, l := range ls {
			b.WriteString(l + "\n")
		}
		if len(els) > 1 {
			b.WriteByte('\n')
		}
	}
	if b.Len() == 0 {
		return ""
	}
	return strings.TrimRight(b.String(), "\n") + "\n"
}

// referencedByLines are the incoming-edge bullets for one element (backlinks).
func (s *site) referencedByLines(e *tree.Element, file string) []string {
	var out []string
	for _, r := range s.ix.Incoming(e.Addr) {
		line := "- " + s.link(file, r.From) + " (" + string(r.Kind) + ")"
		if r.Note != "" {
			line += " — " + s.interpolate(r.Note, r.From, file)
		}
		out = append(out, line)
	}
	return out
}

// seeAlsoLines are the shared-target-sibling bullets for one element.
func (s *site) seeAlsoLines(e *tree.Element, file string) []string {
	var out []string
	for _, addr := range s.seeAlso(e.Addr) {
		out = append(out, "- "+s.link(file, addr))
	}
	return out
}

// renderContents lists a directory page's child pages as links (the map).
func (s *site) renderContents(host, file string) string {
	children := s.docChildren(host)
	if len(children) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n## Contents\n\n")
	for _, c := range children {
		b.WriteString("- " + s.link(file, c.Addr) + "\n")
	}
	return b.String()
}

// block interpolates a multi-line text field and trims trailing whitespace, so
// a YAML block scalar's trailing newline never leaks into an extra blank line —
// spacing stays uniform however the field was authored.
func (s *site) block(text, owner, file string) string {
	return strings.TrimRight(s.interpolate(text, owner, file), " \t\n")
}

// interpolate resolves every {{ }} expression in text. text belongs to owner
// (for relative addresses); file is the page being written (for relative hrefs).
func (s *site) interpolate(text, owner, file string) string {
	return interp.Interpolate(text, owner, func(e interp.Expr) string {
		if e.IsPath {
			if v, ok := s.cfg.Paths[e.PathVar]; ok {
				return v
			}
			return "{{paths." + e.PathVar + "}}" // unresolved: keep it visible
		}
		if e.Addr == "" {
			return ""
		}
		switch e.Accessor {
		case "name":
			return s.nameFor(e.Addr)
		case "path":
			return s.href(file, e.Addr)
		default: // "link" or a bare address
			return s.link(file, e.Addr)
		}
	})
}

// link renders a markdown hyperlink to addr from the page at file; it falls back
// to plain name text when the target can't be resolved to a page.
func (s *site) link(file, addr string) string {
	name := s.nameFor(addr)
	href := s.href(file, addr)
	if href == "" {
		return name
	}
	return "[" + name + "](" + href + ")"
}

// href is the relative markdown href from the page at file to addr, with an
// in-page anchor when the target is an inline element. Same-page targets yield a
// bare fragment.
func (s *site) href(file, addr string) string {
	host, ok := s.host[addr]
	if !ok {
		return ""
	}
	target := s.file[host]
	if target == "" {
		return ""
	}
	anchor := s.anchor[addr]
	if target == file {
		if anchor == "" {
			return "#"
		}
		return "#" + anchor
	}
	rel := slashRel(file, target)
	if anchor != "" {
		rel += "#" + anchor
	}
	return rel
}

// nameFor is the display name of an address, or the address itself if unknown.
func (s *site) nameFor(addr string) string {
	if e, ok := s.ix.Get(addr); ok {
		return displayName(e)
	}
	return addr
}

// fieldLabel humanizes a snake_case field key for a heading
// (acceptance_criteria -> "Acceptance criteria").
func fieldLabel(key string) string {
	words := strings.Split(key, "_")
	if len(words[0]) > 0 {
		words[0] = strings.ToUpper(words[0][:1]) + words[0][1:]
	}
	return strings.Join(words, " ")
}
