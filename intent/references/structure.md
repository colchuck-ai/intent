# Structure

How elements and records map to files under `docs/`. Create a path only when the complexity tier for that element calls for it. Tier 0 elements have no dedicated path — they are recorded on a parent document listed below.

## Complexity Tiers

This framework treats every item in the hierarchy as an **element**, in two categories with different tier rules: **spine elements** (product, jobs, outcomes, risks, requirements, architecture, components) and **record elements** (CRs, PDRs, ADRs). Documents are the current state; records are the history and rationale — they never share a file.

Use the **simplest tier that captures what readers and implementers need**. Demote only when deliberately simplifying scope — never to hide unresolved ambiguity.

### Spine elements

Default to Tier 0. Fixed exceptions: product and architecture are always Tier 1; jobs and risks are always Tier 0.

| Tier | Name | On disk | Promote when |
|---|---|---|---|
| 0 | **Statement** | Minimal pattern lives on the parent document. | (default) |
| 1 | **Own document** | Single `<slug>.md` file. | Needs acceptance criteria, edge cases, cross-references that don't fit inline, or a child file. |
| 2 | **Own directory** | `<slug>/README.md` with child folders for requirements (outcomes only) or supplementary materials. | Needs child folders. |

Per-element promotion rules for outcomes, requirements, and components appear in their [Paths](#paths) sections.

### Record elements

Every CR, PDR, and ADR lives in one of four central directories — `docs/product/crs/`, `docs/product/drs/`, `docs/engineering/crs/`, `docs/engineering/drs/` — never inline, never in per-element folders, never in appendices. The scoped element backlinks via `## See Also` (Tier 1+) or a bold-label list in its inline block (Tier 0). **Capturing a record never promotes the element it concerns.**

| Tier | Name | On disk | Promote when |
|---|---|---|---|
| — | **None** | No record. | Change or decision not worth capturing. |
| 1 | **Own document** | `<slug>.md` in the central directory. | Capturing is warranted. |
| 2 | **Own directory** | `<slug>/README.md` in the central directory. | Needs supplementary files. |

### Product scale profiles

**Minimal** — Two documents (product, architecture) with everything else at Tier 0; no records. See the [minimal directory tree](#minimal-directory-tree).

**Standard** — Outcomes at Tier 1 or 2; requirements and components at Tier 0 or 1 where detail is needed. Records appear in the central directories as material decisions or changes occur.

**Complex** — Full nesting where warranted; use the [maximum directory tree](#maximum-directory-tree) as the ceiling.

## Child rendering

A parent document represents a child element in one of two modes, determined by the child's tier:

- **Tier 0 child — inline.** The child's full content lives on the parent document (minimal pattern plus any bold-labeled sub-lists).
- **Tier 1 or Tier 2 child — reference.** The parent keeps the child's heading or bullet, one line of substance (the child's minimal pattern, or the Responsibility line for components), and a `See [<Logical ID> - <Name>](path).` link to the canonical document. No paragraphs, no sub-lists, no restated detail beyond that line; everything else lives on the child's own document.

The concrete form per element type — which children get their own heading vs. render as a bullet, what bold labels to use for inlined sub-content, where record backlinks attach — is shown in the parent templates: [product.md](../assets/templates/product.md), [architecture.md](../assets/templates/architecture.md), [outcome-document.md](../assets/templates/outcome-document.md), [outcome-directory.md](../assets/templates/outcome-directory.md). Follow them as the canonical form; this section only states the rules every form must obey.

A parent doc mixes modes per child — each child renders by *its own* tier, not the parent's. The product doc therefore reads as a complete outline when every descendant is Tier 0, and stays a scannable outline (one-liner per child plus a link) as children are promoted — never collapsing into a bare table of contents.

**Form rules** (both modes):

- Promoting a child from Tier 0 to Tier 1 collapses the parent's block down to that one-liner plus the link, and moves every bold-labeled sub-list (Risks, Requirements, Relationships, record backlinks, etc.) onto the child.
- Keep the one-liner in sync with the child's canonical statement: if the minimal pattern on the child changes, update the parent's reference block in the same change.
- Never emit an empty heading. Omit any heading whose content — including every subsection beneath it — would be empty; do not leave bare placeholders. This applies to inline record sub-lists (e.g. `**Product Decision Records**`, `**Change Records**`), to `## See Also` when all of its subsections would be empty, and to any other optional section in a template. Subsections of an otherwise-populated section are also omitted individually when empty.
- Never nest a bulleted list inside another bulleted list. Use a bold-labeled group or a sub-heading instead.
- Many-to-many mapping sections (e.g. `## Risk-Requirement Map`, `## Requirement-Component Map`, and their inline counterparts) are flat bulleted lists keyed by `**<Source ID> - <Source Name>**` with comma-separated `<Target ID> - <Target Name>` pairs after the colon. Always include the human-readable name on both sides — IDs alone are unscannable.

## Naming

Logical IDs are assigned at every tier so tracing works even when there is no file. File and folder names (**slugs**) use only the local ID segment for that element type. Parent context is encoded by the directory path — do not repeat ancestor IDs in slugs.

Conventions for the `<name>` and `<NNN>` segments:

- **Name segment** (the `<name>` in slugs like `O<NNN>-<name>`): lowercase kebab-case — words separated by single hyphens, with no spaces, underscores, or capitals. Keep it short and descriptive.
- **Number** (`<NNN>`): zero-padded to three digits (e.g. `001`, `042`); once past `999`, keep the natural width (`1000`+).
- **Assignment**: use the lowest unused number within that element's scope (outcomes are numbered within the product, requirements within their owning outcome, components within the architecture, and so on).
- **No reuse**: never reuse a number after an element is deleted — retired numbers stay retired so historical references and records remain unambiguous.

Cross-tree references in document bodies and record titles use **logical IDs** that qualify from the product or engineering root so they are unambiguous outside their folder:

| Element | Slug (filename/folder) | Logical ID (in document body) |
|---|---|---|
| Job | (recorded on product doc) | `J<NNN>` |
| Outcome | `O<NNN>-<name>` | `O<NNN>` |
| Risk | (recorded inline with the outcome) | `O<NNN>-RSK<NNN>` |
| Requirement | `R<NNN>-<name>` | `O<NNN>-R<NNN>` |
| Component | `C<NNN>-<name>` | `C<NNN>` |
| PDR | `PDR<NNN>-<name>` | `PDR<NNN>` |
| ADR | `ADR<NNN>-<name>` | `ADR<NNN>` |
| CR | `CR<NNN>-<name>` | `CR<NNN>` |

Records carry no element prefix in their logical ID — the record document itself names the element(s) it scopes. PDRs and ADRs live in a single directory each (`docs/product/drs/` and `docs/engineering/drs/`), giving each one a globally unique number. CRs share a single logical-ID namespace (`CR<NNN>`); each file lives in `docs/product/crs/` or `docs/engineering/crs/` depending on which side of the tree the change touches. To resolve a bare `CR<NNN>` reference, follow the link from the element's `## See Also` (or inline `**Change Records**` list), or search both `crs/` directories.

## Paths

All paths are under `docs/`. Create a path only when the element's tier calls for it.

### Product

Tier 1 — own document

Path: `docs/product/README.md`

Templates: [product.md](../assets/templates/product.md) (Tier 1). See [Templates](#templates) for the full index.

### Outcomes

Promote to Tier 1 when risks and requirements need their own document but still fit inline; Tier 2 when they need requirement files or supplementary materials.

Tier 0 — statement

Path: none — recorded on `docs/product/README.md`

Tier 1 — own document

Path: `docs/product/outcomes/O<NNN>-<name>.md`

Tier 2 — own directory

Path: `docs/product/outcomes/O<NNN>-<name>/README.md`

Child folders: `requirements/`. Records about this outcome live in the central `docs/product/drs/` and `docs/product/crs/` directories; the outcome backlinks to them from `## See Also`.

Templates: [outcome-document.md](../assets/templates/outcome-document.md) (Tier 1), [outcome-directory.md](../assets/templates/outcome-directory.md) (Tier 2). See [Templates](#templates) for the full index.

### Requirements

A requirement at Tier 1 or above requires its owning outcome at Tier 2 — the requirement file lives inside the outcome's `requirements/` child folder. Use Tier 2 only when the requirement needs supplementary materials.

Tier 0 — statement

Path: none — recorded on the owning outcome document

Tier 1 — own document

Path: `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>.md`

Tier 2 — own directory

Path: `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>/README.md`

No record child folders — records about this requirement live in the central `docs/product/drs/` and `docs/product/crs/` directories.

Templates: [requirement-document.md](../assets/templates/requirement-document.md) (Tier 1), [requirement-directory.md](../assets/templates/requirement-directory.md) (Tier 2). See [Templates](#templates) for the full index.

### Architecture

Tier 1 — own document

Path: `docs/engineering/README.md`

Templates: [architecture.md](../assets/templates/architecture.md) (Tier 1). See [Templates](#templates) for the full index.

### Components

Promote to Tier 1 for interfaces, behavior, or success criteria; Tier 2 only for large supplementary specs.

Tier 0 — statement

Path: none — recorded on `docs/engineering/README.md`

Tier 1 — own document

Path: `docs/engineering/components/C<NNN>-<name>.md`

Tier 2 — own directory

Path: `docs/engineering/components/C<NNN>-<name>/README.md`

No record child folders — records about this component live in the central `docs/engineering/drs/` and `docs/engineering/crs/` directories.

Templates: [component-document.md](../assets/templates/component-document.md) (Tier 1), [component-directory.md](../assets/templates/component-directory.md) (Tier 2). See [Templates](#templates) for the full index.

### Product decision records

The record document names the affected element(s); each scoped element backlinks to it.

Paths:

- `docs/product/drs/PDR<NNN>-<name>.md` (Tier 1 — own document)
- `docs/product/drs/PDR<NNN>-<name>/README.md` (Tier 2 — own directory)

Templates: [pdr-document.md](../assets/templates/pdr-document.md) (Tier 1), [pdr-directory.md](../assets/templates/pdr-directory.md) (Tier 2). See [Record elements](#record-elements) for tier selection.

### Architectural decision records

The record document names the affected element(s); each scoped element backlinks to it.

Paths:

- `docs/engineering/drs/ADR<NNN>-<name>.md` (Tier 1 — own document)
- `docs/engineering/drs/ADR<NNN>-<name>/README.md` (Tier 2 — own directory)

Templates: [adr-document.md](../assets/templates/adr-document.md) (Tier 1), [adr-directory.md](../assets/templates/adr-directory.md) (Tier 2). See [Record elements](#record-elements) for tier selection.

### Change records

CRs live in one of two central directories, chosen by which side of the tree the change touches:

- `docs/product/crs/` — changes to product, outcomes, risks, requirements, or jobs.
- `docs/engineering/crs/` — changes to architecture or components.

The record document names the affected element(s); each scoped element backlinks to it. A CR that legitimately touches both sides is recorded once on the side most central to the change, and the other side's affected element(s) still backlink to it from their `## See Also`.

Paths:

- `docs/product/crs/CR<NNN>-<name>.md` (Tier 1 — own document)
- `docs/engineering/crs/CR<NNN>-<name>.md` (Tier 1 — own document)
- `docs/product/crs/CR<NNN>-<name>/README.md` (Tier 2 — own directory)
- `docs/engineering/crs/CR<NNN>-<name>/README.md` (Tier 2 — own directory)

Templates: [cr-document.md](../assets/templates/cr-document.md) (Tier 1), [cr-directory.md](../assets/templates/cr-directory.md) (Tier 2). See [Record elements](#record-elements) for tier selection.

### Minimal directory tree

The layout below is the floor for a [minimal product](#product-scale-profiles) — two Tier 1 documents with all other elements at Tier 0 inline on them. No records are present. Capturing the first record adds the relevant central directory (`docs/product/crs/`, `docs/product/drs/`, `docs/engineering/crs/`, or `docs/engineering/drs/`) without touching any element.

```txt
docs/
  product/
    README.md
  engineering/
    README.md
```

### Maximum directory tree

The layout below is the ceiling — create each path only when its tier requires it. `.md` = Tier 1 (own document). `/<slug>/README.md` = Tier 2 (own directory).

```txt
docs/
  product/
    crs/
      CR<NNN>-<name>/
        README.md
      CR<NNN>-<name>.md
    outcomes/
      O<NNN>-<name>.md
      O<NNN>-<name>/
        requirements/
          R<NNN>-<name>/
            README.md
          R<NNN>-<name>.md
        README.md
    drs/
      PDR<NNN>-<name>/
        README.md
      PDR<NNN>-<name>.md
    README.md
  engineering/
    drs/
      ADR<NNN>-<name>/
        README.md
      ADR<NNN>-<name>.md
    components/
      C<NNN>-<name>.md
      C<NNN>-<name>/
        README.md
    crs/
      CR<NNN>-<name>/
        README.md
      CR<NNN>-<name>.md
    README.md
```

## Templates

Every element × tier has one canonical template. The parent templates ([product.md](../assets/templates/product.md), [architecture.md](../assets/templates/architecture.md), [outcome-document.md](../assets/templates/outcome-document.md), [outcome-directory.md](../assets/templates/outcome-directory.md)) are the source of truth for [child rendering](#child-rendering) — the per-element sections above link to specific rows here.

| Element | Tier | Template |
|---|---|---|
| Product | 1 | [product.md](../assets/templates/product.md) |
| Outcome | 1 | [outcome-document.md](../assets/templates/outcome-document.md) |
| Outcome | 2 | [outcome-directory.md](../assets/templates/outcome-directory.md) |
| Requirement | 1 | [requirement-document.md](../assets/templates/requirement-document.md) |
| Requirement | 2 | [requirement-directory.md](../assets/templates/requirement-directory.md) |
| Architecture | 1 | [architecture.md](../assets/templates/architecture.md) |
| Component | 1 | [component-document.md](../assets/templates/component-document.md) |
| Component | 2 | [component-directory.md](../assets/templates/component-directory.md) |
| CR | 1 | [cr-document.md](../assets/templates/cr-document.md) |
| CR | 2 | [cr-directory.md](../assets/templates/cr-directory.md) |
| PDR | 1 | [pdr-document.md](../assets/templates/pdr-document.md) |
| PDR | 2 | [pdr-directory.md](../assets/templates/pdr-directory.md) |
| ADR | 1 | [adr-document.md](../assets/templates/adr-document.md) |
| ADR | 2 | [adr-directory.md](../assets/templates/adr-directory.md) |
