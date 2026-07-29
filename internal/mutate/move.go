package mutate

import (
	"fmt"
	"strings"

	"github.com/colchuck-ai/intent/internal/interp"
	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
)

// MvOpts selects what an mv does. Any combination is allowed; at least one must
// be set. NewParent (equal to the current parent is treated as "no move"),
// Rename, and Before/After position the element among its siblings.
type MvOpts struct {
	NewParent string // resolved full parent address, or "" to keep the current parent
	Rename    string // new key, or "" to keep the current key
	Before    string // sibling key to sit before, or ""
	After     string // sibling key to sit after, or ""
}

// Mv reorders, renames, and/or re-parents the element at addr, then rewrites
// every reference to it — declared edges, relationship keys, and prose {{ }}
// addresses — across the whole tree so nothing dangles (DESIGN §11). It returns
// the element's new address and the number of references rewritten.
//
// A pure reorder (only Before/After) leaves the address unchanged and needs no
// cascade. Any rename or re-parent changes the address; the reference cascade
// runs first, while every element is still at its old address (so relative prose
// resolves correctly), and the structural move follows. Re-parenting is only
// meaningful for the elements with more than one possible parent — outcomes
// (between jobs) and risks/requirements (between outcomes); everything else has a
// fixed location and can only be renamed or reordered.
func Mv(r *model.Root, ix *tree.Index, addr string, opts MvOpts) (string, int, []string, error) {
	e, ok := ix.Get(addr)
	if !ok {
		return "", 0, nil, notFound(addr)
	}
	if e.Kind == tree.KindProduct || e.Kind == tree.KindEngineering {
		return "", 0, nil, fmt.Errorf("cannot move a root")
	}
	if opts.Before != "" && opts.After != "" {
		return "", 0, nil, fmt.Errorf("--before and --after are mutually exclusive")
	}

	oldParent := tree.ParentAddr(addr)
	oldKey := tree.LastSegment(addr)
	seg := collectionSegOf(addr)

	newKey := oldKey
	if opts.Rename != "" {
		if !keyRe.MatchString(opts.Rename) {
			return "", 0, nil, fmt.Errorf("key %q must be lowercase letters, digits, underscores, or hyphens, starting with a letter", opts.Rename)
		}
		newKey = opts.Rename
	}

	ref, after := opts.Before, false
	if opts.After != "" {
		ref, after = opts.After, true
	}
	ref = tree.LastSegment(ref) // accept a full address or a bare sibling key

	reparent := opts.NewParent != "" && opts.NewParent != oldParent

	// Pure reorder: no address change, so no reference cascade.
	if !reparent && opts.Rename == "" {
		if ref == "" {
			return "", 0, nil, fmt.Errorf("mv needs one of a new parent, --rename, --before, or --after")
		}
		if err := reorderInParent(r, oldParent, seg, oldKey, ref, after); err != nil {
			return "", 0, nil, err
		}
		return addr, 0, nil, nil
	}

	// Address changes — compute the new address and validate the destination.
	var newAddr string
	if reparent {
		childSeg := collectionSeg(e.Kind)
		if childSeg == "" {
			return "", 0, nil, fmt.Errorf("a %s has a fixed location; it can be renamed or reordered but not re-parented", e.Kind)
		}
		if opts.NewParent == addr || strings.HasPrefix(opts.NewParent, addr+".") {
			return "", 0, nil, fmt.Errorf("cannot move %s into itself", addr)
		}
		parent, ok := ix.Get(opts.NewParent)
		if !ok {
			return "", 0, nil, notFound(opts.NewParent)
		}
		if !canNest(e.Kind, parent.Kind) {
			return "", 0, nil, fmt.Errorf("a %s cannot be re-parented under a %s (%s)", e.Kind, parent.Kind, opts.NewParent)
		}
		newAddr = opts.NewParent + "." + childSeg + "." + newKey
	} else {
		newAddr = replaceLast(addr, newKey)
	}
	if newAddr != addr {
		if _, exists := ix.Get(newAddr); exists {
			return "", 0, nil, existsErr(newAddr)
		}
	}

	// If the moved subtree renders any document, its new parent must render as a
	// document too — decided now, against the pre-move index, before addresses
	// change.
	carriesDocument := reparent && subtreeRendersDocument(ix, addr)

	n := remapRefs(r, addr, newAddr)

	var promoted []string
	if reparent {
		if err := relocate(r, e.Kind, oldParent, oldKey, opts.NewParent, newKey, ref, after); err != nil {
			return "", 0, nil, err
		}
		// Keep containment valid by construction (DESIGN §5): rather than let a
		// document land under an inline parent and be rejected as E002, promote the
		// destination parent chain — the same mkdir -p cascade `promote` performs.
		if carriesDocument {
			p, err := promote(r, tree.Build(r), opts.NewParent)
			if err != nil {
				return "", 0, nil, err
			}
			promoted = p
		}
	} else {
		if err := renameInParent(r, oldParent, seg, oldKey, newKey); err != nil {
			return "", 0, nil, err
		}
		if ref != "" {
			if err := reorderInParent(r, oldParent, seg, newKey, ref, after); err != nil {
				return "", 0, nil, err
			}
		}
	}
	return newAddr, n, promoted, nil
}

