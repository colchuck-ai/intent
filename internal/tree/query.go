package tree

import "strings"

// Find returns elements whose searchable text contains query (case-insensitive),
// optionally narrowed to a kind and/or domain. An empty query matches everything,
// so `find --type component` lists all components. Results stay in tree order.
//
// Searchable text is the address, name, primary prose, and detail (DESIGN §11:
// search name/statement/story/detail; the address is included so a bare key
// works too).
func (ix *Index) Find(query string, kind Kind, domain Domain) []*Element {
	q := strings.ToLower(strings.TrimSpace(query))
	var out []*Element
	for _, e := range ix.order {
		if kind != "" && e.Kind != kind {
			continue
		}
		if domain != "" && e.Domain != domain {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(e.searchText()), q) {
			continue
		}
		out = append(out, e)
	}
	return out
}

func (e *Element) searchText() string {
	return strings.Join([]string{e.Addr, e.Name, e.Prose, e.Detail}, "\n")
}

// ProseValues returns every free-text field on the element that may carry {{ }}
// interpolation (DESIGN §6): the primary prose, detail, type-specific fields
// (text and list items), and relationship notes. Returned separately (not
// joined) so callers can attribute a match to the field it came from.
func (e *Element) ProseValues() []string {
	var out []string
	if e.Prose != "" {
		out = append(out, e.Prose)
	}
	if e.Detail != "" {
		out = append(out, e.Detail)
	}
	for _, f := range e.Fields {
		if f.Text != "" {
			out = append(out, f.Text)
		}
		out = append(out, f.List...)
	}
	for _, edge := range e.Edges {
		if edge.Note != "" {
			out = append(out, edge.Note)
		}
	}
	return out
}

// TraceNode is one node reached while tracing a neighborhood. Depth 0 is the
// start element; deeper nodes carry the edge and direction that reached them.
type TraceNode struct {
	Addr  string
	Depth int
	Via   EdgeKind // edge that reached this node (empty at depth 0)
	Up    bool     // true if reached by following an incoming ref (upstream)
	From  string   // the address we came from (empty at depth 0)
}

// Trace walks the edge graph outward from start up to maxDepth hops. down
// follows outgoing edges, up follows incoming refs; passing both walks both.
// edgeKinds, when non-empty, restricts which edge kinds are traversed. Nodes are
// returned in breadth-first order, each visited once.
func (ix *Index) Trace(start string, maxDepth int, up, down bool, edgeKinds map[EdgeKind]bool) []TraceNode {
	if maxDepth < 0 {
		maxDepth = 0
	}
	allowed := func(k EdgeKind) bool {
		return len(edgeKinds) == 0 || edgeKinds[k]
	}

	visited := map[string]bool{start: true}
	out := []TraceNode{{Addr: start, Depth: 0}}
	queue := []TraceNode{{Addr: start, Depth: 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if cur.Depth >= maxDepth {
			continue
		}

		var next []TraceNode
		if down {
			if e, ok := ix.byAddr[cur.Addr]; ok {
				for _, edge := range e.Edges {
					if allowed(edge.Kind) {
						next = append(next, TraceNode{Addr: edge.To, Via: edge.Kind, Up: false, From: cur.Addr})
					}
				}
			}
		}
		if up {
			for _, ref := range ix.Incoming(cur.Addr) {
				if allowed(ref.Kind) {
					next = append(next, TraceNode{Addr: ref.From, Via: ref.Kind, Up: true, From: cur.Addr})
				}
			}
		}

		for _, n := range next {
			if visited[n.Addr] {
				continue
			}
			visited[n.Addr] = true
			n.Depth = cur.Depth + 1
			out = append(out, n)
			queue = append(queue, n)
		}
	}
	return out
}
