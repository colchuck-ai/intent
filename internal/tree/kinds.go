package tree

import (
	"fmt"
	"sort"
	"strings"
)

// validKinds is the closed set the --type filter accepts.
var validKinds = map[Kind]bool{
	KindProduct: true, KindEngineering: true,
	KindJob: true, KindOutcome: true, KindRisk: true, KindRequirement: true,
	KindComponent: true, KindPrinciple: true, KindConstraint: true,
	KindPDR: true, KindADR: true, KindCR: true,
}

var validDomains = map[Domain]bool{
	DomainProduct: true, DomainEngineering: true, DomainChange: true,
}

// ParseKind validates a --type value. An empty string means "no filter" and is
// allowed; anything else must be a known kind, so a typo errors instead of
// silently matching nothing.
func ParseKind(s string) (Kind, error) {
	if s == "" {
		return "", nil
	}
	k := Kind(s)
	if !validKinds[k] {
		return "", fmt.Errorf("unknown type %q; valid types: %s", s, keyList(validKinds))
	}
	return k, nil
}

// ParseDomain validates a --domain value. Empty means "no filter".
func ParseDomain(s string) (Domain, error) {
	if s == "" {
		return "", nil
	}
	d := Domain(s)
	if !validDomains[d] {
		return "", fmt.Errorf("unknown domain %q; valid domains: product, engineering, change", s)
	}
	return d, nil
}

func keyList[K ~string](m map[K]bool) string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, string(k))
	}
	sort.Strings(out)
	return strings.Join(out, ", ")
}