// subtreeRendersDocument reports whether addr or any element beneath it renders
// as its own document — the condition under which re-parenting must ensure a
// document-hosting parent.
func subtreeRendersDocument(ix *tree.Index, addr string) bool {
	for _, e := range ix.All() {
		if (e.Addr == addr || strings.HasPrefix(e.Addr, addr+".")) && e.RendersAsDocument() {
			return true
		}
	}
	return false
}

// canNest reports whether an element of kind may be re-parented under a parent of
// parentKind. Only the multi-parent kinds qualify; the rest have a single legal
// home, so re-parenting them is meaningless (and would cross a domain boundary).
func canNest(kind, parentKind tree.Kind) bool {
	switch kind {
	case tree.KindOutcome:
		return parentKind == tree.KindJob
	case tree.KindRisk, tree.KindRequirement:
		return parentKind == tree.KindOutcome
	}
	return false
}

// reorderInParent moves key adjacent to ref within the parent's seg collection.
func reorderInParent(r *model.Root, parentAddr, seg, key, ref string, after bool) error {
	if parentAddr == "" && seg == "change_records" {
		return r.ChangeRecords.Reorder(key, ref, after)
	}
	parent, save, err := locate(r, parentAddr)
	if err != nil {
		return err
	}
	if err := reorderChild(parent, seg, key, ref, after); err != nil {
		return err
	}
	save()
	return nil
}

// renameInParent renames a key in place within the parent's seg collection,
// preserving its position.
func renameInParent(r *model.Root, parentAddr, seg, oldKey, newKey string) error {
	if parentAddr == "" && seg == "change_records" {
		if !r.ChangeRecords.Rename(oldKey, newKey) {
			return fmt.Errorf("cannot rename %q to %q", oldKey, newKey)
		}
		return nil
	}
	parent, save, err := locate(r, parentAddr)
	if err != nil {
		return err
	}
	if err := renameChild(parent, seg, oldKey, newKey); err != nil {
		return err
	}
	save()
	return nil
}

func reorderChild(parent any, seg, key, ref string, after bool) error {
	switch p := parent.(type) {
	case *model.Product:
		switch seg {
		case "jobs":
			return p.Jobs.Reorder(key, ref, after)
		case "decision_records":
			return p.DecisionRecords.Reorder(key, ref, after)
		}
	case *model.Job:
		if seg == "outcomes" {
			return p.Outcomes.Reorder(key, ref, after)
		}
	case *model.Outcome:
		switch seg {
		case "risks":
			return p.Risks.Reorder(key, ref, after)
		case "requirements":
			return p.Requirements.Reorder(key, ref, after)
		}
	case *model.Engineering:
		switch seg {
		case "components":
			return p.Components.Reorder(key, ref, after)
		case "principles":
			return p.Principles.Reorder(key, ref, after)
		case "constraints":
			return p.Constraints.Reorder(key, ref, after)
		case "decision_records":
			return p.DecisionRecords.Reorder(key, ref, after)
		}
	}
	return fmt.Errorf("cannot reorder within %q", seg)
}

