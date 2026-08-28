# Security Modeling — Implementation Plan

Sequencing for [requirements.md](requirements.md). Seven phases.

**Depends on the command resurface.** Code phases need the `verb thing` grammar and the registry, so
they start after `int-xds.6` (Phase 5: remaining verbs and retirement). The help phase additionally
overlaps `int-xds.7` on two files and must follow it.

## Sequencing principle

Storage and index computations land first with no CLI change, because the derived machinery
(boundary resolution, crossings, obligations, CVSS) is where the real work and the real risk sit —
and all of it is unit-testable without a command. The CLI is a thin registry addition by comparison.
Diagrams and help come last, once the vocabulary is fixed.

```
0 ── 1 ── 2 ─┬─ 3 ── 4 ── 6
             │
             └────────── 5 ──┘        (5 also waits on int-xds.7)
```

## Phase 0 — Model, schema, index

**0.1 · Structs.** `model.Root` (`model.go:10-14`) gains `Security`. New: `Container`, `Boundary`,
`TrustLevel`, `Actor`, `Flow`, `Asset`, `Threat`, `Security`. `Engineering` gains `Containers`,
`Boundaries`, `TrustLevels`, `Actors`, `Flows`, `Assets` as `OrderedMap`s.

`Component` gains `Container`, `Boundary`, `Role`, `Stores []string`; **loses** `DataModel`
(`model.go:99`) and `Interfaces` (`model.go:100`).

`DecisionRecord` (`model.go:110-118`) gains an optional `Treatment` — meaningful only for SDRs, which
Phase 2.2 enforces. It goes on the shared struct rather than giving SDRs their own, to keep the
"location implies kind" pattern intact.

Field order is canonical output order (`model.go:4-6`) — place new fields deliberately.

**0.2 · Schema.** New `$defs` for each element; new properties on `component`; remove `data_model`
and `interfaces` from `/$defs/component`. `classification`, `role`, `category`, and `treatment` are
closed enums. `cvss` gets a pattern; `rank` gets `minimum: 0`.

**0.3 · Vocabulary.** `tree.Kind` (`tree.go:24-37`) gains 9 kinds; `Domain` (`tree.go:40-46`) gains
`DomainSecurity`; `EdgeKind` (`tree.go:48-57`) gains `EdgeStorage`. Update `validKinds` and
`validDomains` (`kinds.go:10-19`).

**0.4 · Index walk.** New `walkSecurity()`; extend `walkEngineering()` for the six new maps. Decide
`RendersAsDocument` (`tree.go:104-113`) per new kind — flows and threats want their own pages;
trust levels and assets probably render inline.

**0.5 · Mutate.** `setField` arms for every new type (`fields.go:12`); `addEdge`/`removeEdge` arms
for `EdgeStorage` (`edges.go:14,46`); new container cases in `locate.go`.

**0.6 · Dogfood `intent.yaml`.** Model this repo: one container (the single Go binary, per the
`single-go-binary` constraint), its boundary and trust level, actors (operator, coding agent, git
remote), flows between them, and the tree file itself as a `role: store` component holding the
intent-tree asset. The five existing components (`validate`, `mutate`, `gen`, `help`, `skill`) each
get a `container`.

This phase is the design's first real test. Expect it to surface at least one thing the requirements
got wrong — that is the point of doing it early rather than last.

## Phase 1 — Derived computations

All in `internal/tree`, all unit-testable with no CLI.

**1.1 · `Index.BoundaryOf(addr)`** — walks component → container → boundary, honoring a
component-level override; actors resolve directly. Returns an error for the unresolvable cases that
Phase 2 turns into findings.

**1.2 · `Index.Crossings()`** — flows whose endpoints resolve to different boundaries.

**1.3 · `Index.Obligations(addr)`** — the endpoint-kind matrix from requirements §6. `E` compares
`trust_level.rank` on both sides and is obligated only when the target's is higher.

**1.4 · `Index.Coverage()`** — obligations joined against threats, their `mitigation` sources, and
SDR `effect`s. Returns one of five derived statuses per threat — `mitigated` from the edge,
`accepted`/`transferred`/`avoided` from the affecting SDR's `treatment`, `unresolved` from neither —
and covered/uncovered per obligation.

**1.5 · CVSS.** Vector parse, version dispatch on the `CVSS:x.y/` prefix, score computation, and
`CR` derivation from the classification of assets on the threat's flows.

> ⚠️ **Biggest single implementation risk.** CVSS v3.1 scoring is a specified formula with
> roundup and impact/exploitability subscores — fiddly but tractable. **v4.0 is not a formula**: it
> scores via MacroVector lookup against a table of ~270 entries. Recommend implementing v3.1 scoring
> now, accepting and validating v4.0 vectors, and reporting v4.0 scores as unavailable until a later
> pass. Do not let v4.0 scoring block this phase.

**1.6 · Ordering.** CVSS score → `severity` band → exposure (max classification on the threat's
flows, plus a crossing flag).

## Phase 2 — Validate: severity and new rules

