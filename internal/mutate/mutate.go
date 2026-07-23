// Package mutate is the write path: the surgical, validated edits the CLI makes
// to an intent.yaml model (DESIGN §11). Reads go through a denormalized index
// (package tree); writes come here, where they touch the live *model.Root.
//
// The one hard problem is addressing a struct field that lives deep in the
// nested model. Go map entries aren't addressable, so locate walks the tree,
// copying each level, and returns a pointer to the (copied) target plus a save
// closure that writes the whole chain back. Every verb — Add, Set, Remove, Link,
// Unlink — is built on that pair. The caller is responsible for re-validating
// and serializing the result (only the whole-tree check can prove a write kept
// the tree valid); mutate just performs the structural edit and reports what it
// couldn't do.
package mutate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
)

// keyRe is the snake_case element-key format (mirrors the schema's `key`).
var keyRe = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Fields carries the optional field values a new element may be created with.
// Only the fields relevant to the kind being added are read; the rest are
// ignored. Edge fields hold fully-resolved dotted addresses (the CLI resolves
// suffixes before calling Add).
type Fields struct {
	Name           string
	Story          string // job
	Statement      string // outcome/risk/requirement; also the principle/constraint value
	Responsibility string // component
	Detail         string
	Type           string   // inline | document (promotable types)
	Mitigates      []string // requirement (required, ≥1)
	DependsOn      []string // requirement
	Fulfills       []string // component (required, ≥1)
	Acceptance     []string // requirement
}

// Add inserts a new element of kind under parentAddr with key, populated from f.
// It errors if the key isn't snake_case, the kind can't nest under that parent,
// or the key already exists there. It returns the new element's full address.
//
// Required edges (a requirement's mitigates, a component's fulfills) must be
// supplied here: the schema rejects an element missing them, so the CLI would
// refuse to write a two-step "add then link". Add checks for them up front to
// give a clearer message than the schema would.
func Add(r *model.Root, kind tree.Kind, parentAddr, key string, f Fields) (string, error) {
	if !keyRe.MatchString(key) {
		return "", fmt.Errorf("key %q must be snake_case (letters, digits, underscores; starting with a letter)", key)
	}
	if kind == tree.KindRequirement && len(f.Mitigates) == 0 {
		return "", fmt.Errorf("a requirement needs at least one --mitigates target")
	}
	if kind == tree.KindComponent && len(f.Fulfills) == 0 {
		return "", fmt.Errorf("a component needs at least one --fulfills target")
	}
	if f.Type != "" {
		if err := checkType(f.Type); err != nil {
			return "", err
		}
	}

	parent, save, err := locate(r, parentAddr)
	if err != nil {
		return "", fmt.Errorf("parent %s: %w", parentAddr, err)
	}
	seg := collectionSeg(kind)
	if seg == "" {
		return "", fmt.Errorf("cannot add a %s with this command", kind)
	}
	addr := parentAddr + "." + seg + "." + key

	switch p := parent.(type) {
	case *model.Product:
		if kind != tree.KindJob {
			return "", parentMismatch(kind, parentAddr)
		}
		if _, ok := p.Jobs.Get(key); ok {
			return "", existsErr(addr)
		}
		p.Jobs.Set(key, model.Job{Name: f.Name, Story: f.Story, Detail: f.Detail})
	case *model.Job:
		if kind != tree.KindOutcome {
			return "", parentMismatch(kind, parentAddr)
		}
		if _, ok := p.Outcomes.Get(key); ok {
			return "", existsErr(addr)
		}
		p.Outcomes.Set(key, model.Outcome{Name: f.Name, Type: f.Type, Statement: f.Statement, Detail: f.Detail})
	case *model.Outcome:
		switch kind {
		case tree.KindRisk:
			if _, ok := p.Risks.Get(key); ok {
				return "", existsErr(addr)
			}
			p.Risks.Set(key, model.Risk{Name: f.Name, Statement: f.Statement, Detail: f.Detail})
		case tree.KindRequirement:
			if _, ok := p.Requirements.Get(key); ok {
				return "", existsErr(addr)
			}
			p.Requirements.Set(key, model.Requirement{
				Name: f.Name, Type: f.Type, Statement: f.Statement,
				Mitigates: f.Mitigates, DependsOn: f.DependsOn,
				AcceptanceCriteria: f.Acceptance, Detail: f.Detail,
			})
		default:
			return "", parentMismatch(kind, parentAddr)
		}
	case *model.Engineering:
		switch kind {
		case tree.KindComponent:
			if _, ok := p.Components.Get(key); ok {
				return "", existsErr(addr)
			}
			p.Components.Set(key, model.Component{
				Name: f.Name, Type: f.Type, Responsibility: f.Responsibility,
				Fulfills: f.Fulfills, Detail: f.Detail,
			})
		case tree.KindPrinciple:
			if _, ok := p.Principles.Get(key); ok {
				return "", existsErr(addr)
			}
			p.Principles.Set(key, f.Statement)
		case tree.KindConstraint:
			if _, ok := p.Constraints.Get(key); ok {
				return "", existsErr(addr)
			}
			p.Constraints.Set(key, f.Statement)
		default:
			return "", parentMismatch(kind, parentAddr)
		}
	default:
		return "", fmt.Errorf("cannot add a %s under %s", kind, parentAddr)
	}
	save()
	return addr, nil
}