func renameChild(parent any, seg, oldKey, newKey string) error {
	renamed := false
	switch p := parent.(type) {
	case *model.Product:
		switch seg {
		case "jobs":
			renamed = p.Jobs.Rename(oldKey, newKey)
		case "decision_records":
			renamed = p.DecisionRecords.Rename(oldKey, newKey)
		}
	case *model.Job:
		if seg == "outcomes" {
			renamed = p.Outcomes.Rename(oldKey, newKey)
		}
	case *model.Outcome:
		switch seg {
		case "risks":
			renamed = p.Risks.Rename(oldKey, newKey)
		case "requirements":
			renamed = p.Requirements.Rename(oldKey, newKey)
		}
	case *model.Engineering:
		switch seg {
		case "components":
			renamed = p.Components.Rename(oldKey, newKey)
		case "principles":
			renamed = p.Principles.Rename(oldKey, newKey)
		case "constraints":
			renamed = p.Constraints.Rename(oldKey, newKey)
		case "decision_records":
			renamed = p.DecisionRecords.Rename(oldKey, newKey)
		}
	}
	if !renamed {
		return fmt.Errorf("cannot rename %q to %q", oldKey, newKey)
	}
	return nil
}

// relocate extracts the element from its source parent and inserts it under the
// destination parent. It anchors at the live product root and writes the whole
// chain back after the extract before the insert, so a move within a shared
// ancestor stays consistent.
func relocate(r *model.Root, kind tree.Kind, srcParent, oldKey, dstParent, newKey, ref string, after bool) error {
	switch kind {
	case tree.KindOutcome:
		return relocateOutcome(r, srcParent, oldKey, dstParent, newKey, ref, after)
	case tree.KindRisk, tree.KindRequirement:
		return relocateLeaf(r, kind, srcParent, oldKey, dstParent, newKey, ref, after)
	}
	return fmt.Errorf("cannot re-parent a %s", kind)
}

func relocateOutcome(r *model.Root, srcJobAddr, oldKey, dstJobAddr, newKey, ref string, after bool) error {
	srcJK := tree.LastSegment(srcJobAddr)
	dstJK := tree.LastSegment(dstJobAddr)

	srcJob, ok := r.Product.Jobs.Get(srcJK)
	if !ok {
		return notFound(srcJobAddr)
	}
	val, ok := srcJob.Outcomes.Get(oldKey)
	if !ok {
		return notFound(srcJobAddr + ".outcomes." + oldKey)
	}
	srcJob.Outcomes.Delete(oldKey)
	r.Product.Jobs.Set(srcJK, srcJob)

	dstJob, ok := r.Product.Jobs.Get(dstJK)
	if !ok {
		return notFound(dstJobAddr)
	}
	dstJob.Outcomes.Set(newKey, val)
	if ref != "" {
		if err := dstJob.Outcomes.Reorder(newKey, ref, after); err != nil {
			return err
		}
	}
	r.Product.Jobs.Set(dstJK, dstJob)
	return nil
}

func relocateLeaf(r *model.Root, kind tree.Kind, srcOcAddr, oldKey, dstOcAddr, newKey, ref string, after bool) error {
	srcJK, srcOK := jobAndOutcome(srcOcAddr)
	dstJK, dstOK := jobAndOutcome(dstOcAddr)

	srcJob, ok := r.Product.Jobs.Get(srcJK)
	if !ok {
		return notFound(srcOcAddr)
	}
	srcOc, ok := srcJob.Outcomes.Get(srcOK)
	if !ok {
		return notFound(srcOcAddr)
	}

	var risk model.Risk
	var req model.Requirement
	switch kind {
	case tree.KindRisk:
		if risk, ok = srcOc.Risks.Get(oldKey); !ok {
			return notFound(srcOcAddr + ".risks." + oldKey)
		}
		srcOc.Risks.Delete(oldKey)
	case tree.KindRequirement:
		if req, ok = srcOc.Requirements.Get(oldKey); !ok {
			return notFound(srcOcAddr + ".requirements." + oldKey)
		}
		srcOc.Requirements.Delete(oldKey)
	}
	srcJob.Outcomes.Set(srcOK, srcOc)
	r.Product.Jobs.Set(srcJK, srcJob)

	dstJob, ok := r.Product.Jobs.Get(dstJK)
	if !ok {
		return notFound(dstOcAddr)
	}
	dstOc, ok := dstJob.Outcomes.Get(dstOK)
	if !ok {
		return notFound(dstOcAddr)
	}
	switch kind {
	case tree.KindRisk:
		dstOc.Risks.Set(newKey, risk)
		if ref != "" {
			if err := dstOc.Risks.Reorder(newKey, ref, after); err != nil {
				return err
			}
		}
	case tree.KindRequirement:
		dstOc.Requirements.Set(newKey, req)
		if ref != "" {
			if err := dstOc.Requirements.Reorder(newKey, ref, after); err != nil {
				return err
			}
		}
	}
	dstJob.Outcomes.Set(dstOK, dstOc)
	r.Product.Jobs.Set(dstJK, dstJob)
	return nil
}

