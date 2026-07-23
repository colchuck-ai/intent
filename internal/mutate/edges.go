package mutate

import (
	"fmt"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
)

// addEdge declares an outgoing edge on el. It validates that the element type
// carries that edge kind and that the edge isn't already present (a duplicate
// would be flagged E003). This is where the CLI writes the storage shape so no
// caller hand-builds an edge list (DESIGN §4, §11).
func addEdge(el any, from string, kind tree.EdgeKind, to, note string) error {
	switch v := el.(type) {
	case *model.Requirement:
		switch kind {
		case tree.EdgeMitigates:
			return appendUnique(&v.Mitigates, from, kind, to)
		case tree.EdgeDependsOn:
			return appendUnique(&v.DependsOn, from, kind, to)
		}
	case *model.Component:
		switch kind {
		case tree.EdgeFulfills:
			return appendUnique(&v.Fulfills, from, kind, to)
		case tree.EdgeRelationships:
			if _, ok := v.Relationships.Get(to); ok {
				return dupEdge(from, kind, to)
			}
			v.Relationships.Set(to, note)
			return nil
		}
	case *model.DecisionRecord:
		if kind == tree.EdgeAffects {
			return appendUnique(&v.Affects, from, kind, to)
		}
	case *model.ChangeRecord:
		if kind == tree.EdgeAffects {
			return appendUnique(&v.Affects, from, kind, to)
		}
	}
	return fmt.Errorf("%s cannot declare a %s edge", from, kind)
}

// removeEdge deletes an outgoing edge from el. It errors if the edge isn't
// declared.
func removeEdge(el any, from string, kind tree.EdgeKind, to string) error {
	switch v := el.(type) {
	case *model.Requirement:
		switch kind {
		case tree.EdgeMitigates:
			return removeFrom(&v.Mitigates, from, kind, to)
		case tree.EdgeDependsOn:
			return removeFrom(&v.DependsOn, from, kind, to)
		}
	case *model.Component:
		switch kind {
		case tree.EdgeFulfills:
			return removeFrom(&v.Fulfills, from, kind, to)
		case tree.EdgeRelationships:
			if !v.Relationships.Delete(to) {
				return noEdge(from, kind, to)
			}
			return nil
		}
	case *model.DecisionRecord:
		if kind == tree.EdgeAffects {
			return removeFrom(&v.Affects, from, kind, to)
		}
	case *model.ChangeRecord:
		if kind == tree.EdgeAffects {
			return removeFrom(&v.Affects, from, kind, to)
		}
	}
	return fmt.Errorf("%s has no %s edge to remove", from, kind)
}

func appendUnique(list *[]string, from string, kind tree.EdgeKind, to string) error {
	for _, x := range *list {
		if x == to {
			return dupEdge(from, kind, to)
		}
	}
	*list = append(*list, to)
	return nil
}

func removeFrom(list *[]string, from string, kind tree.EdgeKind, to string) error {
	for i, x := range *list {
		if x == to {
			*list = append((*list)[:i], (*list)[i+1:]...)
			return nil
		}
	}
	return noEdge(from, kind, to)
}

func dupEdge(from string, kind tree.EdgeKind, to string) error {
	return fmt.Errorf("%s already has a %s edge to %s", from, kind, to)
}

func noEdge(from string, kind tree.EdgeKind, to string) error {
	return fmt.Errorf("%s has no %s edge to %s", from, kind, to)
}
