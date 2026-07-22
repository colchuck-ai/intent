// Package model defines the in-memory shape of an intent.yaml file and how to
// load and canonically re-serialize it.
//
// The struct field order below is the canonical field order for output: yaml.v3
// writes struct fields in declaration order, so changing the order here changes
// every generated file. Keep it deliberate.
package model

// Root is the whole intent.yaml: three top-level roots (DESIGN §1).
type Root struct {
	Product       Product                  `yaml:"product"`
	Engineering   Engineering              `yaml:"engineering"`
	ChangeRecords OrderedMap[ChangeRecord] `yaml:"change_records,omitempty"`
}

// Product is the product root: jobs and product decision records (PDRs).
type Product struct {
	Name            string                     `yaml:"name"`
	Summary         string                     `yaml:"summary"`
	Detail          string                     `yaml:"detail,omitempty"`
	Jobs            OrderedMap[Job]            `yaml:"jobs,omitempty"`
	DecisionRecords OrderedMap[DecisionRecord] `yaml:"decision_records,omitempty"`
}

// Engineering is the engineering root: principles, constraints, components, and
// architecture decision records (ADRs).
//
// Principles and constraints are keyed maps of key->statement, not bare lists,
// so they carry stable addresses and can be edge targets (DESIGN §2, decision
// to make them "keyed but shapeless").
type Engineering struct {
	Name            string                     `yaml:"name"`
	Summary         string                     `yaml:"summary"`
	Detail          string                     `yaml:"detail,omitempty"`
	Principles      OrderedMap[string]         `yaml:"principles,omitempty"`
	Constraints     OrderedMap[string]         `yaml:"constraints,omitempty"`
	Components      OrderedMap[Component]       `yaml:"components,omitempty"`
	DecisionRecords OrderedMap[DecisionRecord] `yaml:"decision_records,omitempty"`
}

// Job is a jobs-to-be-done narrative. Jobs are always their own document, so
// they carry no `type`.
type Job struct {
	Name     string              `yaml:"name"`
	Story    string              `yaml:"story"`
	Detail   string              `yaml:"detail,omitempty"`
	Outcomes OrderedMap[Outcome] `yaml:"outcomes,omitempty"`
}

// Outcome is a measurable result the job is hired to produce. Promotable, so it
// may carry `type`.
type Outcome struct {
	Name         string                  `yaml:"name"`
	Type         string                  `yaml:"type,omitempty"`
	Statement    string                  `yaml:"statement"`
	Detail       string                  `yaml:"detail,omitempty"`
	Risks        OrderedMap[Risk]        `yaml:"risks,omitempty"`
	Requirements OrderedMap[Requirement] `yaml:"requirements,omitempty"`
}

// Risk is something that could stop an outcome. Always inline; carries no
// `type` and no edges.
type Risk struct {
	Name      string `yaml:"name"`
	Statement string `yaml:"statement"`
	Detail    string `yaml:"detail,omitempty"`
}

// Requirement is something the product must do. Promotable. `mitigates` points
// at same-outcome risks and is required (at least one); `dependsOn` points at
// requirement leaves anywhere.
type Requirement struct {
	Name               string   `yaml:"name"`
	Type               string   `yaml:"type,omitempty"`
	Statement          string   `yaml:"statement"`
	Mitigates          []string `yaml:"mitigates"`
	DependsOn          []string `yaml:"dependsOn,omitempty"`
	AcceptanceCriteria []string `yaml:"acceptance_criteria,omitempty"`
	Detail             string   `yaml:"detail,omitempty"`
}

// Component is an engineering building block. Promotable. `fulfills` points at
// the requirements it satisfies and is required (at least one);
// `relationships` is a keyed map of path->note because each edge carries a
// note payload.
type Component struct {
	Name            string             `yaml:"name"`
	Type            string             `yaml:"type,omitempty"`
	Responsibility  string             `yaml:"responsibility"`
	Fulfills        []string           `yaml:"fulfills"`
	Relationships   OrderedMap[string] `yaml:"relationships,omitempty"`
	DataModel       string             `yaml:"data_model,omitempty"`
	Interfaces      string             `yaml:"interfaces,omitempty"`
	Behavior        string             `yaml:"behavior,omitempty"`
	EdgeCases       []string           `yaml:"edge_cases,omitempty"`
	SuccessCriteria []string           `yaml:"success_criteria,omitempty"`
	Detail          string             `yaml:"detail,omitempty"`
}

// DecisionRecord is a PDR (under product) or an ADR (under engineering).
// Location implies which one it is, so there is no `kind` field. `affects` is
// domain-pure by rule (checked by the linter, not the schema).
type DecisionRecord struct {
	Name         string   `yaml:"name"`
	Affects      []string `yaml:"affects"`
	Summary      string   `yaml:"summary"`
	Context      string   `yaml:"context,omitempty"`
	Options      []string `yaml:"options,omitempty"`
	Decision     string   `yaml:"decision,omitempty"`
	Consequences []string `yaml:"consequences,omitempty"`
}

// ChangeRecord is a top-level, cross-cutting change. Its `affects` may mix
// domains and levels — it is the one record type allowed to cross domains.
type ChangeRecord struct {
	Name      string   `yaml:"name"`
	Affects   []string `yaml:"affects"`
	Change    string   `yaml:"change"`
	Rationale string   `yaml:"rationale"`
}