// jobAndOutcome splits an outcome address into its job key and outcome key. The
// caller has validated the address is an outcome, so the shape is fixed.
func jobAndOutcome(outcomeAddr string) (string, string) {
	segs := strings.Split(outcomeAddr, ".")
	return segs[2], segs[4]
}

func replaceLast(addr, newKey string) string {
	if i := strings.LastIndex(addr, "."); i >= 0 {
		return addr[:i+1] + newKey
	}
	return newKey
}

// --- reference cascade ---------------------------------------------------

// remapRefs rewrites every reference in the tree from the old address (or an
// address under it) to the new one, and returns how many it changed. It walks
// the model exactly as tree.Build does, so no reference-bearing field is missed:
// declared edge lists, relationship map keys, and every prose field (which may
// carry {{ }} addresses). It runs before the structural move, so each element's
// owner address is still its current one — which prose rewriting needs to
// resolve relative addresses.
func remapRefs(r *model.Root, oldAddr, newAddr string) int {
	remap := makeRemap(oldAddr, newAddr)
	n := 0

	rewriteProse(&r.Product.Summary, "product", remap, &n)
	rewriteProse(&r.Product.Detail, "product", remap, &n)
	for _, jk := range r.Product.Jobs.Keys() {
		job, _ := r.Product.Jobs.Get(jk)
		jAddr := "product.jobs." + jk
		rewriteProse(&job.Story, jAddr, remap, &n)
		rewriteProse(&job.Detail, jAddr, remap, &n)
		for _, ok := range job.Outcomes.Keys() {
			oc, _ := job.Outcomes.Get(ok)
			oAddr := jAddr + ".outcomes." + ok
			rewriteProse(&oc.Statement, oAddr, remap, &n)
			rewriteProse(&oc.Detail, oAddr, remap, &n)
			for _, rk := range oc.Risks.Keys() {
				risk, _ := oc.Risks.Get(rk)
				rAddr := oAddr + ".risks." + rk
				rewriteProse(&risk.Statement, rAddr, remap, &n)
				rewriteProse(&risk.Detail, rAddr, remap, &n)
				oc.Risks.Set(rk, risk)
			}
			for _, qk := range oc.Requirements.Keys() {
				req, _ := oc.Requirements.Get(qk)
				qAddr := oAddr + ".requirements." + qk
				rewriteProse(&req.Statement, qAddr, remap, &n)
				rewriteProse(&req.Detail, qAddr, remap, &n)
				rewriteProseList(req.AcceptanceCriteria, qAddr, remap, &n)
				remapEdges(req.Mitigates, remap, &n)
				remapEdges(req.DependsOn, remap, &n)
				oc.Requirements.Set(qk, req)
			}
			job.Outcomes.Set(ok, oc)
		}
		r.Product.Jobs.Set(jk, job)
	}
	for _, dk := range r.Product.DecisionRecords.Keys() {
		dr, _ := r.Product.DecisionRecords.Get(dk)
		remapDecisionRecord(&dr, "product.decision_records."+dk, remap, &n)
		r.Product.DecisionRecords.Set(dk, dr)
	}

	rewriteProse(&r.Engineering.Summary, "engineering", remap, &n)
	rewriteProse(&r.Engineering.Detail, "engineering", remap, &n)
	for _, pk := range r.Engineering.Principles.Keys() {
		v, _ := r.Engineering.Principles.Get(pk)
		rewriteProse(&v.Statement, "engineering.principles."+pk, remap, &n)
		r.Engineering.Principles.Set(pk, v)
	}
	for _, ck := range r.Engineering.Constraints.Keys() {
		v, _ := r.Engineering.Constraints.Get(ck)
		rewriteProse(&v.Statement, "engineering.constraints."+ck, remap, &n)
		r.Engineering.Constraints.Set(ck, v)
	}
	for _, ck := range r.Engineering.Components.Keys() {
		c, _ := r.Engineering.Components.Get(ck)
		cAddr := "engineering.components." + ck
		rewriteProse(&c.Responsibility, cAddr, remap, &n)
		rewriteProse(&c.Detail, cAddr, remap, &n)
		rewriteProse(&c.DataModel, cAddr, remap, &n)
		rewriteProse(&c.Interfaces, cAddr, remap, &n)
		rewriteProse(&c.Behavior, cAddr, remap, &n)
		rewriteProseList(c.EdgeCases, cAddr, remap, &n)
		rewriteProseList(c.SuccessCriteria, cAddr, remap, &n)
		remapEdges(c.Fulfills, remap, &n)
		remapRelationships(&c, cAddr, remap, &n)
		r.Engineering.Components.Set(ck, c)
	}
	for _, dk := range r.Engineering.DecisionRecords.Keys() {
		dr, _ := r.Engineering.DecisionRecords.Get(dk)
		remapDecisionRecord(&dr, "engineering.decision_records."+dk, remap, &n)
		r.Engineering.DecisionRecords.Set(dk, dr)
	}

	for _, ck := range r.ChangeRecords.Keys() {
		cr, _ := r.ChangeRecords.Get(ck)
		cAddr := "change_records." + ck
		remapEdges(cr.Affects, remap, &n)
		rewriteProse(&cr.Change, cAddr, remap, &n)
		rewriteProse(&cr.Rationale, cAddr, remap, &n)
		r.ChangeRecords.Set(ck, cr)
	}
	return n
}