**2.1 · Severity.** `Finding` (`validate.go:45-49`) gains `Severity`. `String()`
(`validate.go:53-55`) includes it while keeping the trailing `→ intent help <code>` pointer.
`build` (`build.go:79-84`) refuses on **errors only**; warnings print and pass.

**2.2 · Error codes** — one per requirements §7 error: flow endpoint of the wrong kind; actor with
no boundary; component with neither container nor boundary override; container with no boundary;
boundary with no trust level; two trust levels sharing a rank; `treatment` set on a PDR or ADR;
threat with no `against` — **carrying one exemption**, for a threat an SDR with `treatment: avoid`
affects, since avoidance deleted the flow it pointed at.

**2.3 · Warning codes** — uncovered obligation; threat with neither mitigation nor SDR; actor with
no flows; boundary with no members; asset no flow carries; store with no `storage`.

> Small call to make here: warnings should get a `W0xx` prefix rather than continuing `E0xx`, so
> `intent help W001` reads correctly and the code itself carries the weight. The existing E001–E005
> are untouched.

**2.4 · `validate --strict`** — treat warnings as errors, so CI can adopt the coverage gate when the
team is ready rather than on day one.

**2.5 · Help topics** for every new code, in `internal/help/content/errors/`.

## Phase 3 — Registry and CLI

**3.1 · Registry.** Add 9 nouns to the elements table, `storage` to relationships, and a `--threat`
target on `mitigation`. Cobra groups already exist from the resurface, so the new words slot into
the Elements and Relationships groups.

**3.2 · `--domain {product,security}`** on `add requirement`, enforcing the rule that any noun with
more than one home must say which.

**3.3 · Coverage reads.** `list threat --uncovered`, `list threat --status transferred`,
`list flow --crossing`, and the coverage table as a `list` mode — all flags on existing verbs, no new
verb. `--status` takes any of the five derived statuses, which is what makes compliance export
("what have we transferred, and under which contract?") a query rather than a reading exercise.

**3.4 · Verify `impact` / `trace`** traverse the new edges. Both work off the registry and the edge
graph, so this should be assertion, not implementation.

## Phase 4 — Diagrams and generated docs

**4.1 · Architecture diagram** — Mermaid `flowchart`, subgraph per container, nodes = components and
actors, edges = collaborations labelled with their notes.

**4.2 · DFD** — Mermaid, subgraph per boundary, node shapes by `role` and kind (process rounded,
store cylinder, actor stadium), edges = flows labelled with what they carry.

**4.3 · `docs/security/`** — per-threat pages, the crossing worklist, the coverage table, a link to
the DFD.

**4.4 · Gen paths.** Slug and path rules for the new kinds (`internal/gen/path.go`), and the edge
label map (`render.go:23`) for `storage`.

Diagrams emit only when their inputs exist, following `appendText`/`appendList`
(`tree.go:282-294`). No flag.

## Phase 5 — Help content

Depends on `int-xds.7`, which rewrites the 41 existing topics for 12 nouns. This adds the security
vocabulary on top. Overlap is confined to `ref/commands.md` and the relationships topic — everything
else is new files.

**5.1 · `concept/`** — 8 new pages: container, boundary, trust-level, actor, flow, asset, threat,
sdr. Update `concept/component.md` for `role`, `storage`, and the removed fields.

**5.2 · `ref/`** — the same 8, plus updates to `ref/component.md` and `ref/requirement.md` (the
same-outcome line is now wrong: `mitigates` may target a threat).

**5.3 · `judgment/`** — this is the stated mitigation for a vocabulary going from 12 nouns to 21, so
it is not optional: `threat-vs-risk`, `boundary-vs-container`, `what-counts-as-an-asset`,
`accept-or-mitigate`, `flow-granularity`.

**5.4 · `guide/`** — one end-to-end walkthrough: model the topology, generate the worklist, discharge
or accept each obligation.

## Phase 6 — Skill, CI, dogfood

**6.1 · Regenerate the skill.** `install skill` renders from the embedded help, so agents get the
security vocabulary only once Phase 5 lands.

**6.2 · CI.** `intent validate` on errors always; `intent validate --strict` as a separate,
opt-in step.

**6.3 · `intent build`** and commit the regenerated `docs/engineering/` and new `docs/security/`.
Confirm `intent build --check` is clean.

## Risks

**CVSS v4.0 scoring** (Phase 1.5) — a lookup table, not a formula. Ship v3.1 scoring; validate v4.0
vectors without scoring them.

**Noun count** — 12 → 21 tree nouns, engineering 4 → 10. The mitigations are Phase 5.3's judgment
topics and cobra's grouped help. If Phase 0.6's dogfood makes the vocabulary feel unusable at this
repo's scale, that is the signal to cut before building further.

**The `E` obligation on lateral crossings** — `trust_level.rank` should suppress most false
positives, but if the worklist is still noisy after Phase 0.6, the escape hatch is the `trusts`
relationship named as out of scope in requirements §13.

**Two plans, one vocabulary** — the noun and verb tables appear in both requirements docs. Treat
`command-resurface/requirements.md` as canonical for the grammar and this one as canonical for the
schema.
