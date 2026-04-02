# Intent

## Elements

This section defines the [product](#product), [engineering](#engineering), and [record](#records) elements used throughout the intent framework.

### Product

A product is a solution people hire to get a job done.

It is defined not by its features but by:

- The job it serves
- The outcome customers use to measure success
- The risks that threaten those outcomes
- The requirements that constrain how it delivers

#### Job

A job is a goal someone is trying to achieve in a specific situation - the progress they want to make, not the product they use.

A job is usually:

- Functional: what they're trying to get done (the task or result).
- Situational: when/where/why it matters (context, trigger, constraints).
- Stable: it persists even as tools and features change.

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

An outcome is a measurable way the customer judges success while (or after) getting a job done. It answers: "What would 'better' look like?" in terms they care about - not in terms of your product.

An outcome is usually:

- Customer-centric: stated from their perspective (e.g., minimize, maximize, etc.)
- Measureable in principle: you could score or compare (e.g., time, effort, errors, confidence, cost, risk) - even if some are qualitative scales in practice.
- Solution-free: it doesn't name your feature or technology; it names the result they want (e.g., "time to find X", "likelihood of making a mistake")

Minimal Pattern:

> [Verb] [Unit of Measure] [Object]

- Verb: usually minimize or maximize
- Unit of Measure: usually time, number, or likelihood
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

A risk is a condition or event that negatively impacts a desired outcome.

A risk is usually:

- Outcome-linked: tied to a specific desired outcome.
- Measureable: scored by likelihood x impact on the outcome.
- Customer-centric: stated from the customer's perspective, not the team's.
- Solution-free: names the condition or event, not a product or feature.

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

A requirement is what the product must do or respect to mitigate one or more risks.

A requirement is usually:

- Risk-linked: tied to one or more specific risks.
- Solution-aware: this is where the product enters the picture — what it must do, not how it's built.
- Verifiable: you can test whether it's met or not.
- Necessary: it exists to mitigate one or more risks, not for its own sake.

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

#### Architecture

An architecture is the set of decisions that define how components are organized, communicate, and constrain each other.

An architecture is usually:

- Structural: defines how parts relate, not what they do individually.
- Requirement-driven: shaped by product requirements.

Examples:

**Good**

> The system separates review analysis from comparison management, communicating through a shared option model. Both are stateless; session persistence is handled by the Comparison Session Store.

Why it's good: structural (describes relationships and boundaries), requirement-driven (the separation traces to distinct product requirements).

**Bad**

> We use microservices and React.

Why it's bad: names technologies, not structure. Doesn't describe how parts relate or why.

#### Component

A component is a distinct, implementable part of the system that fulfills one or more requirements.

A component is usually:

- Bounded: has a clear scope and responsibility.
- Requirement-linked: exists to fulfill one or more requirements.
- Implementable: concrete enough to build, test, and deploy.

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

#### Change Record (CR)

A change record documents a modification to a product or engineering element — what changed, why, and what it affects.

A change record is usually:

- Scoped: tied to a specific element (product, outcome, requirement, architecture, component).
- Traceable: explains what changed and why.
- Temporal: captures a point-in-time modification.

Examples:

**Good**

> Updated the Review Analyzer component to normalize star ratings across sources. Previously, ratings were passed through raw, causing contradictions to be missed when sources used different scales. This change affects the Review Analyzer component and the "surface contradictions across reviews" requirement.

Why it's good: scoped (names the element changed), traceable (explains what changed and why), and temporal (describes a before/after).

**Bad**

> Fixed some stuff in the backend.

Why it's bad: no scope, no rationale, no traceability to any element.

#### Product Decision Record (PDR)

A product decision record documents a product decision — the context, the options considered, and the choice made.

A product decision record is usually:

- Context-driven: explains the situation that prompted the decision.
- Option-aware: captures alternatives that were considered.
- Justified: states why the chosen option was selected.
- Consequential: documents the expected consequences — both positive and negative — of the decision.

Examples:

**Good**

> We decided to limit comparisons to three options at a time. Users comparing more than three reported decision fatigue and longer completion times. We considered unlimited comparisons and a five-option cap. Three balances thoroughness with cognitive load. Consequence: users comparing more than three options must run multiple sessions.

Why it's good: context-driven, option-aware, justified, and states the consequence.

**Bad**

> We're only showing three options.

Why it's bad: no context, no alternatives considered, no justification, no consequences.

#### Architectural Decision Record (ADR)

An architectural decision record documents an engineering decision — the context, the options considered, and the choice made.

An architectural decision record is usually:

- Context-driven: explains the situation that prompted the decision.
- Option-aware: captures alternatives that were considered.
- Justified: states why the chosen option was selected.
- Consequential: documents the expected consequences — both positive and negative — of the decision.

Examples:

**Good**

> We decided to make the Review Analyzer stateless, storing no session data. The Comparison Session Store already handles persistence, and duplicating state would create sync risks. We considered a stateful analyzer with local caching. Consequence: the Review Analyzer must re-fetch review data on every request.

Why it's good: context-driven, option-aware, justified, and states the consequence.

**Bad**

> Review Analyzer is stateless.

Why it's bad: states the decision but not the context, alternatives, justification, or consequences.

## Files

### Product Document 

Path: `docs/product/README.md`

Template:

```md
# <Product Name>

<What job(s) it serves and for whom.>

## Jobs

### J<NNN> - <Job Name>

> When [Situation], I want to [Goal], so I can [Outcome].

#### Outcomes

- **J<NNN>-O<NNN>** - <Outcome Name>: [Verb] [Unit of Measure] [Object]
- **J<NNN>-O<NNN>** - <Outcome Name>: [Verb] [Unit of Measure] [Object]
```

### Outcome Documents

Paths: `docs/product/outcomes/J<NNN>-O<NNN>-<name>/README.md`


Template:

```md
# J<NNN>-O<NNN> - <Outcome Name>

[Verb] [Unit of Measure] [Object]

## Risks

- **J<NNN>-O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]
- **J<NNN>-O<NNN>-RSK<NNN>** - <Risk Name>: [Condition/Event] [Negative Impact on Outcome]

## Requirements

- **J<NNN>-O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]
- **J<NNN>-O<NNN>-R<NNN>** - <Requirement Name>: [Product/Solution] must [Capability/Constraint]

## Risk-Requirement Map

- **J<NNN>-O<NNN>-RSK<NNN>**: J<NNN>-O<NNN>-R<NNN>, J<NNN>-O<NNN>-R<NNN>
- **J<NNN>-O<NNN>-RSK<NNN>**: J<NNN>-O<NNN>-R<NNN>
```

### Requirement Documents

Paths:

- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/README.md`


Simple Template (`...-R<NNN>-<name>.md`):

```md
# J<NNN>-O<NNN>-R<NNN> - <Requirement Name>

[Product/Solution] must [Capability/Constraint]

## Detail

<Expanded description of the requirement.>

## Edge Cases

- <Condition>: <Expected Behavior>
- <Condition>: <Expected Behavior>

## Examples

### <Scenario Name>

- Input: <Input>
- Expected Output: <Expected Output>
- Verification: <Verification Step(s)>

## Acceptance Criteria

- [ ] <Concrete, testable condition>
- [ ] <Concrete, testable condition>

## Dependencies

- J<NNN>-O<NNN>-R<NNN>
- J<NNN>-O<NNN>-R<NNN>
```


Nested Template (`...-R<NNN>-<name>/README.md`):

```md
# J<NNN>-O<NNN>-R<NNN> - <Requirement Name>

[Product/Solution] must [Capability/Constraint]

## Detail

<Expanded description of the requirement.>

## Edge Cases

- <Condition>: <Expected Behavior>
- <Condition>: <Expected Behavior>

## Examples

### <Scenario Name>

- Input: <Input>
- Expected Output: <Expected Output>
- Verification: <Verification Step(s)>

## Acceptance Criteria

- [ ] <Concrete, testable condition>
- [ ] <Concrete, testable condition>

## Dependencies

- J<NNN>-O<NNN>-R<NNN>
- J<NNN>-O<NNN>-R<NNN>

## See Also

### Product Decision Records

- [<PDR Name>](pdrs/<PDR filename>)

### Change Records

- [<CR Name>](crs/<CR filename>)

### Other Materials

- [<Material Name>](<path>)
```

### Architecture Document

Path: `docs/engineering/README.md`


Template:

```md
# <Architecture Name>

<How components are organized, communicate, and constrain each other.>

## Principles

- <Principle>
- <Principle>

## Constraints

- <Constraint>
- <Constraint>

## Technology Choices

- <Technology>: <Rationale>
- <Technology>: <Rationale>

## Components

### C<NNN> - <Component Name>

<Responsibility>

#### Relationships

- **C<NNN>**: <How they communicate or depend>
- **C<NNN>**: <How they communicate or depend>

### C<NNN> - <Component Name>

<Responsibility>

#### Relationships

- **C<NNN>**: <How they communicate or depend>

## Requirement-Component Map

- **J<NNN>-O<NNN>-R<NNN>**: C<NNN>, C<NNN>
- **J<NNN>-O<NNN>-R<NNN>**: C<NNN>
```

### Component Documents

Paths:

- `docs/engineering/components/C<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/README.md`


Simple Template (`...C<NNN>-<name>.md`):

```md
# C<NNN> - <Component Name>

<Responsibility>

## Data model

<Owned or exposed entities, schemas, invariants. Omit section if trivial.>

## Interfaces

<APIs, events, protocols, or formats — inputs, outputs, versioning. Omit section if trivial.>

## Behavior

<Non-obvious logic: ordering, retries, idempotency, conflict handling, or core algorithms. Omit section if straightforward.>

## Edge cases

- <Condition>: <Expected handling>
- <Condition>: <Expected handling>

## Relationships

- **C<NNN>**: <How this component communicates or depends on the other>
- **C<NNN>**: <How this component communicates or depends on the other>

## Success criteria

- <Engineering-level check — e.g. SLO, invariant, load or correctness assumption>
- <Engineering-level check>

## Notes

<Optional. Constraints, runbooks, or links not covered above.>
```


Nested Template (`...C<NNN>-<name>/README.md`):

```md
# C<NNN> - <Component Name>

<Responsibility>

<Optional paragraph: scope, boundaries, or context when the one-liner is not enough.>

## Data model

<Owned or exposed entities, schemas, invariants. Link to supplementary files in this folder if the model is large. Omit section if trivial.>

## Interfaces

<APIs, events, protocols, or formats — inputs, outputs, versioning. Link to OpenAPI specs, event catalogs, or similar in this folder if needed. Omit section if trivial.>

## Behavior

<Non-obvious logic: ordering, retries, idempotency, conflict handling, or core algorithms. Omit section if straightforward.>

## Edge cases

- <Condition>: <Expected handling>
- <Condition>: <Expected handling>

## Relationships

- **C<NNN>**: <How this component communicates or depends on the other>
- **C<NNN>**: <How this component communicates or depends on the other>

## Success criteria

- <Engineering-level check — e.g. SLO, invariant, load or correctness assumption>
- <Engineering-level check>

## Notes

<Optional. Constraints, runbooks, or operational detail not covered above.>

## See Also

### Architectural Decision Records

- [<ADR Name>](adrs/<ADR filename>)

### Change Records

- [<CR Name>](crs/<CR filename>)

### Other Materials

- [<Material Name>](<path>)
```

### Product Decision Records

Paths:

- `docs/product/pdrs/PRD<NNN>-<name>.md`
- `docs/product/pdrs/PRD<NNN>-<name>/README.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/pdrs/J<NNN>-O<NNN>-PRD<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/pdrs/J<NNN>-O<NNN>-PRD<NNN>-<name>/README.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/pdrs/J<NNN>-O<NNN>-R<NNN>-PRD<NNN>-<name>.md`
- `docs/product/outcomes/J<NNN>-O<NNN>-<name>/requirements/J<NNN>-O<NNN>-R<NNN>-<name>/pdrs/J<NNN>-O<NNN>-R<NNN>-PRD<NNN>-<name>/README.md`

`<PDR ID>` is the filename prefix before `-<name>` — e.g. `PRD<NNN>`, `J<NNN>-O<NNN>-PRD<NNN>`, or `J<NNN>-O<NNN>-R<NNN>-PRD<NNN>`, depending on which path above the file lives under.

Simple Template (`...<PDR ID>-<name>.md`):

```md
# <PDR ID> - <Decision Name>

<What this decision is about — a concise summary.>

## Context

<The situation, constraints, or triggers that prompted the decision.>

## Options

- <Option A>: <tradeoffs or notes>
- <Option B>: <tradeoffs or notes>

## Decision

<What was chosen and why — the justification.>

## Consequences

- <Positive or enabling consequence>
- <Negative, cost, or follow-up consequence>
```


Nested Template (`...<PDR ID>-<name>/README.md`):

```md
# <PDR ID> - <Decision Name>

<What this decision is about — a concise summary.>

<Optional paragraph: supporting detail, stakeholder input, or links to files in this folder (research, mockups, data).>

## Context

<The situation, constraints, or triggers that prompted the decision.>

## Options

- <Option A>: <tradeoffs or notes>
- <Option B>: <tradeoffs or notes>

## Decision

<What was chosen and why — the justification.>

## Consequences

- <Positive or enabling consequence>
- <Negative, cost, or follow-up consequence>

## See Also

### Change Records

- [<CR Name>](<path>)

### Related product elements

- [<Outcome, job, requirement, or other parent doc>](<path>)

### Other Materials

- [<Material Name>](<path>)
```

### Architectural Decision Records

Paths:

- `docs/engineering/adrs/ADR<NNN>-<name>.md`
- `docs/engineering/adrs/ADR<NNN>-<name>/README.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/C<NNN>-ADR<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/C<NNN>-ADR<NNN>-<name>/README.md`

`<ADR ID>` is the filename prefix before `-<name>` — e.g. `ADR<NNN>` or `C<NNN>-ADR<NNN>`, depending on which path above the file lives under.

Simple Template (`...<ADR ID>-<name>.md`):

```md
# <ADR ID> - <Decision Name>

<What this engineering decision is about — a concise summary.>

## Context

<The situation, constraints, or triggers that prompted the decision.>

## Options

- <Option A>: <tradeoffs or notes>
- <Option B>: <tradeoffs or notes>

## Decision

<What was chosen and why — the justification.>

## Consequences

- <Positive or enabling consequence>
- <Negative, cost, or follow-up consequence>
```


Nested Template (`...<ADR ID>-<name>/README.md`):

```md
# <ADR ID> - <Decision Name>

<What this engineering decision is about — a concise summary.>

<Optional paragraph: supporting detail, benchmarks, spikes, or links to files in this folder.>

## Context

<The situation, constraints, or triggers that prompted the decision.>

## Options

- <Option A>: <tradeoffs or notes>
- <Option B>: <tradeoffs or notes>

## Decision

<What was chosen and why — the justification.>

## Consequences

- <Positive or enabling consequence>
- <Negative, cost, or follow-up consequence>

## See Also

### Change Records

- [<CR Name>](<path>)

### Related engineering elements

- [<Architecture, component, or other doc>](<path>)

### Other Materials

- [<Material Name>](<path>)
```


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

Simple Template (`...<CR ID>-<name>.md`):

```md
# <CR ID> - <Change Name>

<What changed — name the scoped element(s) (e.g. outcome, requirement, architecture, component) and summarize the modification.>

## Change

<Before and after, or fuller narrative.>

## Rationale

<Why the change was made.>

## Affects

<What else this modification impacts — elements, requirements, or follow-up work.>
```


Nested Template (`...<CR ID>-<name>/README.md`):

```md
# <CR ID> - <Change Name>

<What changed — name the scoped element(s) (e.g. outcome, requirement, architecture, component) and summarize the modification.>

<Optional paragraph: diffs, migration notes, or links to files in this folder.>

## Change

<Before and after, or fuller narrative.>

## Rationale

<Why the change was made.>

## Affects

<What else this modification impacts — elements, requirements, or follow-up work.>

## See Also

### Decision Records

- [<PDR or ADR title>](<path>)

### Related elements

- [<Outcome, requirement, component, architecture, or other doc>](<path>)

### Other Materials

- [<Material Name>](<path>)
```

## Structure

```txt
docs/
  product/
    crs/
      PROD-CR<NNN>-<name>/
        README.md
      PROD-CR<NNN>-<name>.md
    outcomes/
      J<NNN>-O<NNN>-<name>/
        crs/
          J<NNN>-O<NNN>-CR<NNN>-<name>/
            README.md
          J<NNN>-O<NNN>-CR<NNN>-<name>.md
        pdrs/
          J<NNN>-O<NNN>-PDR<NNN>-<name>/
            README.md
          J<NNN>-O<NNN>-PDR<NNN>-<name>.md
        requirements/
          J<NNN>-O<NNN>-R<NNN>-<name>/
            crs/
              J<NNN>-O<NNN>-R<NNN>-CR<NNN>-<name>/
                README.md
              J<NNN>-O<NNN>-R<NNN>-CR<NNN>-<name>.md
            pdrs/
              J<NNN>-O<NNN>-R<NNN>-PDR<NNN>-<name>/
                README.md
              J<NNN>-O<NNN>-R<NNN>-PDR<NNN>-<name>.md
            README.md
          J<NNN>-O<NNN>-R<NNN>-<name>.md
        README.md
    pdrs/
      PDR<NNN>-<name>/
        README.md
      PDR<NNN>-<name>.md
    README.md
  engineering/
    adrs/
      ADR<NNN>-<name>/
        README.md
      ADR<NNN>-<name>.md
    components/
      C<NNN>-<name>/
        adrs/
          C<NNN>-ADR<NNN>-<name>/
            README.md
          C<NNN>-ADR<NNN>-<name>.md
        crs/
          C<NNN>-CR<NNN>-<name>/
            README.md
          C<NNN>-CR<NNN>-<name>.md
        README.md
      C<NNN>-<name>.md
    crs/
      ENG-CR<NNN>-<name>/
        README.md
      ENG-CR<NNN>-<name>.md
    README.md
```

## Templates

- product readme
- outcome readme
- requirement readme
- engineering readme
- component readme
- product decision record
- architectural decision record
- change record

## Verbs

- create
- read
- update
- delete

## Phases of a modification
 
1. draft
2. judge
3. review

## Resources

- product
- outcome
- requirement
- architecture
- component

## Review Criteria for Outcomes, Requirements, and Components

- scope is well defined
- overlap is minimized
- when overlapping, elements agree
- their union is complete
- in agreement with their parent (outcomes -> jobs, requirements -> outcome, components -> requirements)
