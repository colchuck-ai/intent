// Package interp is the one owner of the {{ }} interpolation syntax (DESIGN §6).
//
// An expression is disambiguated by prefix: a `paths.*` config variable, an
// absolute element address (product./engineering./change_records.), or a
// relative one (leading dot, resolved against the owning element's container).
// A reserved terminal accessor (.name/.path/.link) selects what to render; it is
// not part of the address.
//
// Two consumers share this parser so the syntax can't drift between them:
// validation (Phase 2) needs only the element addresses, to check that each one
// resolves (E001); generation (Phase 3) needs the parsed expression plus a way
// to rewrite prose, so it can turn an address+accessor into text or a markdown
// link. ParseExpr is the shared core; Addresses and Interpolate sit on top.
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

// Expr is one parsed {{ }} expression. Exactly one of the two forms is set:
// a config path variable (IsPath, with PathVar) or an element reference (Addr,
// with an optional Accessor). An expression that carries no address at all — a
// lone accessor like {{.link}} — parses to the zero Expr (Addr and PathVar both
// empty).
type Expr struct {
	Raw      string // the original inner text, e.g. ".reference_guides.link"
	IsPath   bool   // true for a paths.* config variable
	PathVar  string // the variable name after "paths." (only when IsPath)
	Addr     string // the resolved full dotted element address (accessor stripped)
	Accessor string // the terminal accessor: name | path | link, or "" if none
}

// Ref is one element address found in prose, kept for validation callers.
type Ref struct {
	Expr string // the original inner expression, e.g. ".reference_guides.link"
	Addr string // the resolved full dotted address (accessor stripped)
}

// ParseExpr parses one inner expression against ownerAddr (the address of the
// element the prose belongs to), used to resolve relative (leading-dot)
// addresses against its container.
func ParseExpr(inner, ownerAddr string) Expr {
	inner = strings.TrimSpace(inner)
	e := Expr{Raw: inner}

	if inner == "paths" || strings.HasPrefix(inner, "paths.") {
		e.IsPath = true
		e.PathVar = strings.TrimPrefix(inner, "paths.")
		if e.PathVar == "paths" { // the bare "paths" case
			e.PathVar = ""
		}
		return e
	}

	body := inner
	if strings.HasPrefix(body, ".") {
		body = strings.TrimPrefix(body, ".")
		addr, acc := splitAccessor(body)
		if addr == "" {
			return e // a lone accessor carries no address
		}
		e.Addr = container(ownerAddr) + "." + addr
		e.Accessor = acc
		return e
	}

	addr, acc := splitAccessor(body)
	e.Addr = addr
	e.Accessor = acc
	return e
}

// Addresses returns the element addresses interpolated in prose, in order of
// appearance. paths.* expressions and lone accessors carry no address and are
// omitted.
func Addresses(prose, ownerAddr string) []Ref {
	var out []Ref
	for _, m := range exprPattern.FindAllStringSubmatch(prose, -1) {
		e := ParseExpr(m[1], ownerAddr)
		if e.Addr == "" {
			continue
		}
		out = append(out, Ref{Expr: e.Raw, Addr: e.Addr})
	}
	return out
}

// Interpolate rewrites every {{ }} expression in prose by calling repl with the
// parsed Expr and substituting its return value. It is the rendering entry point
// for generation.
func Interpolate(prose, ownerAddr string, repl func(Expr) string) string {
	return exprPattern.ReplaceAllStringFunc(prose, func(m string) string {
		inner := exprPattern.FindStringSubmatch(m)[1]
		return repl(ParseExpr(inner, ownerAddr))
	})
}

// splitAccessor separates a trailing .name/.path/.link accessor from a dotted
// address. It returns the address and the accessor ("" if none). A string that
// is only an accessor returns an empty address.
func splitAccessor(addr string) (string, string) {
	if i := strings.LastIndex(addr, "."); i >= 0 && accessors[addr[i+1:]] {
		return addr[:i], addr[i+1:]
	}
	if accessors[addr] {
		return "", addr
	}
	return addr, ""
}

// container is the namespace a relative address resolves within: the owner's
// parent for a nested element, or the root itself when the owner is a root.
func container(ownerAddr string) string {
	if i := strings.LastIndex(ownerAddr, "."); i >= 0 {
		return ownerAddr[:i]
	}
	return ownerAddr
}
