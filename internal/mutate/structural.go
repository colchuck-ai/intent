package mutate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
)

// SetType changes the render type of the element at addr, keeping the
// inline/document containment invariant valid by construction (DESIGN §5) rather
// than letting an invalid state be created and then rejected.
//
//   - document: promotion cascades up — every inline promotable ancestor is
//     promoted too, so a document never sits under an inline parent (like
//     `mkdir -p`). cascade is ignored.
//   - inline: demotion is the guarded direction — it is refused when the element
//     has document descendants (that would orphan their files) unless cascade is
//     set, which demotes the whole subtree explicitly.
//
// It returns the addresses whose type actually changed (the target first, then
// any ancestors/descendants), for the CLI to report. An empty result means the
// element was already at the requested type — a no-op, not an error.
func SetType(r *model.Root, ix *tree.Index, addr, value string, cascade bool) ([]string, error) {
	if err := checkType(value); err != nil {
		return nil, err
	}
	target, ok := ix.Get(addr)
	if !ok {
		return nil, notFound(addr)
	}
	if !promotableKind(target.Kind) {
		return nil, fmt.Errorf("%s is a %s; only outcome, requirement, and component carry a type", addr, target.Kind)
	}
	if value == "document" {
		return promote(r, ix, addr)
	}
	return demote(r, ix, addr, cascade)
}

// promote sets addr and every inline promotable ancestor to document.
func promote(r *model.Root, ix *tree.Index, addr string) ([]string, error) {
	var changed []string
	for a := addr; a != ""; a = tree.ParentAddr(a) {
		e, ok := ix.Get(a)
		if !ok || !promotableKind(e.Kind) || e.IsDocument {
			continue // non-promotable ancestors (jobs, roots) already host documents
		}
		changed = append(changed, a)
	}
	for _, a := range changed {
		if err := setTypeAt(r, a, "document"); err != nil {
			return nil, err
		}
	}
	return changed, nil
}

// demote sets addr (and, with cascade, its document descendants) to inline. It
// stores inline as the empty type so it serializes by omission — the canonical,
// minimal-diff form (DESIGN §5, default inline).
func demote(r *model.Root, ix *tree.Index, addr string, cascade bool) ([]string, error) {
	target, _ := ix.Get(addr)
	if !target.IsDocument {
		return nil, nil // already inline
	}
	var docDesc []string
	for _, e := range ix.All() {
		if strings.HasPrefix(e.Addr, addr+".") && e.IsDocument {
			docDesc = append(docDesc, e.Addr)
		}
	}
	sort.Strings(docDesc)
	if len(docDesc) > 0 && !cascade {
		return nil, &demoteBlockedError{addr: addr, descendants: docDesc}
	}
	changed := append([]string{addr}, docDesc...)
	for _, a := range changed {
		if err := setTypeAt(r, a, ""); err != nil {
			return nil, err
		}
	}
	return changed, nil
}

// demoteBlockedError reports a refused demotion and names the document
// descendants that block it, pointing at the --cascade escape hatch.
type demoteBlockedError struct {
	addr        string
	descendants []string
}

func (e *demoteBlockedError) Error() string {
	return fmt.Sprintf("cannot demote %s to inline — it has %d document descendant(s): %s (pass --cascade to demote the whole subtree)",
		e.addr, len(e.descendants), strings.Join(e.descendants, ", "))
}

// setTypeAt writes the type field of the promotable element at addr.
func setTypeAt(r *model.Root, addr, value string) error {
	el, save, err := locate(r, addr)
	if err != nil {
		return err
	}
	switch v := el.(type) {
	case *model.Outcome:
		v.Type = value
	case *model.Requirement:
		v.Type = value
	case *model.Component:
		v.Type = value
	default:
		return fmt.Errorf("%s does not carry a type", addr)
	}
	save()
	return nil
}

func promotableKind(k tree.Kind) bool {
	return k == tree.KindOutcome || k == tree.KindRequirement || k == tree.KindComponent
}