func remapDecisionRecord(dr *model.DecisionRecord, addr string, remap func(string) string, n *int) {
	remapEdges(dr.Affects, remap, n)
	rewriteProse(&dr.Summary, addr, remap, n)
	rewriteProse(&dr.Context, addr, remap, n)
	rewriteProse(&dr.Decision, addr, remap, n)
	rewriteProseList(dr.Options, addr, remap, n)
	rewriteProseList(dr.Consequences, addr, remap, n)
}

func remapRelationships(c *model.Component, owner string, remap func(string) string, n *int) {
	if c.Relationships.Len() == 0 {
		return
	}
	var next model.OrderedMap[string]
	for _, k := range c.Relationships.Keys() {
		note, _ := c.Relationships.Get(k)
		nk := remap(k)
		if nk != k {
			*n++
		}
		rewriteProse(&note, owner, remap, n)
		next.Set(nk, note)
	}
	c.Relationships = next
}

func makeRemap(oldAddr, newAddr string) func(string) string {
	prefix := oldAddr + "."
	return func(a string) string {
		if a == oldAddr {
			return newAddr
		}
		if strings.HasPrefix(a, prefix) {
			return newAddr + a[len(oldAddr):]
		}
		return a
	}
}

func rewriteProse(s *string, owner string, remap func(string) string, n *int) {
	if *s == "" {
		return
	}
	if out := interp.Rewrite(*s, owner, remap(owner), remap); out != *s {
		*s = out
		*n++
	}
}

func rewriteProseList(list []string, owner string, remap func(string) string, n *int) {
	for i := range list {
		rewriteProse(&list[i], owner, remap, n)
	}
}

func remapEdges(list []string, remap func(string) string, n *int) {
	for i := range list {
		if v := remap(list[i]); v != list[i] {
			list[i] = v
			*n++
		}
	}
}
