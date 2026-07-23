package tree

import (
	"fmt"
	"sort"
	"strings"
)

// NotFoundError is returned when a query matches no element.
type NotFoundError struct{ Query string }

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("no element matches %q", e.Query)
}

// AmbiguousError is returned when a suffix query matches more than one element.
// Candidates are full addresses, sorted, so the operator can pick a longer
// suffix.
type AmbiguousError struct {
	Query      string
	Candidates []string
}

func (e *AmbiguousError) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%q is ambiguous; candidates:", e.Query)
	for _, c := range e.Candidates {
		fmt.Fprintf(&b, "\n  %s", c)
	}
	return b.String()
}

// Resolve turns a query into a single element. It accepts a full dotted address
// or the shortest unambiguous dotted-segment suffix (DESIGN §11). An exact
// address always wins; otherwise a suffix must match exactly one element, or the
// candidates are reported for disambiguation.
func (ix *Index) Resolve(query string) (*Element, error) {
	query = strings.TrimSpace(query)
	if e, ok := ix.byAddr[query]; ok {
		return e, nil
	}

	qseg := strings.Split(query, ".")
	var matches []*Element
	for _, e := range ix.order {
		if segmentSuffix(strings.Split(e.Addr, "."), qseg) {
			matches = append(matches, e)
		}
	}

	switch len(matches) {
	case 0:
		return nil, &NotFoundError{Query: query}
	case 1:
		return matches[0], nil
	default:
		cands := make([]string, len(matches))
		for i, m := range matches {
			cands[i] = m.Addr
		}
		sort.Strings(cands)
		return nil, &AmbiguousError{Query: query, Candidates: cands}
	}
}

// segmentSuffix reports whether suffix is a whole-segment tail of addr. Matching
// is on segment boundaries, so "logical_ids" does not match a segment named
// "stable_logical_ids".
func segmentSuffix(addr, suffix []string) bool {
	if len(suffix) == 0 || len(suffix) > len(addr) {
		return false
	}
	off := len(addr) - len(suffix)
	for i := range suffix {
		if addr[off+i] != suffix[i] {
			return false
		}
	}
	return true
}

// LastSegment returns the final dotted segment of an address — the element key.
func LastSegment(addr string) string {
	if i := strings.LastIndex(addr, "."); i >= 0 {
		return addr[i+1:]
	}
	return addr
}

// ParentAddr returns the address of the element that contains addr, or "" when
// addr is top-level. It drops the trailing map-key layer and the element key
// (e.g. product.jobs.j.outcomes.o -> product.jobs.j), because a nested element's
// address interleaves a collection key (jobs, outcomes, ...) with each element
// key.
func ParentAddr(addr string) string {
	segs := strings.Split(addr, ".")
	if len(segs) <= 2 {
		return ""
	}
	return strings.Join(segs[:len(segs)-2], ".")
}
