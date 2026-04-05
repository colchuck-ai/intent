---
name: intent
description: Structure and maintain product and engineering documentation using the intent framework. Use when defining jobs, outcomes, risks, requirements, architecture, components, or decision records with full traceability.
---

# Intent

A framework for structuring product and engineering intent as a traceable hierarchy — from jobs and outcomes through risks and requirements down to architecture and components — with decision and change records that keep the whole narrative coherent over time.

## Global Guidance

### Plain Language

Write so a junior PM or engineer can read and understand on the first pass.

- Use short, common words. If a simpler word means the same thing, use it.
- Keep sentences short. One idea per sentence.
- Avoid jargon, acronyms, and abstract nouns unless they are defined in this framework.
- Say what you mean directly. Don't hedge or pad.
- If you need a technical term, use it — but don't dress up a simple idea in complex language.

### Candor

Say what you know, what you don't, and what you're unsure about.

- If you lack information to make a claim, say so. Don't fill the gap with plausible-sounding guesses.
- If a decision could go either way, name the tradeoffs instead of picking one and presenting it as obvious.
- Mark assumptions as assumptions. Don't state them as facts.
- "I don't know" and "we need more information" are valid and expected answers.

## Elements

This section defines the [product](#product), [engineering](#engineering), and [record](#records) elements used throughout the intent framework (ending with [Records](#records), then [Context Tracing](#context-tracing) as the last subsection). [Files](#files) and [Structure](#structure) follow and document paths and on-disk layout for each element's documents.

### Product

A product is a solution people hire to get a job done.

It is defined not by its features but by:

- The job it serves
- The outcome customers use to measure success
- The risks that threaten those outcomes
- The requirements that constrain how it delivers

The product has no ancestors above it in this framework. Its descendants—jobs, outcomes, risks, and requirements—should read as one story with the product: nothing downstream should contradict what the product claims.

#### Job

A job is a goal someone is trying to achieve in a specific situation — the progress they want to make, not the product they use. It is functional, situational, and stable across changes to tools and features.

A job should align with the product. Outcomes under that job should express progress toward the job without pulling against the product narrative.

Sibling jobs on the same product should be well scoped: overlap should be minimal; where two jobs touch the same situation, they should agree; together they should cover the product's scope without gaps.

Minimal Pattern:

> When [Situation], I want to [Goal], so I can [Outcome].

- Situation: when, where, or why the job arises — the context or trigger that creates the need.
- Goal: what the person is trying to accomplish — the progress they want to make.
- Outcome: why it matters — the benefit or result they expect from achieving the goal.

Examples:

**Good**

> When I'm comparing options before a big purchase, I want to feel confident I'm not missing a better alternative, so I can buy without second-guessing.

Why it's good: it's functional (comparing options), situational (before a big purchase), and solution-free (no product named). The situation provides the trigger, the goal states the progress, and the outcome explains why it matters.

**Bad**

> I want a product comparison app with filters and ratings.

Why it's bad: it's a solution, not a job. Name is a product. Skips the situation and outcome entirely. 

#### Outcome

An outcome is a measurable way the customer judges success while (or after) getting a job done. It answers: "What would 'better' look like?" in terms they care about — not in terms of your product. It is stated from the customer's perspective, measurable in principle, and free of solution naming.

An outcome should align with its owning job and with the product. Risks tied to this outcome should express harms to this outcome without orphan or misaligned risks relative to it.

Sibling outcomes under the same product should be well scoped: overlap should be minimal; where peer outcomes touch the same situation, they should agree; taken together, they should cover the jobs they serve. Cousin outcomes should stay distinct in scope where possible; where they bear on the same situation, they should agree; taken together with peer outcomes, their scope should be complete without contradiction.

Minimal Pattern:

> [Verb] [Unit of Measure] [Object]

- Verb: minimize or maximize (among others)
- Unit of Measure: time, number, or likelihood
- Object: the specific thing being measured

Examples:

**Good**

- Minimize the time to identify the most relevant alternatives.
- Minimize the likelihood of overlooking a meaningful difference between options.
- Maximize the confidence that the chosen option meets the most important criteria.

Why they're good: each is customer-centric (stated from their perspective), measurable in principle (time, likelihood, confidence), and solution-free (no feature or technology named). They follow the minimal pattern with a clear verb, unit of measure, and object.

**Bad**

- Show me a side-by-side comparison table
- Give me at least 10 results

Why they're bad: they're feature requests, not measurable success metrics. They name a solution and aren't stated from the customer's perspective of what "better" looks like. 

#### Risk

A risk is a condition or event that negatively impacts a desired outcome. It is tied to that outcome, measurable by likelihood and impact on the outcome, stated from the customer's perspective, and free of product or feature naming.

A risk should belong to its outcome. Requirements that mitigate it should connect to that outcome without orphan or misaligned risks.

Sibling risks under the same outcome should avoid duplicating the same failure mode; where two risks bear on the same requirement, the narrative should agree. Cousin risks should stay distinct where possible; where mitigations span outcomes, shared facts and requirement links should agree.

Minimal Pattern:

> [Condition/Event] [Negative Impact on Outcome]

- Condition/Event: what might go wrong — a situation, trigger, or circumstance that could occur.
- Negative Impact on Outcome: how it hurts — the specific way it degrades a desired outcome.

Examples:

**Good**

- Conflicting reviews increase the time to identify the most relevant alternatives.
- Time pressure increases the likelihood of overlooking a meaningful difference between options.
- Post-purchase discovery of a better option reduces confidence that the chosen option meets the most important criteria.

Why they're good: each is outcome-linked (directly references a specific desired outcome), measurable by likelihood × impact, customer-centric (stated as conditions the customer faces), and solution-free (no product or feature named). The condition/event is distinct from the negative impact.

**Bad**

- The app might crash.
- Users might not like the UI.

Why they're bad: their product or implementation concerns are not customer-centric conditions tied to a desired outcome. They are solution-aware and not measurable by likelihood x impact on a specific outcome. 

#### Requirement

A requirement is what the product must do or respect to mitigate one or more risks. It is tied to specific risks, states what the product must do (not how it is built), is verifiable, and exists only to mitigate those risks.

A requirement should align with its outcome and the risks it mitigates, and with the product story; architecture and components should honor it without drift.

Sibling requirements under the same outcome should be well scoped: overlap should be minimal; where they overlap, they should agree. Cousin requirements should impose compatible demands when traced through shared architecture; together, requirements should cover the mitigations the outcomes need.

Minimal Pattern:

> [Product/Solution] must [Capability/Constraint]

- Product/Solution: the product or solution being constrained.
- Capability/Constraint: what it must do or respect.

Examples:

**Good**

- The product must surface contradictions across reviews.
- The product must allow saving and resuming a comparison.
- The product must confirm no better-matching options exist at time of decision.

Why they're good: each is risk-linked (directly references one or more specific risks), solution-aware but not implementation-specific (says what the product must do, not how), and verifiable (you can test whether it's met).

