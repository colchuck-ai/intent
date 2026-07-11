# PDR001 - Risk Over ODI Opportunity Scoring

Intent expresses unmet customer needs as risks (likelihood × impact) tied to an outcome, rather than adopting ODI's importance/satisfaction opportunity score.

## Context

Intent's product vocabulary otherwise follows Ulwick/ODI (jobs, outcomes, needs-first sequencing). ODI ranks outcomes by an opportunity score (importance − satisfaction) to prioritize where to invest. Intent needed a way to carry "what's unmet and why it matters" all the way down to engineering-facing requirements, so that a requirement's existence is traceable to a specific customer-facing condition rather than only to a ranked outcome.

## Options

- Adopt ODI's importance/satisfaction opportunity score on outcomes, and let requirements cite the underserved outcome directly: keeps closer fidelity to ODI, but leaves no explicit link between a requirement and the specific condition it addresses — a requirement can cite a low-scoring outcome without stating what actually threatens it.
- Add Risk as a first-class element between outcome and requirement, expressed as likelihood × impact, and require every requirement to name the risk it mitigates: adds a vocabulary term beyond ODI, but makes the outcome → requirement mitigation trace legible.

## Decision

Chose Risk as a first-class element. A requirement traces to a named risk, and a risk traces to a named outcome — a requirement with no linked risk is recognizable as a feature request, not a requirement. ODI's scoring ranks outcomes relative to each other for prioritization; Intent's risk model instead states, for each outcome, what could go wrong and how a requirement addresses it, which is the trace this framework exists to keep intact.

## Consequences

- Every requirement must name at least one risk; the bundle linter flags a requirement with no linked risk as a feature-request smell.
- Intent does not produce a ranked backlog or opportunity score the way ODI does — prioritizing across outcomes is out of scope and must come from elsewhere.
