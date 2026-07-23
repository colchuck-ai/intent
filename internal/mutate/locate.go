package mutate

import (
	"fmt"
	"strings"

	"github.com/colchuck-ai/intent/internal/model"
)

// locate resolves a full dotted address to a pointer to its element and a save
// closure. Roots return a pointer straight into the model (save is a no-op);
// every nested map entry is a copy, and its save closure writes the copy — and
// every ancestor copy above it — back into the model. Mutate the returned
// pointer, then call save.
//
// The returned pointer is one of: *model.Product, *model.Engineering,
// *model.Job, *model.Outcome, *model.Risk, *model.Requirement, *model.Component,
// *model.DecisionRecord, *model.ChangeRecord, or *string (a principle or
// constraint, whose value is its statement). Callers type-switch on it.
func locate(r *model.Root, addr string) (any, func(), error) {
	segs := strings.Split(addr, ".")
	noop := func() {}

	switch segs[0] {
	case "product":
		if len(segs) == 1 {
			return &r.Product, noop, nil
		}
		switch {
		case segs[1] == "jobs" && len(segs) >= 3:
			return locateInJob(r, segs)
		case segs[1] == "decision_records" && len(segs) == 3:
			dv, ok := r.Product.DecisionRecords.Get(segs[2])
			if !ok {
				return nil, nil, notFound(addr)
			}
			return &dv, func() { r.Product.DecisionRecords.Set(segs[2], dv) }, nil
		}
	case "engineering":
		if len(segs) == 1 {
			return &r.Engineering, noop, nil
		}
		if len(segs) == 3 {
			switch segs[1] {
			case "components":
				cv, ok := r.Engineering.Components.Get(segs[2])
				if !ok {
					return nil, nil, notFound(addr)
				}
				return &cv, func() { r.Engineering.Components.Set(segs[2], cv) }, nil
			case "principles":
				pv, ok := r.Engineering.Principles.Get(segs[2])
				if !ok {
					return nil, nil, notFound(addr)
				}
				return &pv, func() { r.Engineering.Principles.Set(segs[2], pv) }, nil
			case "constraints":
				cv, ok := r.Engineering.Constraints.Get(segs[2])
				if !ok {
					return nil, nil, notFound(addr)
				}
				return &cv, func() { r.Engineering.Constraints.Set(segs[2], cv) }, nil
			case "decision_records":
				dv, ok := r.Engineering.DecisionRecords.Get(segs[2])
				if !ok {
					return nil, nil, notFound(addr)
				}
				return &dv, func() { r.Engineering.DecisionRecords.Set(segs[2], dv) }, nil
			}
		}
	case "change_records":
		if len(segs) == 2 {
			cv, ok := r.ChangeRecords.Get(segs[1])
			if !ok {
				return nil, nil, notFound(addr)
			}
			return &cv, func() { r.ChangeRecords.Set(segs[1], cv) }, nil
		}
	}
	return nil, nil, fmt.Errorf("cannot address %q", addr)
}

// locateInJob handles the product.jobs.* subtree (job → outcome → risk /
// requirement), building each level's save closure over its ancestors'.
func locateInJob(r *model.Root, segs []string) (any, func(), error) {
	addr := strings.Join(segs, ".")
	jk := segs[2]
	job, ok := r.Product.Jobs.Get(jk)
	if !ok {
		return nil, nil, notFound(addr)
	}
	saveJob := func() { r.Product.Jobs.Set(jk, job) }
	if len(segs) == 3 {
		return &job, saveJob, nil
	}

	if segs[3] != "outcomes" || len(segs) < 5 {
		return nil, nil, fmt.Errorf("cannot address %q", addr)
	}
	ok2 := segs[4]
	oc, ok := job.Outcomes.Get(ok2)
	if !ok {
		return nil, nil, notFound(addr)
	}
	saveOc := func() { job.Outcomes.Set(ok2, oc); saveJob() }
	if len(segs) == 5 {
		return &oc, saveOc, nil
	}

	if len(segs) != 7 {
		return nil, nil, fmt.Errorf("cannot address %q", addr)
	}
	leaf := segs[6]
	switch segs[5] {
	case "risks":
		rv, ok := oc.Risks.Get(leaf)
		if !ok {
			return nil, nil, notFound(addr)
		}
		return &rv, func() { oc.Risks.Set(leaf, rv); saveOc() }, nil
	case "requirements":
		qv, ok := oc.Requirements.Get(leaf)
		if !ok {
			return nil, nil, notFound(addr)
		}
		return &qv, func() { oc.Requirements.Set(leaf, qv); saveOc() }, nil
	}
	return nil, nil, fmt.Errorf("cannot address %q", addr)
}

// deleteChild removes key from the collection named seg on the parent element.
// It reports whether the key was present.
func deleteChild(parent any, seg, key string) bool {
	switch p := parent.(type) {
	case *model.Product:
		if seg == "jobs" {
			return p.Jobs.Delete(key)
		}
		if seg == "decision_records" {
			return p.DecisionRecords.Delete(key)
		}
	case *model.Job:
		if seg == "outcomes" {
			return p.Outcomes.Delete(key)
		}
	case *model.Outcome:
		switch seg {
		case "risks":
			return p.Risks.Delete(key)
		case "requirements":
			return p.Requirements.Delete(key)
		}
	case *model.Engineering:
		switch seg {
		case "components":
			return p.Components.Delete(key)
		case "principles":
			return p.Principles.Delete(key)
		case "constraints":
			return p.Constraints.Delete(key)
		case "decision_records":
			return p.DecisionRecords.Delete(key)
		}
	}
	return false
}