**Bad**

- Build a review aggregation microservice.
- Add a save button to the comparison page.

Why they're bad: they're implementation details, not requirements. They specify *how* to build, not *what* the product must do. They aren't tied to any specific risks.

### Engineering

Engineering is how the product is built. It is defined by the architecture that shapes it and the components that compose it.

Engineering should align with product intent: the architecture and components should trace to the product's requirements and outcomes; components should realize the architecture and fulfill its obligations without contradicting upstream intent.

Product-facing and engineering-facing descriptions that refer to the same requirement, component, or boundary should agree.

#### Architecture

An architecture is the set of decisions that define how components are organized, communicate, and constrain each other. It describes structure and relationships, not individual features in isolation, and is driven by product requirements.

The architecture should carry every requirement it is meant to satisfy and stay aligned with the product and outcomes above it; components should realize the structure without contradicting those requirements. Structure, principles, and technology choices within the architecture should not contradict each other.

Examples:

**Good**

> The system separates review analysis from comparison management, communicating through a shared option model. Both are stateless; session persistence is handled by the Comparison Session Store.

Why it's good: structural (describes relationships and boundaries), requirement-driven (the separation traces to distinct product requirements).

**Bad**

> We use microservices and React.

Why it's bad: names technologies, not structure. Doesn't describe how parts relate or why.

#### Component

