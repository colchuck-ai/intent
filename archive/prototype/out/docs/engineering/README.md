# Intent

Intent is delivered as a set of files an AI assistant loads at runtime — a skill definition, reference guides, document templates, and a validation script — that together teach the assistant how to place, draft, and check product and engineering documentation.

## Principles

- Default every element to the lowest tier (an inline statement) and promote only when a concrete need forces it — never scaffold structure ahead of need.
- Keep documents (current state) and records (history and rationale) in separate files; a document and its records never share a file.
- Make every element resolvable by a stable logical ID before any file exists for it.

## Constraints

- Records live only in the four central directories (`docs/product/crs/`, `docs/product/drs/`, `docs/engineering/crs/`, `docs/engineering/drs/`) — never inline, never in a per-element folder.
- Slugs carry only the local ID segment; ancestor context comes from the directory path, not the slug.

## Technology Choices

- Markdown files under `docs/`: human- and agent-readable, diffable in git, no runtime or database needed to operate.
- Python standard library (PEP 723, no external deps) for the linter script: runs anywhere Python 3 runs with no install step, so it can run as the mechanical first pass before human judgment in the Judge loop.

## Components

### C001 - Skill Definition

Declares Intent's activation triggers, principles, gotchas, and the reference index an assistant loads to operate the framework.

**Relationships**

- **C002 - Reference Guides**: points into the reference guides for full rules rather than restating them.
- **C004 - Bundle Linter**: lists the linter under "Available scripts" as the mechanical first pass for the Judge loop.

### C002 - Reference Guides

Defines what each element is, where it lives, how to modify it, and what context to trace, across `elements.md`, `structure.md`, `workflows.md`, and `context-tracing.md`.

**Relationships**

- **C003 - Templates**: the paths and tiers defined here are what the templates instantiate.

### C003 - Templates

Provides one canonical per-element × tier document skeleton, fixing required vs. omittable sections and the inline-vs-reference child rendering form.

**Relationships**

- **C002 - Reference Guides**: each template implements the paths and tier rules defined there.

**Architectural Decision Records**

- [ADR001 - Separate Tier Templates](drs/ADR001-separate-tier-templates.md)

### C004 - Bundle Linter

Runs a mechanical first pass over a `docs/` tree — logical-ID format/uniqueness/no-reuse, backlink bidirectionality, empty headings, requirements with no linked risk, and CR tree-side placement — before human judgment is applied in the Judge loop.

**Relationships**

- **C002 - Reference Guides**: encodes the rules defined there as automated checks.

## Requirement-Component Map

- **O001-R001 - Stable Logical IDs**: C001 - Skill Definition, C002 - Reference Guides
- **O001-R002 - One Tracing Path Per Element**: C002 - Reference Guides
- **O002-R001 - Alignment Check Before Finalizing**: C002 - Reference Guides
- **O002-R002 - Capture Change Records**: C002 - Reference Guides, C003 - Templates
- **O002-R003 - Mechanical Linter First Pass**: C004 - Bundle Linter
- **O003-R001 - Minimal-First Tier Ladder**: C002 - Reference Guides
- **O003-R002 - Canonical Path/Slug Convention**: C002 - Reference Guides, C003 - Templates
