// Package tree turns the in-memory intent model into a flat, addressable index.
//
// Phase 0 owns the raw model (load/serialize). This package sits on top of it:
// it walks the roots once into a flat list of Elements, each carrying its full
// dotted-path address, its outgoing declared edges, and enough denormalized text
// to display and search. It also builds the reverse edge map so "what references
// this?" is a lookup, not a re-walk.
//
// The two jobs everything downstream leans on are addressing (resolve a full
// path or the shortest unambiguous suffix) and the edge graph (forward edges on
// each Element, incoming refs via Incoming).
package tree

import (
	"sort"

	"github.com/colchuck-ai/intent/internal/model"
)

// Kind names an addressable element type. It is what the CLI's --type filter
// matches and what show/tree/find print in brackets.
type Kind string

const (
	KindProduct     Kind = "product"     // the product root
	KindEngineering Kind = "engineering" // the engineering root
	KindJob         Kind = "job"
	KindOutcome     Kind = "outcome"
	KindRisk        Kind = "risk"
	KindRequirement Kind = "requirement"
	KindComponent   Kind = "component"
	KindPrinciple   Kind = "principle"
	KindConstraint  Kind = "constraint"
	KindPDR         Kind = "pdr" // product decision record
	KindADR         Kind = "adr" // architecture decision record
	KindCR          Kind = "cr"  // change record
)

// Domain groups elements for the --domain filter.
type Domain string

const (
	DomainProduct     Domain = "product"
	DomainEngineering Domain = "engineering"
	DomainChange      Domain = "change"
)

// EdgeKind names a declared reference type (DESIGN §4).
type EdgeKind string

const (
	EdgeMitigates     EdgeKind = "mitigates"
	EdgeDependsOn     EdgeKind = "dependsOn"
	EdgeFulfills      EdgeKind = "fulfills"
	EdgeAffects       EdgeKind = "affects"
	EdgeRelationships EdgeKind = "relationships"
)

// Edge is one outgoing declared reference from an element.
type Edge struct {
	Kind EdgeKind
	To   string // full dotted address, or a root literal (product/engineering)
	Note string // only set for relationships edges
}

// Ref is the reverse of an Edge: some element points at us.
type Ref struct {
	Kind EdgeKind
	From string
	Note string
}

// Field is a type-specific structured field shown by `show` (acceptance
// criteria, a component's behavior, a record's options, ...). Exactly one of
// Text or List is populated; List being non-nil marks it a list field.
type Field struct {
	Label string
	Text  string
	List  []string
}

// Element is a single addressable node, denormalized for display and search so
// the CLI never has to reach back into the raw model.
type Element struct {
	Addr       string
	Kind       Kind
	Domain     Domain
	Name       string // display name; empty for principles/constraints (key is identity)
	ProseField string // the name of the primary prose field (statement/story/…); "" if none
	Prose      string // the primary prose text
	Detail     string
	Fields     []Field // type-specific structured content
	IsDocument bool
	Level      int    // nesting level, for outline indentation in `tree`
	Edges      []Edge // outgoing declared references
}

// RendersAsDocument reports whether the element hosts its own file/directory
// rather than being embedded inline in its parent's page. Roots, jobs, and
// records always do; the promotable types (outcome/requirement/component) do
// when their authored type is document; risks, principles, and constraints never
// do. This is the taxonomy the containment invariant (DESIGN §5) is checked
// against, kept next to the Kind constants and IsDocument that feed it.
func (e *Element) RendersAsDocument() bool {
	switch e.Kind {
	case KindProduct, KindEngineering, KindJob, KindPDR, KindADR, KindCR:
		return true
	case KindOutcome, KindRequirement, KindComponent:
		return e.IsDocument
	default:
		return false
	}
}

// Index is the flat, addressable view of one intent tree.
type Index struct {
	Root     *model.Root
	order    []*Element          // elements in tree (insertion) order
	byAddr   map[string]*Element // address -> element
	incoming map[string][]Ref    // address -> refs pointing at it
}

