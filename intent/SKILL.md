# Intent

## Elements

### Product

#### Job

A job is a goal someone is trying to achieve in a specific situation - the progress they want to make, not the product they use.

A job is usually:

- Functional: what they're trying to get done (the task or result).
- Situational: when/where/why it matters (context, trigger, constraints).
- Stable: it persists even as tools and features change.

Minimal Pattern:

> When [situation], I want to [goal], so I can [outcome].

Example:

> When I'm comparing options before a big purchase, I want to feel confident I'm not missing a better alternative, so I can buy without second-guessing.

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

#### Risk

#### Requirement

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
