# Structure

How elements and records map to files under `docs/`. Create a path only when the complexity tier for that element calls for it. Tier 0 elements have no dedicated path — they are recorded on a parent document listed below.

## Complexity Tiers

Not every element needs its own file, folder, or decision history. Use the **simplest tier that still captures what readers and implementers need**. The hard default: **every element starts at Tier 0** — a statement recorded inline on its parent document — unless a concrete trigger below forces promotion. (Fixed exceptions: the product and architecture are always Tier 1; jobs and risks are always Tier 0.)

**Elements** use three tiers:

| Tier | Name | On disk |
|---|---|---|
| 0 | **Statement** | No dedicated file. The minimal pattern lives on the parent element. |
| 1 | **Simple** | Single `<slug>.md` file. |
| 2 | **Nested** | Folder `<slug>/README.md` with optional child folders for requirements, records, or materials. |

**Records** use the same three tiers, plus **none** when no record is warranted:

| Tier | Name | On disk |
|---|---|---|
| — | **None** | No record. |
| 0 | **Inline** | No dedicated file. Recorded in the appendix of the element's document, or of the parent document when the element is Tier 0. |
| 1 | **Simple** | Single record file. |
| 2 | **Nested** | Record folder with supplementary materials. |

**Logical IDs** are assigned at every tier so tracing works even when there is no file.

The intent framework scales from a single product document to a fully nested tree with records and supplementary materials.

**Promote an element from Tier 0 → Tier 1** when it needs ANY of: its own acceptance criteria; enumerated edge cases; dependencies or cross-references that don't fit inline; its own change or decision record file; or a child file or folder.

**Promote an element from Tier 1 → Tier 2** when it needs child folders — requirement files under an outcome, record files (`crs/`, `pdrs/`, `adrs/`), or supplementary materials that need a stable home.

**Promote a record** by the same kind of concrete trigger: none → inline when an entry is worth capturing; inline → simple when the entry needs its own file; simple → nested when the record needs supplementary files.

Demote only when you are deliberately simplifying scope — never to hide unresolved ambiguity.

#### By element

**Product and jobs** — The product is always Tier 1 (`docs/product/README.md`). Jobs are always Tier 0 statements on that document, keyed `J<NNN>`.

**Outcomes and risks** — Outcomes start at Tier 0 on the product document. Promote to Tier 1 when risks and requirements need their own document but still fit inline. Promote to Tier 2 when you need requirement files, outcome-scoped record files, or supplementary materials. Risks are always Tier 0 on the outcome.

**Requirements** — Start at Tier 0 on the owning outcome. Promote to Tier 1 for detail, edge cases, acceptance criteria, or dependencies (requires a Tier 2 nested outcome folder). Promote to Tier 2 when the requirement accumulates record files or supplementary materials.

**Architecture and components** — Architecture is always Tier 1 (`docs/engineering/README.md`). Components start at Tier 0 on the architecture document; promote to Tier 1 for interfaces, behavior, or success criteria; to Tier 2 for record files or large supplementary specs.

**Records** — Change records use all four levels (none, inline, simple, nested). When the scoped parent has no `crs/` subfolder, use inline. A Tier-0 element's inline CR lives in the appendix of the parent document that hosts the element, keyed by the element-scoped logical ID. Product- and engineering-scoped simple and nested CRs do not require promoting the parent. PDRs and ADRs use none, simple, or nested only.

#### Product scale profiles

