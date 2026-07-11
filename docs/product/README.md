# Intent

Intent is an Agent Skill that teaches an AI assistant to write product and engineering documentation as one traceable tree, so every component traces back to a requirement, every requirement to a risk, and every risk to a real customer outcome. It serves PMs and engineers who use an AI assistant to draft and maintain that documentation as their product evolves.

## Jobs

### Understand The Rationale Behind An Element

> When I need to understand why a requirement, component, or decision exists before I touch it, I want to trace it back through its parent, siblings, and the customer outcome it serves, so I can be confident I'm not about to contradict something upstream.

#### O001 - Fast Rationale Lookup

Minimize the time to trace a documented component or requirement back to the customer outcome it serves, when auditing or onboarding to the docs tree.

**Risks**

- **O001-RSK001** - Tracing Slows As Tree Grows: A growing docs tree increases the time to trace a component back to the outcome it serves.
- **O001-RSK002** - Ambiguous Or Premature Reference: An element referenced before its file exists, or an ambiguous ID, increases the time spent re-tracing when a first attempt resolves to the wrong element.

**Requirements**

- **O001-R001** - Stable Logical IDs: Intent must assign every element a stable logical ID that resolves correctly whether or not a file yet exists for it.
- **O001-R002** - One Tracing Path Per Element: Intent must define one fixed vertical-and-horizontal tracing path per element type, so there is exactly one correct route to follow.

**Risk-Requirement Map**

- **O001-RSK001 - Tracing Slows As Tree Grows**: O001-R002 - One Tracing Path Per Element
- **O001-RSK002 - Ambiguous Or Premature Reference**: O001-R001 - Stable Logical IDs

### Change Documentation Without Drift

> When I'm creating, updating, or deleting product or engineering documentation — often with an AI assistant doing much of the drafting — I want the change validated against every element it touches before it's finalized, so the docs stay one coherent, traceable story instead of drifting piecemeal.

#### O002 - Coherent Change

Minimize the likelihood that a change to one element leaves a linked element — its parent, siblings, or a record about it — contradicting it.

**Risks**

- **O002-RSK001** - Linked Elements Fall Out Of Sync: Frequent AI-assisted edits increase the likelihood that a change to one element isn't reflected in the elements linked to it, leaving them in contradiction.
- **O002-RSK002** - Undocumented Change Reads As Drift: An edit made without recording what changed and why increases the likelihood that a resulting contradiction goes undetected, since there's nothing distinguishing a deliberate change from unintentional drift.

**Requirements**

- **O002-R001** - Alignment Check Before Finalizing: Intent must validate every change against its parent, descendants, siblings, and lateral links with a pass/fail alignment check before the change is finalized.
- **O002-R002** - Capture Change Records: Intent must capture a Change Record for any update or delete that is material to an outcome's success, stating what changed and why.
- **O002-R003** - Mechanical Linter First Pass: Intent must provide a mechanical linter that checks logical-ID and backlink violations as a first pass before human judgment is applied.

**Risk-Requirement Map**

- **O002-RSK001 - Linked Elements Fall Out Of Sync**: O002-R001 - Alignment Check Before Finalizing, O002-R003 - Mechanical Linter First Pass
- **O002-RSK002 - Undocumented Change Reads As Drift**: O002-R002 - Capture Change Records

#### O003 - Right-Sized Placement

Minimize the effort to determine where a new or changed piece of content belongs in the docs tree and at what level of detail.

**Risks**

- **O003-RSK001** - Wrong Tier For The Content: Uncertainty about complexity tiers increases the effort to place new content, since without a default an author must weigh every promotion trigger from scratch for each element.
- **O003-RSK002** - No Canonical Path: An undefined path or naming convention increases the effort to find or place a file correctly, since the author must invent a convention rather than look one up.

**Requirements**

- **O003-R001** - Minimal-First Tier Ladder: Intent must define a minimal-first tier ladder per element type with explicit promotion triggers, defaulting to inline until a concrete need forces promotion.
- **O003-R002** - Canonical Path/Slug Convention: Intent must define one canonical path and slug convention per element and tier so placement is deterministic.

**Risk-Requirement Map**

- **O003-RSK001 - Wrong Tier For The Content**: O003-R001 - Minimal-First Tier Ladder
- **O003-RSK002 - No Canonical Path**: O003-R002 - Canonical Path/Slug Convention

## See Also

### Product Decision Records

- [PDR001 - Risk Over ODI Opportunity Scoring](drs/PDR001-risk-over-odi-scoring.md)
