// Package interp extracts element addresses from prose interpolation.
//
// Prose fields may embed {{ }} expressions (DESIGN §6), disambiguated by prefix:
// a `paths.*` config variable, an absolute element address (product./engineering.
// /change_records.), or a relative one (leading dot, resolved against the owning
// element's container). A reserved terminal accessor (.name/.path/.link) selects
// what to render; it is not part of the address.
//
// Phase 2 needs only the addresses, so it can check that each one resolves (E001).
// The rendering side — turning an address plus accessor into text or a markdown
// link — arrives with generation in Phase 3.
package interp

import (
	"regexp"
	"strings"
)

// exprPattern matches one {{ ... }} interpolation and captures its inner text.
var exprPattern = regexp.MustCompile(`\{\{\s*(.*?)\s*\}\}`)

// accessors are the reserved terminal tokens; they select a rendering of the
// target, not a deeper address segment.
var accessors = map[string]bool{"name": true, "path": true, "link": true}

// Ref is one element address found in prose.
type Ref struct {
	Expr string // the original inner expression, e.g. ".reference_guides.link"
	Addr string // the resolved full dotted address (accessor stripped)
}

// Addresses returns the element addresses interpolated in prose, in order of
// appearance. ownerAddr is the address of the element the prose belongs to, used
// to resolve relative (leading-dot) addresses against its container. paths.*
// expressions are config path variables, not elements, and are omitted.
func Addresses(prose, ownerAddr string) []Ref {
	var out []Ref
	for _, m := range exprPattern.FindAllStringSubmatch(prose, -1) {
		expr := strings.TrimSpace(m[1])
		addr := resolve(expr, ownerAddr)
		if addr == "" {
			continue
		}
		out = append(out, Ref{Expr: expr, Addr: addr})
	}
	return out
}

// resolve turns one inner expression into a full dotted address, or "" if it is
// a config path variable or carries no address (e.g. a lone accessor).
func resolve(expr, ownerAddr string) string {
	if expr == "paths" || strings.HasPrefix(expr, "paths.") {
		return ""
	}
	if strings.HasPrefix(expr, ".") {
		body := stripAccessor(strings.TrimPrefix(expr, "."))
		if body == "" {
			return ""
		}
		return container(ownerAddr) + "." + body
	}
	return stripAccessor(expr)
}

// stripAccessor drops a trailing .name/.path/.link token from a dotted address.
func stripAccessor(addr string) string {
	if i := strings.LastIndex(addr, "."); i >= 0 && accessors[addr[i+1:]] {
		return addr[:i]
	}
	if accessors[addr] {
		return ""
	}
	return addr
}

// container is the namespace a relative address resolves within: the owner's
// parent for a nested element, or the root itself when the owner is a root.
func container(ownerAddr string) string {
	if i := strings.LastIndex(ownerAddr, "."); i >= 0 {
		return ownerAddr[:i]
	}
	return ownerAddr
}
