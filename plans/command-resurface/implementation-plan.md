# Command Resurface — Implementation Plan

Sequencing for [requirements.md](requirements.md). Nine phases. Every file:line below is verified
against the tree at time of writing.

## Sequencing principle

The CLI surface is the *last* thing that should churn and the *first* thing everything else depends
on. So: vocabulary changes land first behind the old CLI (Phase 0–1), the new surface is generated
from one data table rather than hand-written (Phase 2), verbs are built over that table (Phase 3–5),
and the help content — the actual bulk of the work — comes last, once the examples it quotes are
real (Phase 6).

```
0 ─┬─ 2 ─┬─ 3 ─┬─ 5 ── 6 ── 7
   │     │     │
1 ─┘     └─ 4 ─┘        └──── 8
```

## Phase 0 — Vocabulary and storage

Land the renames while the old CLI still works, so this diff is reviewable on its own.

**0.1 · Drop `success_criteria`**
- `internal/model/model.go:103` — remove `Component.SuccessCriteria`.
- `internal/tree/tree.go:242` — remove its `appendList`.
- `internal/schema/intent.schema.json` — remove from `/$defs/component` (which is
  `additionalProperties: false`, so this is a hard reject afterwards).
- `internal/help/content/ref/component.md:19`, `concept/component.md:17` — remove the rows.
- `intent.yaml` — 0 occurrences. Nothing to migrate.

**0.2 · `relationships` → `collaborations`**
- `internal/model/model.go:98` — field and `yaml:"collaborations,omitempty"` tag.
- `internal/tree/tree.go:56` — `EdgeRelationships` → `EdgeCollaboration`, value `"collaboration"`.
- `internal/tree/tree.go:236` — the walk.
- `internal/gen/render.go:14,23` — the edge-order slice and **the explicit label map**
  (`"Relationships"` → `"Collaborations"`). Note this is *not* covered by `fieldLabel`
  (`render.go:299`), which only humanizes snake_case *field* keys; edge labels are a hand-written
  map.
- `internal/validate/validate.go:115` — the `checkDuplicates` special-case.
- `internal/mutate/edges.go:27,61` — the `addEdge`/`removeEdge` arms.
- `internal/cli/link.go:19` — the `edgeKinds` map key (dies in 5.2, but must compile until then).
- `internal/schema/intent.schema.json` — `/$defs/component`.
- `internal/help/content/ref/edges.md:15`.
- `intent.yaml` — 4 occurrences: lines 214, 225, 236, 250.

**0.3 · `KindEngineering` → `KindArchitecture`**

Only six non-test sites. The subtlety worth stating in the commit message: **the Kind *value*
becomes `"architecture"`; the element's *address* stays `engineering`.**

- `internal/tree/tree.go:26` — `KindArchitecture Kind = "architecture"`.
- `internal/tree/tree.go:106` — `RendersAsDocument`.
- `internal/tree/tree.go:205` — the walk. `Addr: "engineering"` is unchanged; only `Kind:` moves.
- `internal/tree/kinds.go:11` — `validKinds`.
- `internal/validate/validate.go:167` — the E005 root exemption.
- `internal/mutate/move.go:39` — the "roots have a fixed home" guard.
- `DomainEngineering` and every `"engineering"` address literal are untouched.

**Exits when:** `go test ./...` green, `intent validate` clean, `intent build` produces a diff only
where the "Relationships" heading became "Collaborations".

## Phase 1 — Index capabilities

Additive. No CLI change; each lands with its own unit tests.

**1.1 · `ResolveIn(query, ResolveOpts{Kind, Scope})`** in `internal/tree/address.go`
- Exact full address still wins first (`address.go:39`), then `Kind` filter, then `Scope` prefix
  filter, then `segmentSuffix` (`address.go:69`) over the narrowed set.
- `Resolve(q)` becomes `ResolveIn(q, ResolveOpts{})` — byte-for-byte today's behavior, which is the
  regression test.
- `NotFoundError`/`AmbiguousError` gain the scope in their messages (`no requirement matches "auth"`).

**1.2 · `Index.Impact(addr)`** in `internal/tree/query.go`
- Transitive incoming closure, BFS, visit-once, address-sorted.
- A dedicated method rather than a sentinel on `Trace`: `Trace` clamps negative `maxDepth` to 0
  (`query.go:74-76`) and that clamp is load-bearing.

**1.3 · `Index.DerivedCriteria(componentAddr)`** in `internal/tree/query.go`
- Union of `acceptance_criteria` across the component's `fulfills` edges. Dedup, stable order,
  attribute each line to its source requirement so `list` can show provenance.

## Phase 2 — The registry (keystone)

Roughly 100 subcommands are implied by the matrix in requirements §3. They must be **generated from
data**, not hand-written. New file `internal/cli/registry.go` with three tables and no cobra:

```go
var elements = []Element{...}          // word, kind, parent kind, parent flag, scalars, parts, verbs
var parts = []Part{...}                // word, owner kind, list field, derived bool
var relationships = []Relationship{...} // word, edge kind, source/target kind, source/target flag, note
```

- Scalar specs carry only **flag name → model field name**, because `mutate.Set` is already generic
  by field-name string (`internal/mutate/fields.go:12`) and already rejects list fields and edges
  with good errors (`fields.go:171,179`). Heavy reuse; no new setter layer.
- Cobra group IDs live here too, feeding `AddGroup`/`GroupID` (cobra v1.10.2).

**Exits when:** a consistency test proves every entry is internally coherent — every element's parent
flag matches its parent kind, every part's owner is a real kind, every relationship's edge kind
exists in `tree`, and every word is unique across all three tables.

