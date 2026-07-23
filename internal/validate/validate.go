// Package validate is the intent linter: the meaning-level checks the schema
// can't express (DESIGN §9). The schema owns shape (required fields, snake_case
// keys, at-least-one edges); this package owns the four semantic rules:
//
//	E001  a declared edge or prose {{ }} address that doesn't resolve
//	E002  a document element under an inline parent (hand-edit backstop, §5)
//	E003  the same target listed twice in one edge list
//	E004  a decision record whose affects leaves its own domain (§4)
//
// Every finding carries its code so the caller can end the message with the
// `→ intent help E0NN` pointer that makes the fix one hop away.
package validate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/colchuck-ai/intent/internal/interp"
	"github.com/colchuck-ai/intent/internal/tree"
)

// Code identifies which rule a finding violates.
type Code string

const (
	E001 Code = "E001" // dangling reference
	E002 Code = "E002" // inline/document containment
	E003 Code = "E003" // duplicate reference in a list
	E004 Code = "E004" // record cross-domain affects
)

// Finding is one rule violation, located at an element address.
type Finding struct {
	Code   Code
	Addr   string
	Detail string
}

// String renders the finding as a single operator-facing line, ending with the
// help pointer required by DESIGN §9.
func (f Finding) String() string {
	return fmt.Sprintf("%s  %s: %s → intent help %s", f.Code, f.Addr, f.Detail, f.Code)
}

// Check runs every rule over the index and returns the findings sorted for
// stable output (by address, then code, then detail). An empty result means the
// tree is valid.
func Check(ix *tree.Index) []Finding {
	var fs []Finding
	for _, e := range ix.All() {
		fs = append(fs, checkDangling(ix, e)...)
		fs = append(fs, checkDuplicates(e)...)
		fs = append(fs, checkDomainScope(e)...)
	}
	fs = append(fs, checkContainment(ix)...)

	sort.Slice(fs, func(i, j int) bool {
		if fs[i].Addr != fs[j].Addr {
			return fs[i].Addr < fs[j].Addr
		}
		if fs[i].Code != fs[j].Code {
			return fs[i].Code < fs[j].Code
		}
		return fs[i].Detail < fs[j].Detail
	})
	return fs
}

// checkDangling flags declared edges and prose {{ }} addresses that don't
// resolve to a real element (E001).
func checkDangling(ix *tree.Index, e *tree.Element) []Finding {
	var fs []Finding
	for _, edge := range e.Edges {
		if _, ok := ix.Get(edge.To); !ok {
			fs = append(fs, Finding{
				Code:   E001,
				Addr:   e.Addr,
				Detail: fmt.Sprintf("declared %s reference to %q does not resolve", edge.Kind, edge.To),
			})
		}
	}
	for _, prose := range e.ProseValues() {
		for _, ref := range interp.Addresses(prose, e.Addr) {
			if _, ok := ix.Get(ref.Addr); !ok {
				fs = append(fs, Finding{
					Code:   E001,
					Addr:   e.Addr,
					Detail: fmt.Sprintf("prose reference {{%s}} does not resolve (%s)", ref.Expr, ref.Addr),
				})
			}
		}
	}
	return fs
}

// checkDuplicates flags a target listed twice within the same edge list (E003).
// relationships is a keyed map, so it dedups structurally and is exempt.
func checkDuplicates(e *tree.Element) []Finding {
	var fs []Finding
	seen := map[tree.EdgeKind]map[string]bool{}
	for _, edge := range e.Edges {
		if edge.Kind == tree.EdgeRelationships {
			continue
		}
		if seen[edge.Kind] == nil {
			seen[edge.Kind] = map[string]bool{}
		}
		if seen[edge.Kind][edge.To] {
			fs = append(fs, Finding{
				Code:   E003,
				Addr:   e.Addr,
				Detail: fmt.Sprintf("duplicate %s reference to %q", edge.Kind, edge.To),
			})
			continue
		}
		seen[edge.Kind][edge.To] = true
	}
	return fs
}

// checkDomainScope flags a decision record whose affects targets leave its own
// domain (E004). PDRs must stay under product.*, ADRs under engineering.*;
// change records are the cross-cutting exception and are exempt.
func checkDomainScope(e *tree.Element) []Finding {
	var domain string
	switch e.Kind {
	case tree.KindPDR:
		domain = "product"
	case tree.KindADR:
		domain = "engineering"
	default:
		return nil
	}
	var fs []Finding
	for _, edge := range e.Edges {
		if edge.Kind != tree.EdgeAffects {
			continue
		}
		if edge.To != domain && !strings.HasPrefix(edge.To, domain+".") {
			fs = append(fs, Finding{
				Code:   E004,
				Addr:   e.Addr,
				Detail: fmt.Sprintf("affects %q outside its domain (must be under %s)", edge.To, domain),
			})
		}
	}
	return fs
}

// checkContainment flags a document element whose parent is inline (E002). A
// document renders to its own file and needs a directory to live in; an inline
// parent has none. The CLI keeps this valid by construction (§5), so a violation
// only reaches here via a hand-edit.
func checkContainment(ix *tree.Index) []Finding {
	var fs []Finding
	for _, e := range ix.All() {
		parent := tree.ParentAddr(e.Addr)
		if parent == "" {
			continue
		}
		p, ok := ix.Get(parent)
		if !ok || !e.RendersAsDocument() || p.RendersAsDocument() {
			continue
		}
		fs = append(fs, Finding{
			Code:   E002,
			Addr:   e.Addr,
			Detail: fmt.Sprintf("document element under inline parent %s", parent),
		})
	}
	return fs
}
