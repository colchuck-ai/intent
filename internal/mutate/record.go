package mutate

import (
	"fmt"
	"strings"

	"github.com/colchuck-ai/intent/internal/model"
	"github.com/colchuck-ai/intent/internal/tree"
)

// RecordFields carries the values a new record is created with. Which fields
// apply depends on the record kind: a decision record (PDR/ADR) uses Summary and
// the optional Context/Options/Decision/Consequences; a change record uses
// Change and Rationale. Affects holds fully-resolved dotted addresses (the CLI
// resolves suffixes and root literals before calling AddRecord).
type RecordFields struct {
	Name         string
	Affects      []string
	Summary      string   // decision record (PDR/ADR)
	Context      string   // decision record
	Options      []string // decision record
	Decision     string   // decision record
	Consequences []string // decision record
	Change       string   // change record (CR)
	Rationale    string   // change record
}

// AddRecord creates a record of kind (KindPDR/KindADR/KindCR) with key, populated
// from f, and returns its full address. Location implies the record type, so a
// PDR lands under product.decision_records, an ADR under
// engineering.decision_records, and a CR at top-level change_records (DESIGN §2).
//
// Required fields are checked up front for a clearer message than the schema
// gives: every record needs at least one --affects; a decision record needs a
// --summary; a change record needs --change and --rationale. Decision records
// are domain-pure — a PDR's affects must stay under product.*, an ADR's under
// engineering.* — so a cross-domain target is refused here rather than surfacing
// only as the post-write E004 (DESIGN §4).
func AddRecord(r *model.Root, kind tree.Kind, key string, f RecordFields) (string, error) {
	if !keyRe.MatchString(key) {
		return "", fmt.Errorf("key %q must be lowercase letters, digits, underscores, or hyphens, starting with a letter", key)
	}
	if len(f.Affects) == 0 {
		return "", fmt.Errorf("a record needs at least one --affects target")
	}

	switch kind {
	case tree.KindPDR:
		if err := requireDecisionFields(f); err != nil {
			return "", err
		}
		if err := checkDomainPure(f.Affects, "product", "pdr"); err != nil {
			return "", err
		}
		addr := "product.decision_records." + key
		if _, ok := r.Product.DecisionRecords.Get(key); ok {
			return "", existsErr(addr)
		}
		r.Product.DecisionRecords.Set(key, decisionRecord(f))
		return addr, nil
	case tree.KindADR:
		if err := requireDecisionFields(f); err != nil {
			return "", err
		}
		if err := checkDomainPure(f.Affects, "engineering", "adr"); err != nil {
			return "", err
		}
		addr := "engineering.decision_records." + key
		if _, ok := r.Engineering.DecisionRecords.Get(key); ok {
			return "", existsErr(addr)
		}
		r.Engineering.DecisionRecords.Set(key, decisionRecord(f))
		return addr, nil
	case tree.KindCR:
		if f.Change == "" || f.Rationale == "" {
			return "", fmt.Errorf("a change record needs --change and --rationale")
		}
		addr := "change_records." + key
		if _, ok := r.ChangeRecords.Get(key); ok {
			return "", existsErr(addr)
		}
		r.ChangeRecords.Set(key, model.ChangeRecord{
			Name: f.Name, Affects: f.Affects, Change: f.Change, Rationale: f.Rationale,
		})
		return addr, nil
	default:
		return "", fmt.Errorf("cannot create a %s record (kinds: pdr, adr, cr)", kind)
	}
}

func requireDecisionFields(f RecordFields) error {
	if f.Summary == "" {
		return fmt.Errorf("a decision record needs a --summary")
	}
	return nil
}

func decisionRecord(f RecordFields) model.DecisionRecord {
	return model.DecisionRecord{
		Name: f.Name, Affects: f.Affects, Summary: f.Summary,
		Context: f.Context, Options: f.Options, Decision: f.Decision, Consequences: f.Consequences,
	}
}

// checkDomainPure verifies every affects target stays within domain (the bare
// root literal or something under it), mirroring the E004 rule. recordType names
// the record in the error (pdr/adr).
func checkDomainPure(affects []string, domain, recordType string) error {
	for _, a := range affects {
		if a != domain && !strings.HasPrefix(a, domain+".") {
			return fmt.Errorf("a %s record's affects must stay under %s.* (got %q)", recordType, domain, a)
		}
	}
	return nil
}