A component is a distinct, implementable part of the system that fulfills one or more requirements. It has a clear scope, maps to those requirements, and is concrete enough to build, test, and deploy.

A component should align with the architecture and with every requirement it fulfills; behavior should match the product line above it.

Sibling components should be well scoped: overlap in responsibility should be minimal; where responsibilities meet, they should agree; together they should cover the architecture's obligations without gaps.

Minimal Pattern:

> [Name] — [Responsibility]

- Name: what the component is called.
- Responsibility: what it does.

Examples:

**Good**

- Review Analyzer — surfaces contradictions across reviews.
- Comparison Session Store — persists and restores in-progress comparisons.

Why they're good: bounded (clear scope), requirement-linked (maps to a specific product requirement), and named by what they do.

**Bad**

- Backend — handles stuff.
- Utils — shared code.

Why they're bad: unbounded scope, not linked to any requirement, names describe location or convenience rather than responsibility.

### Records

Records capture the decisions and changes made to product and engineering elements over time.

A record should tie to the product or engineering elements it concerns and remain coherent with those elements' ancestors and descendants on the trace.

Records at related scope should not contradict each other on the same fact without an explicit superseding record; where two records touch the same interface or decision surface, they should agree.

#### Change Record (CR)

A change record documents a modification to a product or engineering element — what changed, why, and what it affects. It is scoped to a specific element, explains what changed and why, and reads as a point-in-time modification.

A change record should tie to the element it names and remain coherent with that element's ancestors and descendants on the trace.

Other change records that bear on the same fact or touch the same interface should not tell conflicting stories; where they affect shared boundaries, they should agree on facts.

Examples:

**Good**

> Updated the Review Analyzer component to normalize star ratings across sources. Previously, ratings were passed through raw, causing contradictions to be missed when sources used different scales. This change affects the Review Analyzer component and the "surface contradictions across reviews" requirement.

Why it's good: scoped (names the element changed), traceable (explains what changed and why), and temporal (describes a before/after).

**Bad**

> Fixed some stuff in the backend.

Why it's bad: no scope, no rationale, no traceability to any element.

#### Product Decision Record (PDR)

A product decision record documents a product decision — the context, the options considered, and the choice made. It states the situation, alternatives, rationale, and consequences (including costs and follow-ups).

A product decision record should stay aligned with the product outcomes and requirements it affects; consequences should read through to downstream elements.

Product decision records at overlapping or related scope should not contradict without an explicit superseding record; where decisions overlap, they should agree on facts and intent.

Examples:

**Good**

> We decided to limit comparisons to three options at a time. Users comparing more than three reported decision fatigue and longer completion times. We considered unlimited comparisons and a five-option cap. Three balances thoroughness with cognitive load. Consequence: users comparing more than three options must run multiple sessions.

Why it's good: context-driven, option-aware, justified, and states the consequence.

**Bad**

> We're only showing three options.

Why it's bad: no context, no alternatives considered, no justification, no consequences.

#### Architectural Decision Record (ADR)

An architectural decision record documents an engineering decision — the context, the options considered, and the choice made. It states the situation, alternatives, rationale, and consequences (including costs and follow-ups).

An architectural decision record should fit the architecture and components it constrains; tradeoffs should trace to requirements and product intent above.

Architectural decision records at overlapping or related scope should not contradict without an explicit superseding record; where decisions overlap, they should agree on facts and intent.

Examples:

**Good**

> We decided to make the Review Analyzer stateless, storing no session data. The Comparison Session Store already handles persistence, and duplicating state would create sync risks. We considered a stateful analyzer with local caching. Consequence: the Review Analyzer must re-fetch review data on every request.

Why it's good: context-driven, option-aware, justified, and states the consequence.

**Bad**

> Review Analyzer is stateless.

Why it's bad: states the decision but not the context, alternatives, justification, or consequences.

### Context Tracing

See [references/context-tracing.md](references/context-tracing.md) for the full vertical and horizontal tracing rules per element type.

## Files

### Product Document 

Path: `docs/product/README.md`

Template: [product.md](assets/templates/product.md)

### Outcome Documents

Paths: `docs/product/outcomes/J<NNN>-O<NNN>-<name>/README.md`