// Build walks the model into an index. It preserves insertion order (DESIGN §3)
// and records every declared edge in both directions.
func Build(r *model.Root) *Index {
	ix := &Index{
		Root:     r,
		byAddr:   map[string]*Element{},
		incoming: map[string][]Ref{},
	}
	ix.walkProduct()
	ix.walkEngineering()
	ix.walkChangeRecords()
	ix.buildIncoming()
	return ix
}

func (ix *Index) add(e *Element) {
	ix.order = append(ix.order, e)
	ix.byAddr[e.Addr] = e
}

func (ix *Index) walkProduct() {
	p := ix.Root.Product
	ix.add(&Element{
		Addr: "product", Kind: KindProduct, Domain: DomainProduct,
		Name: p.Name, ProseField: "summary", Prose: p.Summary, Detail: p.Detail, Level: 0,
	})

	for _, jk := range p.Jobs.Keys() {
		job, _ := p.Jobs.Get(jk)
		jAddr := "product.jobs." + jk
		ix.add(&Element{
			Addr: jAddr, Kind: KindJob, Domain: DomainProduct,
			Name: job.Name, ProseField: "story", Prose: job.Story, Detail: job.Detail, Level: 1,
		})

		for _, ok := range job.Outcomes.Keys() {
			oc, _ := job.Outcomes.Get(ok)
			oAddr := jAddr + ".outcomes." + ok
			ix.add(&Element{
				Addr: oAddr, Kind: KindOutcome, Domain: DomainProduct,
				Name: oc.Name, ProseField: "statement", Prose: oc.Statement, Detail: oc.Detail,
				IsDocument: oc.Type == "document", Level: 2,
			})

			for _, rk := range oc.Risks.Keys() {
				risk, _ := oc.Risks.Get(rk)
				ix.add(&Element{
					Addr: oAddr + ".risks." + rk, Kind: KindRisk, Domain: DomainProduct,
					Name: risk.Name, ProseField: "statement", Prose: risk.Statement, Detail: risk.Detail, Level: 3,
				})
			}

			for _, reqk := range oc.Requirements.Keys() {
				req, _ := oc.Requirements.Get(reqk)
				e := &Element{
					Addr: oAddr + ".requirements." + reqk, Kind: KindRequirement, Domain: DomainProduct,
					Name: req.Name, ProseField: "statement", Prose: req.Statement, Detail: req.Detail,
					IsDocument: req.Type == "document", Level: 3,
				}
				for _, m := range req.Mitigates {
					e.Edges = append(e.Edges, Edge{Kind: EdgeMitigates, To: m})
				}
				for _, d := range req.DependsOn {
					e.Edges = append(e.Edges, Edge{Kind: EdgeDependsOn, To: d})
				}
				if len(req.AcceptanceCriteria) > 0 {
					e.Fields = append(e.Fields, Field{Label: "acceptance_criteria", List: req.AcceptanceCriteria})
				}
				ix.add(e)
			}
		}
	}

	for _, dk := range p.DecisionRecords.Keys() {
		dr, _ := p.DecisionRecords.Get(dk)
		ix.add(decisionRecordElement("product.decision_records."+dk, KindPDR, DomainProduct, dr))
	}
}