// Set updates a single scalar field on the element at addr. List fields
// (acceptance_criteria, a record's options, ...) and edges are not settable here
// — edges go through Link — so an attempt to set one is a clear error.
func Set(r *model.Root, addr, field, value string) error {
	el, save, err := locate(r, addr)
	if err != nil {
		return err
	}
	if err := setField(el, addr, field, value); err != nil {
		return err
	}
	save()
	return nil
}

// Remove deletes the element at addr (and, structurally, its whole subtree). It
// does not itself guard against dangling references — the caller checks incoming
// refs first for a targeted message, and the post-write validate is the backstop
// (DESIGN §9).
func Remove(r *model.Root, addr string) error {
	parentAddr := tree.ParentAddr(addr)
	key := tree.LastSegment(addr)

	// Top-level change records have no parent element.
	segs := strings.Split(addr, ".")
	if len(segs) == 2 && segs[0] == "change_records" {
		if !r.ChangeRecords.Delete(key) {
			return notFound(addr)
		}
		return nil
	}
	if parentAddr == "" {
		return fmt.Errorf("cannot remove the root %q", addr)
	}

	parent, save, err := locate(r, parentAddr)
	if err != nil {
		return err
	}
	deleted := deleteChild(parent, collectionSegOf(addr), key)
	if !deleted {
		return notFound(addr)
	}
	save()
	return nil
}

// Link adds an outgoing edge of kind from the element at `from` to `to` (a full
// dotted address, or a root literal for affects). note is used only for a
// relationships edge. It errors if the element type can't carry that edge kind
// or the edge is already declared (a duplicate would fail validation as E003).
func Link(r *model.Root, kind tree.EdgeKind, from, to, note string) error {
	el, save, err := locate(r, from)
	if err != nil {
		return err
	}
	if err := addEdge(el, from, kind, to, note); err != nil {
		return err
	}
	save()
	return nil
}

// Unlink removes an outgoing edge of kind from `from` to `to`. It errors if no
// such edge is declared.
func Unlink(r *model.Root, kind tree.EdgeKind, from, to string) error {
	el, save, err := locate(r, from)
	if err != nil {
		return err
	}
	if err := removeEdge(el, from, kind, to); err != nil {
		return err
	}
	save()
	return nil
}

// collectionSeg maps an addable kind to the map segment it lives under, used to
// compute a new element's address. Records are created via `record` (Phase 5),
// so they are not addable here and map to "".
func collectionSeg(kind tree.Kind) string {
	switch kind {
	case tree.KindJob:
		return "jobs"
	case tree.KindOutcome:
		return "outcomes"
	case tree.KindRisk:
		return "risks"
	case tree.KindRequirement:
		return "requirements"
	case tree.KindComponent:
		return "components"
	case tree.KindPrinciple:
		return "principles"
	case tree.KindConstraint:
		return "constraints"
	default:
		return ""
	}
}

// collectionSegOf returns the map segment that holds the element at addr — the
// second-to-last dotted segment (e.g. "risks" in ...outcomes.O.risks.R).
func collectionSegOf(addr string) string {
	segs := strings.Split(addr, ".")
	if len(segs) < 2 {
		return ""
	}
	return segs[len(segs)-2]
}

func parentMismatch(kind tree.Kind, parentAddr string) error {
	return fmt.Errorf("a %s cannot be added under %s", kind, parentAddr)
}

func existsErr(addr string) error {
	return fmt.Errorf("%s already exists", addr)
}

func notFound(addr string) error {
	return &tree.NotFoundError{Query: addr}
}