Template: [outcome.md](assets/templates/outcome.md)

### Requirement Documents

Paths:

- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/README.md`


Simple Template (`...-R<NNN>-<name>.md`): [requirement-simple.md](assets/templates/requirement-simple.md)


Nested Template (`...-R<NNN>-<name>/README.md`): [requirement-nested.md](assets/templates/requirement-nested.md)

### Architecture Document

Path: `docs/engineering/README.md`


Template: [architecture.md](assets/templates/architecture.md)

### Component Documents

Paths:

- `docs/engineering/components/C<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/README.md`


Simple Template (`...C<NNN>-<name>.md`): [component-simple.md](assets/templates/component-simple.md)


Nested Template (`...C<NNN>-<name>/README.md`): [component-nested.md](assets/templates/component-nested.md)

### Product Decision Records

Paths:

- `docs/product/pdrs/PRD<NNN>-<name>.md`
- `docs/product/pdrs/PRD<NNN>-<name>/README.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/pdrs/J<NNN>-O<NNN>-PRD<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/pdrs/J<NNN>-O<NNN>-PRD<NNN>-<name>/README.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/pdrs/J<NNN>-O<NNN>-R<NNN>-PRD<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/pdrs/J<NNN>-O<NNN>-R<NNN>-PRD<NNN>-<name>/README.md`

`<PDR ID>` is the filename prefix before `-<name>` — e.g. `PRD<NNN>`, `J<NNN>-O<NNN>-PRD<NNN>`, or `J<NNN>-O<NNN>-R<NNN>-PRD<NNN>`, depending on which path above the file lives under.

Simple Template (`...<PDR ID>-<name>.md`): [pdr-simple.md](assets/templates/pdr-simple.md)


Nested Template (`...<PDR ID>-<name>/README.md`): [pdr-nested.md](assets/templates/pdr-nested.md)

### Architectural Decision Records

Paths:

- `docs/engineering/adrs/ADR<NNN>-<name>.md`
- `docs/engineering/adrs/ADR<NNN>-<name>/README.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/C<NNN>-ADR<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/C<NNN>-ADR<NNN>-<name>/README.md`

`<ADR ID>` is the filename prefix before `-<name>` — e.g. `ADR<NNN>` or `C<NNN>-ADR<NNN>`, depending on which path above the file lives under.

Simple Template (`...<ADR ID>-<name>.md`): [adr-simple.md](assets/templates/adr-simple.md)


Nested Template (`...<ADR ID>-<name>/README.md`): [adr-nested.md](assets/templates/adr-nested.md)


### Change Records

Paths:

- `docs/product/crs/PROD-CR<NNN>-<name>.md`
- `docs/product/crs/PROD-CR<NNN>-<name>/README.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/crs/J<NNN>-O<NNN>-CR<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/crs/J<NNN>-O<NNN>-CR<NNN>-<name>/README.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/crs/J<NNN>-O<NNN>-R<NNN>-CR<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/crs/J<NNN>-O<NNN>-R<NNN>-CR<NNN>-<name>/README.md`
- `docs/engineering/crs/ENG-CR<NNN>-<name>.md`
- `docs/engineering/crs/ENG-CR<NNN>-<name>/README.md`
- `docs/engineering/components/C<NNN>-<name>/crs/C<NNN>-CR<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/crs/C<NNN>-CR<NNN>-<name>/README.md`

`<CR ID>` is the filename prefix before `-<name>` — e.g. `PROD-CR<NNN>`, `J<NNN>-O<NNN>-CR<NNN>`, `J<NNN>-O<NNN>-R<NNN>-CR<NNN>`, `ENG-CR<NNN>`, or `C<NNN>-CR<NNN>`, depending on which path above the file lives under.

Simple Template (`...<CR ID>-<name>.md`): [cr-simple.md](assets/templates/cr-simple.md)


Nested Template (`...<CR ID>-<name>/README.md`): [cr-nested.md](assets/templates/cr-nested.md)

## Structure

See [references/structure.md](references/structure.md) for the full `docs/` directory tree.

## Workflows

See [references/workflows.md](references/workflows.md) for the read, create, update, and delete workflows including the three-session modification process.