## Phase 3 — Read verbs

Depends on 1, 2.

- **3.1** `show <element> <loc>`, `tree [<element> <loc>]` — locator positional, kind-scoped
  resolution.
- **3.2** `list [<thing>] [query] --domain` — one command serving all three classes. Elements filter
  by kind (retiring `find --type`); parts read the owner's list, with `criterion --component`
  routing to `DerivedCriteria` and showing provenance; relationships list from either end.
- **3.3** `trace <element> <loc> --via --up --down --depth` (`--edges` → `--via`), and
  `impact <element> <loc>` over `Index.Impact`.

## Phase 4 — Write verbs

Depends on 1, 2. Parallel with 3.

- **4.1** `add <element> <key>` — parent flag, scalars, part seeds, relationship seeds. Keeps the
  required-at-creation rule (`internal/cli/add.go:38-42`).
- **4.2** `set <element> <loc>` — scalars (now multiple per call), plus `--key`, `--<parent>`,
  `--before`/`--after`, absorbing `mv` (`internal/cli/mv.go:78-80`). Loses `type` and `--cascade`.
- **4.3** `rm <element> <loc>`.
- **4.4** Parts: `add`/`rm`/`list`. `--index N` (1-based, matching what `list` prints) or `--match`,
  exactly one required. `add criterion --component X` refused with a message naming the derived view.
- **4.5** Relationships: `add`/`rm`. Three things to carry over:
  - Root-literal targets — `--to architecture` must write `engineering`. The translation belongs at
    `internal/mutate/locate.go:39`, which already has the `case "engineering":` arm.
  - The **dangling-target fallback** from `resolveEdgeTarget` (`internal/cli/link.go:73-85`): match
    the query against addresses the source actually declares, so a hand-edited break stays fixable.
  - Endpoint kind-checking, which `addEdge` does not do today (`internal/mutate/edges.go:14`).
- **4.6** `promote` / `demote --cascade`, taking over the cascade logic from `mutate.SetType`.

**4.7 · Footers and mutate-layer error text.** Easy to miss, and user-facing:
- `internal/cli/footer.go:49,55,61` build `intent show <addr>` — the new form needs the element's
  **kind**, so the footer helpers must take the element, not just the address.
- `internal/cli/footer.go:90` emits `intent link <kind> <a> <b>` → must become
  `intent add <relationship-word> --<source> <a> --to <b>`.
- `internal/mutate/fields.go:171,179` reference `intent link` / `intent set … type` in errors.
- `internal/cli/promote.go:16-17`, `add.go:39,47`, `set.go:17`, `rm.go:51` all quote old commands.

## Phase 5 — Remaining verbs and retirement

Depends on 3, 4.

- **5.1** `build --check` (folding `check`, same `--out` resolution at `internal/cli/build.go:65`);
  `install skill` (from `install-skill`); `help` unchanged; cobra groups wired on `add`/`rm`/`list`.
- **5.2** Delete `find`, `affects`, `check`, `link`, `unlink`, `mv`, `record`, `install-skill` and the
  old `add`/`set`/`rm` forms, plus their tests. `newRootCmd` (`internal/cli/root.go:37`) shrinks to
  the 14 verbs.

## Phase 6 — Help content

Depends on 5, so every quoted command is real. This is the bulk of the work: 41 files in
`internal/help/content/`. Six independent chunks — parallelizable.

- **6.1** `ref/commands.md`, `ref/addressing.md` (the kind-scoped resolution rules), `ref/syntax.md`.
- **6.2** `ref/edges.md` → rewrite as relationships; consider renaming the topic slug.
- **6.3** `ref/*` per-element pages (job, outcome, risk, requirement, component, pdr, adr, cr), and
  **add the two missing ones**: `ref/principle.md`, `ref/constraint.md`.
- **6.4** `concept/*` pages, including the CRC mapping in `concept/component.md`; **add
  `concept/principle.md`, `concept/constraint.md`** — also currently missing.
- **6.5** `guide/*` — five files: `add-requirement`, `draft-judge-review`, `edit`, `record-decision`,
  `trace`.
- **6.6** `errors/*` — five files; E001–E005 messages name commands.

## Phase 7 — Skill and dogfood

Depends on 6. The help *is* the skill.

- **7.1** `internal/skill/render.go:54` and `agentsmd.go:15` hardcode `intent install-skill`.
  - ⚠️ `agentsBeginMarker` (`agentsmd.go:15`) is matched by exact `strings.Index`
    (`agentsmd.go:53`). Changing its text orphans previously-installed blocks and produces a
    duplicate section on the next install. **Freeze the marker string** and comment why — it is an
    identifier, not documentation. (Acceptable alternative at v0: change it and hand-clean the one
    affected repo.)
- **7.2** `intent.yaml` prose naming old commands — e.g. the acceptance criterion at line 105
  (`add, set, link, mv, promote, and record …`) and component `behavior` fields.
- **7.3** `intent build`, commit regenerated `docs/product/` and `docs/engineering/`, confirm
  `intent build --check` is clean.

## Phase 8 — Test sweep

Depends on 5. Runs alongside 6–7.

- `internal/cli/*_test.go` is written against the flat surface throughout — the largest test rewrite.
- `structural_test.go`, `help_test.go` (topic index counts change in Phase 6), `footer_test.go`
  (Phase 4.7 changes footer output).
- Add coverage for the genuinely new behavior: kind-scoped resolution and its ambiguity messages,
  transitive `impact`, derived component criteria, endpoint kind-checking, root-literal translation.

## Deliberately out of scope

Version bump and release; any migration tooling (v0, single consumer, hand-edit); `--stereotype` or
any other new component field.