**Minimal** — Two documents (product and architecture) with jobs, outcomes, risks, requirements, and components at Tier 0. Inline CRs in appendices when worth capturing. See the [minimal directory tree](#minimal-directory-tree).

**Standard** — Outcomes at Tier 1 or 2. Requirements and components at Tier 0 or Tier 1 where detail is needed. Records when material decisions or changes occur.

**Complex** — Full nesting where warranted; use the [maximum directory tree](#maximum-directory-tree) as the ceiling.

#### Choosing a tier

**Elements:** (1) minimal pattern enough on parent → Tier 0; (2) own document, no child folders → Tier 1; (3) child folders needed → Tier 2.

**Change records:** (1) trivial edit → none; (2) brief appendix entry → inline; (3) own file, no supplementary folder → simple; (4) supplementary files → nested.

**Decision records:** (1) no decision worth capturing → none; (2) own file → simple; (3) supplementary files → nested.

**Promote to Tier 1 when** the element needs any of: its own acceptance criteria; enumerated edge cases; dependencies or cross-references that don't fit inline; its own record file; or a child file. **Promote to Tier 2 when** it needs child folders (requirement files, record files, or supplementary materials). **Promote a record when** the entry needs its own file (inline → simple) or needs supplementary files (simple → nested).

**Stay lower when** the product is small and the parent stays readable; detail is still provisional.

## Naming

File and folder names (**slugs**) use only the local ID segment for that element type. Parent context is encoded by the directory path — do not repeat ancestor IDs in slugs.

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
| Risk | (recorded on outcome doc) | `O<NNN>-RSK<NNN>` |
| Requirement | `R<NNN>-<name>` | `O<NNN>-R<NNN>` |
| Component | `C<NNN>-<name>` | `C<NNN>` |
| Product PDR | `PDR<NNN>-<name>` | `PDR<NNN>` |
| Outcome PDR | `PDR<NNN>-<name>` | `O<NNN>-PDR<NNN>` |
| Requirement PDR | `PDR<NNN>-<name>` | `O<NNN>-R<NNN>-PDR<NNN>` |
| Inline CR (Tier-0 element) | (recorded in parent doc appendix) | `<element-id>-CR<NNN>` |
| Product CR | `CR<NNN>-<name>` | `PROD-CR<NNN>` |
| Outcome CR | `CR<NNN>-<name>` | `O<NNN>-CR<NNN>` |
| Requirement CR | `CR<NNN>-<name>` | `O<NNN>-R<NNN>-CR<NNN>` |
| Architecture ADR | `ADR<NNN>-<name>` | `ADR<NNN>` |
| Component ADR | `ADR<NNN>-<name>` | `C<NNN>-ADR<NNN>` |
| Engineering CR | `CR<NNN>-<name>` | `ENG-CR<NNN>` |
| Component CR | `CR<NNN>-<name>` | `C<NNN>-CR<NNN>` |

Product- and engineering-level CRs share the same slug form; the directory path (`docs/product/crs/` vs `docs/engineering/crs/`) distinguishes the files, and the `PROD-`/`ENG-` prefix keeps their logical IDs unambiguous.

## Paths

All paths are under `docs/`. Create a path only when the element's tier calls for it.

### Product

Tier 1 — simple

Path: `docs/product/README.md`

Template: [product.md](../assets/templates/product.md)

### Outcomes

Tier 0 — statement

Path: none — recorded on `docs/product/README.md`

Tier 1 — simple

Path: `docs/product/outcomes/O<NNN>-<name>.md`

Template: [outcome-simple.md](../assets/templates/outcome-simple.md)

Tier 2 — nested

Path: `docs/product/outcomes/O<NNN>-<name>/README.md`

Template: [outcome-nested.md](../assets/templates/outcome-nested.md)

Child folders: `requirements/`, `pdrs/`, `crs/`

### Requirements

Tier 0 — statement

Path: none — recorded on the owning outcome document

Tier 1 — simple

Path: `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>.md`

Template: [requirement-simple.md](../assets/templates/requirement-simple.md)

Requires a Tier 2 nested outcome folder.

Tier 2 — nested

Path: `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>/README.md`

Template: [requirement-nested.md](../assets/templates/requirement-nested.md)

Child folders: `pdrs/`, `crs/`

### Architecture

Tier 1 — simple

Path: `docs/engineering/README.md`

Template: [architecture.md](../assets/templates/architecture.md)

### Components

Tier 0 — statement

Path: none — recorded on `docs/engineering/README.md`

Tier 1 — simple

Path: `docs/engineering/components/C<NNN>-<name>.md`

Template: [component-simple.md](../assets/templates/component-simple.md)

Tier 2 — nested

Path: `docs/engineering/components/C<NNN>-<name>/README.md`

Template: [component-nested.md](../assets/templates/component-nested.md)

Child folders: `adrs/`, `crs/`

### Product decision records

Tier 1 — simple

Paths:

- `docs/product/pdrs/PDR<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/pdrs/PDR<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>/pdrs/PDR<NNN>-<name>.md`

Template: [pdr-simple.md](../assets/templates/pdr-simple.md)

Tier 2 — nested

Paths:

- `docs/product/pdrs/PDR<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/pdrs/PDR<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>/pdrs/PDR<NNN>-<name>/README.md`

Template: [pdr-nested.md](../assets/templates/pdr-nested.md)

### Architectural decision records

Tier 1 — simple

Paths:

- `docs/engineering/adrs/ADR<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/ADR<NNN>-<name>.md`

Template: [adr-simple.md](../assets/templates/adr-simple.md)

Tier 2 — nested

Paths:

- `docs/engineering/adrs/ADR<NNN>-<name>/README.md`
- `docs/engineering/components/C<NNN>-<name>/adrs/ADR<NNN>-<name>/README.md`

Template: [adr-nested.md](../assets/templates/adr-nested.md)

### Change records

Tier 0 — inline

Path: none — recorded in the appendix of the parent document that hosts the Tier-0 element (the element has no document of its own). Key the entry by the element-scoped logical ID, e.g. `O<NNN>-CR<NNN>` for a Tier-0 outcome on `docs/product/README.md`, or `C<NNN>-CR<NNN>` for a Tier-0 component on `docs/engineering/README.md`.

Template: [cr-inline.md](../assets/templates/cr-inline.md)

Tier 1 — simple

Paths:

- `docs/product/crs/CR<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/crs/CR<NNN>-<name>.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>/crs/CR<NNN>-<name>.md`
- `docs/engineering/crs/CR<NNN>-<name>.md`
- `docs/engineering/components/C<NNN>-<name>/crs/CR<NNN>-<name>.md`

Template: [cr-simple.md](../assets/templates/cr-simple.md)

Outcome-, requirement-, and component-scoped Tier 1 paths require a Tier 2 nested parent. Product- and engineering-scoped Tier 1 paths do not.

Tier 2 — nested

Paths:

- `docs/product/crs/CR<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/crs/CR<NNN>-<name>/README.md`
- `docs/product/outcomes/O<NNN>-<name>/requirements/R<NNN>-<name>/crs/CR<NNN>-<name>/README.md`
- `docs/engineering/crs/CR<NNN>-<name>/README.md`
- `docs/engineering/components/C<NNN>-<name>/crs/CR<NNN>-<name>/README.md`

Template: [cr-nested.md](../assets/templates/cr-nested.md)

Same scoping rules as Tier 1.

### Minimal directory tree

The layout below is the floor for a [minimal product](#product-scale-profiles) — two Tier 1 documents with all other elements at Tier 0 inline on them. Material changes may appear as inline CRs in appendices; no other paths are required.

```txt
docs/
  product/
    README.md
  engineering/
    README.md
```

### Maximum directory tree

The layout below is the ceiling — create each path only when its tier requires it. `.md` = Tier 1 simple file. `/<slug>/README.md` = Tier 2 nested folder.

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
        crs/
          CR<NNN>-<name>/
            README.md
          CR<NNN>-<name>.md
        pdrs/
          PDR<NNN>-<name>/
            README.md
          PDR<NNN>-<name>.md
        requirements/
          R<NNN>-<name>/
            crs/
              CR<NNN>-<name>/
                README.md
              CR<NNN>-<name>.md
            pdrs/
              PDR<NNN>-<name>/
                README.md
              PDR<NNN>-<name>.md
            README.md
          R<NNN>-<name>.md
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
          ADR<NNN>-<name>/
            README.md
          ADR<NNN>-<name>.md
        crs/
          CR<NNN>-<name>/
            README.md
          CR<NNN>-<name>.md
        README.md
      C<NNN>-<name>.md
    crs/
      CR<NNN>-<name>/
        README.md
      CR<NNN>-<name>.md
    README.md
```