func (ix *Index) walkEngineering() {
	e := ix.Root.Engineering
	ix.add(&Element{
		Addr: "engineering", Kind: KindEngineering, Domain: DomainEngineering,
		Name: e.Name, ProseField: "summary", Prose: e.Summary, Detail: e.Detail, Level: 0,
	})

	for _, k := range e.Principles.Keys() {
		v, _ := e.Principles.Get(k)
		ix.add(&Element{
			Addr: "engineering.principles." + k, Kind: KindPrinciple, Domain: DomainEngineering,
			ProseField: "statement", Prose: v, Level: 1,
		})
	}
	for _, k := range e.Constraints.Keys() {
		v, _ := e.Constraints.Get(k)
		ix.add(&Element{
			Addr: "engineering.constraints." + k, Kind: KindConstraint, Domain: DomainEngineering,
			ProseField: "statement", Prose: v, Level: 1,
		})
	}

	for _, ck := range e.Components.Keys() {
		c, _ := e.Components.Get(ck)
		el := &Element{
			Addr: "engineering.components." + ck, Kind: KindComponent, Domain: DomainEngineering,
			Name: c.Name, ProseField: "responsibility", Prose: c.Responsibility, Detail: c.Detail,
			IsDocument: c.Type == "document", Level: 1,
		}
		for _, f := range c.Fulfills {
			el.Edges = append(el.Edges, Edge{Kind: EdgeFulfills, To: f})
		}
		for _, rk := range c.Relationships.Keys() {
			note, _ := c.Relationships.Get(rk)
			el.Edges = append(el.Edges, Edge{Kind: EdgeRelationships, To: rk, Note: note})
		}
		el.Fields = appendText(el.Fields, "data_model", c.DataModel)
		el.Fields = appendText(el.Fields, "interfaces", c.Interfaces)
		el.Fields = appendText(el.Fields, "behavior", c.Behavior)
		el.Fields = appendList(el.Fields, "edge_cases", c.EdgeCases)
		el.Fields = appendList(el.Fields, "success_criteria", c.SuccessCriteria)
		ix.add(el)
	}

	for _, dk := range e.DecisionRecords.Keys() {
		dr, _ := e.DecisionRecords.Get(dk)
		ix.add(decisionRecordElement("engineering.decision_records."+dk, KindADR, DomainEngineering, dr))
	}
}

func (ix *Index) walkChangeRecords() {
	for _, ck := range ix.Root.ChangeRecords.Keys() {
		cr, _ := ix.Root.ChangeRecords.Get(ck)
		el := &Element{
			Addr: "change_records." + ck, Kind: KindCR, Domain: DomainChange,
			Name: cr.Name, ProseField: "change", Prose: cr.Change, Level: 0,
		}
		for _, a := range cr.Affects {
			el.Edges = append(el.Edges, Edge{Kind: EdgeAffects, To: a})
		}
		el.Fields = appendText(el.Fields, "rationale", cr.Rationale)
		ix.add(el)
	}
}

func decisionRecordElement(addr string, kind Kind, domain Domain, dr model.DecisionRecord) *Element {
	e := &Element{
		Addr: addr, Kind: kind, Domain: domain,
		Name: dr.Name, ProseField: "summary", Prose: dr.Summary, Level: 1,
	}
	for _, a := range dr.Affects {
		e.Edges = append(e.Edges, Edge{Kind: EdgeAffects, To: a})
	}
	e.Fields = appendText(e.Fields, "context", dr.Context)
	e.Fields = appendList(e.Fields, "options", dr.Options)
	e.Fields = appendText(e.Fields, "decision", dr.Decision)
	e.Fields = appendList(e.Fields, "consequences", dr.Consequences)
	return e
}

func appendText(fields []Field, label, text string) []Field {
	if text == "" {
		return fields
	}
	return append(fields, Field{Label: label, Text: text})
}

func appendList(fields []Field, label string, list []string) []Field {
	if len(list) == 0 {
		return fields
	}
	return append(fields, Field{Label: label, List: list})
}

func (ix *Index) buildIncoming() {
	for _, e := range ix.order {
		for _, edge := range e.Edges {
			ix.incoming[edge.To] = append(ix.incoming[edge.To], Ref{Kind: edge.Kind, From: e.Addr, Note: edge.Note})
		}
	}
}

// All returns every element in tree order. The slice is the internal one; treat
// it as read-only.
func (ix *Index) All() []*Element { return ix.order }

// Get looks up an element by its exact full address.
func (ix *Index) Get(addr string) (*Element, bool) {
	e, ok := ix.byAddr[addr]
	return e, ok
}

// Incoming returns the refs that point at addr, in a stable order sorted by the
// source address (DESIGN §8: derived views are address-sorted).
func (ix *Index) Incoming(addr string) []Ref {
	refs := append([]Ref(nil), ix.incoming[addr]...)
	sort.Slice(refs, func(i, j int) bool {
		if refs[i].From != refs[j].From {
			return refs[i].From < refs[j].From
		}
		return refs[i].Kind < refs[j].Kind
	})
	return refs
}
