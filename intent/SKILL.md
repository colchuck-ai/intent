# Intent

## Elements

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

A requirement is what the product must do or respect to mitigate a risk.

A requirement is usually:

- Risk-linked: tied to a specific risk.
- Solution-aware: this is where the product enters the picture — what it must do, not how it's built.
- Verifiable: you can test whether it's met or not.
- Necessary: it exists to mitigate a risk, not for its own sake.

Minimal Pattern:

> [Product/Solution] must [Capability/Constraint]

- Product/Solution: the product or solution being constrained.
- Capability/Constraint: what it must do or respect.

Examples:

**Good**

- The product must surface contradictions across reviews.
- The product must allow saving and resuming a comparison.
- The product must confirm no better-matching options exist at time of decision.

Why they're good: each is risk-linked (directly references a specific risk), solution-aware but not implementation-specific (says what the product must do, not how), and verifiable (you can test whether it's met).

**Bad**

- Build a review aggregation microservice.
- Add a save button to the comparison page.

Why they're bad: they're implementation details, not requirements. They specify *how* to build, not *what* the product must do. They aren't tied to a specific risk.

### Engineering

#### Architecture

#### Component

### Records

### Change Record (CR)

### Product Decision Record (PDR)

### Architectural Decision Record (ADR)

## Files

### Product Document 

Path: `docs/product/README.md`

Template:

```md
```

### Outcome Documents

Paths: `docs/product/outcomes/O<NNN>-<name>/README.md`


Template:

```md
```

### Requirement Documents

Paths:

- `docs/product/outcomes/O<NNN>-<name>/requirements/O<NNN>-R<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/O<NNN>-R<NNN>-<name>/README.md`


Simple Template (`...-R<NNN>-<name>.md`):

```md
```


Nested Template (`...-R<NNN>-<name>/README.md`):

```md
```

### Architecture Document

Path: `docs/engineering/README.md`


Template:

```md
```

### Component Documents

Paths:

- `docs/engineering/components/C<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/README.md`


Simple Template (`...C<NNN>-<name>.md`):

```md
```


Nested Template (`...C<NNN>-<name>/README.md`):

```md
```

### Product Decision Records

Paths:

- `docs/product/pdrs/PRD<NNN>-<name>.md`
- `docs/product/pdrs/PRD<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/pdrs/O<NNN>-PRD<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/pdrs/O<NNN>-PRD<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/O<NNN>-R<NNN>-<name>/pdrs/O<NNN>-R<NNN>-PRD<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/O<NNN>-R<NNN>-<name>/pdrs/O<NNN>-R<NNN>-PRD<NNN>-<name>/README.md`

Simple Template (`...PRD<NNN>-<name>.md`):

```md
```


Nested Template (`...PRD<NNN>-<name>/README.md`):

```md
```

### Architectural Decision Records

Paths:

- `docs/engineering/adrs/ADR<NNN>-<name>.md`
- `docs/engineering/adrs/ADR<NNN>-<name>/README.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/C<NNN>-ADR<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/C<NNN>-ADR<NNN>-<name>/README.md`

Simple Template (`...-ADR<NNN>-<name>.md`):

```md
```


Nested Template (`...-ADR<NNN>-<name>/README.md`):

```md
```


### Change Records

Paths:

- `docs/product/crs/PROD-CR<NNN>-<name>.md`
- `docs/product/crs/PROD-CR<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/crs/O<NNN>-CR<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/crs/O<NNN>-CR<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/O<NNN>-R<NNN>-<name>/crs/O<NNN>-R<NNN>-CR<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/O<NNN>-R<NNN>-<name>/crs/O<NNN>-R<NNN>-CR<NNN>-<name>/README.md`
- `docs/engineering/crs/ENG-CR<NNN>-<name>.md`
- `docs/engineering/crs/ENG-CR<NNN>-<name>/README.md`
- `docs/engineering/components/C<NNN>-<name>/crs/C<NNN>-CR<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/crs/C<NNN>-CR<NNN>-<name>/README.md`

Simple Template (`...-CR<NNN>-<name>.md`):

```md
```


Nested Template (`...-CR<NNN>-<name>/README.md`):

```md
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
      O<NNN>-<name>/
        crs/
          O<NNN>-CR<NNN>-<name>/
            README.md
          O<NNN>-CR<NNN>-<name>.md
        pdrs/
          O<NNN>-PDR<NNN>-<name>/
            README.md
          O<NNN>-PDR<NNN>-<name>.md
        requirements/
          O<NNN>-R<NNN>-<name>/
            crs/
              O<NNN>-R<NNN>-CR<NNN>-<name>/
                README.md
              O<NNN>-R<NNN>-CR<NNN>-<name>.md
            pdrs/
              O<NNN>-R<NNN>-PDR<NNN>-<name>/
                README.md
              O<NNN>-R<NNN>-PDR<NNN>-<name>.md
            README.md
          O<NNN>-R<NNN>-<name>.md
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
