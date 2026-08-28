# Security Modeling — Requirements

Architecture diagrams, data flow diagrams, and a STRIDE threat model derived from the intent tree.

Status: settled. Depends on the `verb thing` grammar in
[../command-resurface/requirements.md](../command-resurface/requirements.md); sequencing lives in
[implementation-plan.md](implementation-plan.md).

## 1. Premise

`intent` **is** the threat model, not an input to one. Threats live in the tree, are discharged by
requirements or accepted by decision records, and CI checks coverage.

The reason to do this here rather than in a dedicated threat-modeling tool: a component still
fulfills requirements, so a boundary crossing traces up through requirement → risk → outcome → job.
Most threat models float free of product intent. This one doesn't.

What falls out of the specification, authored nowhere:

- an **architecture diagram** — what the system is
- a **data flow diagram** — how data moves, with trust boundaries drawn
- a **boundary-crossing worklist** — the STRIDE-per-interaction obligations
- a **coverage table** — which obligations are discharged, and how

## 2. Four domains

`model.Root` (`internal/model/model.go:10-14`) gains `Security`. `tree.Domain`
(`internal/tree/tree.go:40-46`) gains `DomainSecurity`.

| domain | owns |
|---|---|
| product | jobs, outcomes, risks, requirements, PDRs |
| engineering | architecture, containers, boundaries, trust levels, components, actors, flows, assets, principles, constraints, ADRs |
| security | threats, requirements, SDRs |
| change | CRs |

## 3. Engineering — structural facts

| noun | fields |
|---|---|
| `container` | `name`, `boundary` (address, required) |
| `boundary` | `name`, `statement`, `trust_level` (address, required once any trust level exists) |
| `trust_level` | `name`, `rank` (int ≥ 0, **unique across trust levels**), `statement` |
| `actor` | `name`, `statement`, `boundary` (address, required) — **no `fulfills`** |
| `flow` | keyed element · `name`, `from`, `to`, `data[]`, `note` |
| `asset` | `name`, `classification`: `public \| internal \| confidential \| secret` |

### Boundary resolution

```
trust_level  (name, rank)
     ▲
     │ trust_level: <addr>
boundary
     ▲                    ▲
     │ boundary:          │ boundary:  (override only)
container              component
     ▲                       ▲
     │ container:            │
component                 actor ──────┘  boundary: <addr>
```

An **actor** names its boundary directly; a **component** inherits from its container, with a
per-component override for in-process trust changes (a sandbox, a privilege drop). The asymmetry is
deliberate: a container is a deployment unit *you own*, and you don't deploy the internet, a partner
API, or the operator.

Several boundaries may share a trust level — that is how peer zones are modeled, and it is what
keeps Elevation from firing on lateral movement. Uniqueness belongs on `trust_level.rank`, where two
names for one rung would be meaningless.

### `flow`

A keyed element, not a relationship, for two reasons. Threats attach to a flow address
(STRIDE-per-interaction), and a relationship keyed by target address permits only **one** flow per
(source, target) pair — `addEdge` rejects a duplicate target (`internal/mutate/edges.go:29`) — while
`api-gateway → auth-service` carrying an auth token and the same pair carrying audit records are two
flows with different classifications.

`from` and `to` are scalar dotted paths, per the `flat-dotted-path-edges-over-nested-maps` ADR. Both
endpoints are a component or an actor. Flows are directed and one-way: a request/response pair is
two flows, which is correct, because each direction crosses the boundary separately.

### `component` changes

Gains `container`, an optional `boundary` override, `role: process | store`, and `stores[]` via the
new `storage` relationship.

**Loses `data_model`** (`model.go:99`, `ref/component.md:15`) — a data schema is the clearest case of
the implementation dump `judgment:altitude` warns against, and it is used 0 times in this repo.
**Loses `interfaces`** (`model.go:100`, `ref/component.md:16`) — subsumed by flows, which now declare
what crosses a component's edge as typed elements with classified payloads; what remained was
endpoint signatures, the same altitude as the schema.

**Keeps `behavior`** (`model.go:101`) — used, and the only component field carrying `{{ }}`
interpolation through `ProseValues` (`internal/tree/query.go:37`). **Keeps `edge_cases`** as the
`edge-case` part; nothing else in the tree records that a failure mode was considered.

## 4. Security — the annotation layer

| noun | fields |
|---|---|
| `security` | root · `name`, `summary`, `detail` |
| `threat` | `name`, `statement`, `category` (STRIDE enum), `against[]` (flow or component addresses), `severity`: `low\|medium\|high\|critical`, `cvss` (vector string, optional) |
| `requirement` | **reused, not duplicated** — `security.requirements.<r>`, sibling to threats |
| `sdr` | security decision record · same shape as PDR/ADR, `affects` ⊆ `security.*`, plus `treatment`: `accept \| transfer \| avoid` |

### Requirements live in two domains

A security requirement is a requirement whose *location* is security. This is the existing PDR/ADR
pattern: the same struct in two places, where *"location implies which one it is, so there is no
`kind` field"* (`model.go:107-109`).

`Element.Domain` alone distinguishes them, so **no new Kind is needed** — unlike PDR/ADR, which
needed two Kinds because E004 differs per kind. `list requirement --domain security` just works.

Structurally, `security{threats, requirements}` mirrors `outcome{risks, requirements}` exactly.

This is deliberately **not** a `control` noun. By the framework's own definition
(`concept/requirement.md`), a condition that must be true and is testable *is* a requirement; a
security-side duplicate would be the `success_criteria` mistake in a new place.

### Same-outcome scoping is documentation, not enforcement

`Requirement.Mitigates` is documented as pointing at same-outcome risks (`model.go:77-78`,
`ref/requirement.md:14`), but nothing enforces it: the only checks are schema-level (`mitigates`
required, non-empty), `mutate.go:59` (≥1 for a requirement), and E001/E003. So `mitigation` pointing
at a threat needs no loosening — only a docs correction.

## 5. Relationships

`mitigation` gains a `--threat` target alongside `--risk`. Both are kind-named, both differ from the
source, so R3 holds.

One new relationship:

| word | connects | flags | storage |
|---|---|---|---|
| `storage` | component[role=store] → asset | `--component` `--asset` | `Component.Stores []string` |

Target kind is fixed and differs from the source, so per R3 it is `--asset`, not `--to`. A flat
dotted-path list like `Fulfills` (`model.go:97`). No note, no payload.

`storage` earns its place on three counts. The DFD's fourth element type is incomplete without it —
a store with no declared contents is a named box. Inferring contents from inbound flows is *wrong*,
not merely inconvenient: a store holds data predating the current design, and a flow can pass
through a component without persisting, so the inbound set is neither a subset nor a superset. And
it is what makes the store's own gate computable, since CVSS `CR` derives from
`asset.classification`.

## 6. Derived — nothing authored

Per `derive-only-mechanical-duals`: *"Derive only the mechanical duals of what's declared."*

- a component's boundary, from its container, unless overridden
- **crossing** = source boundary ≠ target boundary
- a threat's **resolution status**, from the graph plus the SDR's `treatment` (§8)
- the obligated STRIDE set per interaction (below)
- CVSS base / temporal / environmental / overall scores, from the vector
- CVSS `CR`, from the classification of assets on the threat's flows
- ordering: CVSS score → `severity` band → exposure (max classification + crossing flag)
- the architecture diagram, the DFD, the crossing worklist, the coverage table

### Obligated categories, from endpoint kinds

STRIDE-per-interaction assigns categories by the element types involved, not a fixed row. A fixed
`T/I/D` is the per-*element* data-flow row and would leave Spoofing and Elevation unchecked
everywhere.

| interaction (crossing a boundary) | obligated |
|---|---|
| actor → component | **S** · T · I · D |
| component → actor | T · I · D · **R** |
| component → component | T · I · D · **E** only when target rank > source rank |
| component ↔ store | T · I · D |
| store (element-level) | T · I · D · **R** |

Actors and processes need no separate gate: the interactions that involve them carry S and R. E
fires only on genuine elevation, which is what `trust_level.rank` is for.

Obligations are a **floor, not a ceiling** — nothing prevents authoring a threat on a same-boundary
flow; the gate simply won't demand one.

### CVSS

Author the **vector**; derive every score. A CVSS score is a deterministic arithmetic function of
its vector — a mechanical dual — and storing the numbers would churn four values per metric change,
against `clean-diffs-are-a-design-goal`.

```yaml
security.threats.token-replay:
  category: information-disclosure
  severity: high
  cvss: "CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:N/A:N/E:P/RL:O"
  # base 7.5 · temporal 6.7 · environmental 8.1 · overall 8.1 — all derived
  # CR:H derived from asset "auth-token" (secret)
```

Version needs no decision: a vector self-declares (`CVSS:3.1/…` vs `CVSS:4.0/…`), so accept both and
dispatch on the prefix. The `severity` band uses CVSS's own vocabulary (`low|medium|high|critical`,
matching its 0.1–3.9 / 4.0–6.9 / 7.0–8.9 / 9.0–10.0 ratings) so a tree with some threats scored and
some banded still sorts coherently.

DREAD is deliberately excluded: Microsoft dropped it from the SDL for poor inter-assessor
reproducibility, averaging five subjective 1–10 ratings manufactures precision, and Discoverability
rewards security through obscurity.

## 7. Validate gains severity

`validate.Finding` (`internal/validate/validate.go:45-49`) gains `Severity: error | warning`.
`build` refuses on **errors only** — today it refuses on any finding at all
(`internal/cli/build.go:79-84`), which would mean an incomplete threat model blocks doc generation
and incremental adoption is impossible.

**Errors** — the tree cannot be evaluated:

- a flow endpoint that is not a component or actor
- an `actor` with no `boundary`
- a `component` with neither a `container` nor a `boundary` override
- a `container` with no `boundary`
- a `boundary` with no `trust_level`, once any trust level exists
- two `trust_level`s sharing a `rank`
- a `threat` with no `against` — **unless** an SDR with `treatment: avoid` affects it (§8)
- `treatment` set on a PDR or an ADR, where it is meaningless

**Warnings** — work is outstanding:

- an uncovered obligation on a crossing or store
- a `threat` with neither a mitigation nor an SDR accepting it
- an `actor` with no flows
- a `boundary` with no members
- an `asset` no flow carries
- a `component` with `role: store` and no `storage`

Every new noun earns its place by participating — the same discipline as *"every component should
earn its place by fulfilling a requirement"* (`concept/component.md:22`).

## 8. Discharge, and the five resolution states

A threat's resolution status is **derived**, never authored. Three states come from the graph alone;
the other two come from the SDR's `treatment`, which is the only new authored field:

| status | derivation |
|---|---|
| `mitigated` | ≥1 incoming `mitigation` from a security requirement |
| `accepted` | incoming `effect` from an SDR with `treatment: accept` |
| `transferred` | incoming `effect` from an SDR with `treatment: transfer` |
| `avoided` | incoming `effect` from an SDR with `treatment: avoid` |
| `unresolved` | neither a mitigation nor an SDR effect |

This is the standard four-Ts treatment taxonomy — treat, tolerate, transfer, terminate — with
"treat" carried by the mitigation edge rather than a field.

**Treatment belongs on the SDR, not the threat.** An authored `status` on the threat would duplicate
state the graph already knows, which is the objection that retired `success_criteria` and stored CVSS
scores; and it would create a drift surface (`status: mitigated` with no mitigation edge) requiring a
validate rule to police a contradiction the duplication itself introduced. Putting `treatment` on the
decision keeps the threat's status fully derived, and records *why* alongside *what* — a transfer
without a recorded contract is not auditable.

`treatment` values are exactly the non-mitigation postures. An SDR declaring "we decided to mitigate"
would be redundant with the edge. One SDR carries one treatment; mixed postures need two SDRs, which
is better hygiene anyway.

**An avoided threat is the only threat allowed to have no `against`.** If you avoided it by deleting
the flow, the target is gone — so that error carries one exemption, and avoided threats become a
permanent record of paths not taken. That is the same instinct the framework already applies to
decision records: keep the reasoning, not just the current state.

### Discharge

An obligation is met by **≥1 security requirement** pointing at the threat via `mitigation`, or
**accepted by an SDR** pointing at the threat via `effect` — with any of the three treatments.

The SDR path needs zero new schema: `effect` already targets any element, and threats live under
`security.*` so `checkDomainScope` (`validate.go:137-160`) holds — it hard-prefixes each record
type's `affects` to its own domain (`validate.go:152`), which is precisely why an **ADR cannot** waive
a security threat and a third record type is required.

One SDR can accept many threats, so a blanket decision ("Elevation is not applicable between peer
zones") collapses the burden and records context, options, decision, and consequences — what an
auditor wants, and what a one-line rationale field would not have given.

## 9. Generated docs

Diagrams are Mermaid: text-based, diffable, rendered by GitHub, zero dependencies.

- `docs/engineering/` — the architecture diagram (subgraph per container) and the DFD (subgraph per
  boundary, node shapes by `role` and kind).
- `docs/security/` — the crossing worklist, the coverage table, per-threat pages, and a link to the
  DFD.

Engineering owns topology; security owns analysis. Diagrams emit **when their inputs exist** — no
flag — following `appendText`/`appendList`, which already omit empty sections
(`internal/tree/tree.go:282-294`).

## 10. Grammar impact

**21 tree nouns** (12 → +6 engineering: `container`, `boundary`, `trust_level`, `actor`, `flow`,
`asset`; +3 security: `security`, `threat`, `sdr`). **6 relationships** (+`storage`). **4 parts**
(unchanged). Three component fields deleted.

`--domain {product,security}` is required on `add requirement`, per the rule: **`--domain` is
required whenever a noun has more than one possible home.** Today that is only `requirement`.

```
add requirement token-ttl  --domain product  --outcome fast-checkout --mitigation cart-abandon
add requirement bind-token --domain security                          --mitigation token-replay
```

Cobra command groups (`AddGroup`/`GroupID`, cobra v1.10.2) keep `intent add --help` legible as the
namespace grows.

The noun count is the single largest risk to this design. Engineering goes from 4 nouns to 10. The
mitigation is judgment topics and grouped help, not schema.

## 11. Layering

**Product is a pure sink** — product intent reviews cleanly with no security section present.

`engineering ↔ security` is mutually referencing: threats point at flows, and components fulfill
security requirements. That cycle is contained to two domains and is the price of tracing a control
to the thing that implements it. Full acyclicity was never a stated requirement.

## 12. Decisions, and where the reasoning changed

| # | decision |
|---|---|
| 1 | intent **is** the threat model, not an input to one |
| 2 | threats attach to flows, as keyed elements — STRIDE-per-interaction |
| 3 | `boundary` (trust) and `container` (deployment) are separate nouns |
| 4 | boundary declared on the container; component may override |
| 5 | a **requirement** discharges a threat — no `control` noun |
| 6 | `threat.category` is a STRIDE enum; CI checks coverage, not count |
| 7 | waivers are SDRs pointing at the threat via `effect` |
| 8 | `asset` survives with a closed `classification` enum |
| 9 | a third **security** domain, thin: `threat`, `requirement`, `sdr` |
| 10 | CVSS vector authored, scores derived; DREAD excluded; `severity` band CVSS-aligned |
| 11 | `trust_level` is an element with a unique `rank`; ties across boundaries allowed |
| 12 | `data_model` and `interfaces` dropped; `behavior` and `edge-case` kept |
| 13 | obligated categories derived from endpoint kinds |
| 14 | `Finding` gains severity; build refuses on errors only |
| 15 | security gets no `principles`/`constraints` — those stay engineering |
| 16 | threat status is derived; `treatment` (`accept\|transfer\|avoid`) lives on the SDR |

Reversals worth recording, so the reasoning isn't re-litigated:

- **`component.role` is not load-bearing for the flow gate.** It is needed for DFD node shape and
  for the store gate. An earlier draft claimed the flow gate required it.
- **A fixed `T/I/D` row was wrong.** That is the per-element data-flow row; the per-interaction
  method makes the category set a function of endpoint kinds.
- **Nesting threats under flows was over-argued.** `ix.Incoming` (`tree.go:314`) already makes the
  reverse join cheap, so pointers cost far less than claimed — and they let one threat span several
  flows, which is the common case.
- **Full acyclicity was a tiebreaker, not a requirement.** Containing the cycle to
  engineering ↔ security while keeping product a sink is the property that matters.
- **Uniqueness on trust belongs on the rank, not on boundary assignment.** Forcing distinct levels
  per boundary reinstates a total order and mis-gates peer zones in both directions.

## 13. Out of scope

Container nesting (containers are flat; C4's system > container > component is not modeled).
Boundary nesting. A `trusts` relationship for asymmetric peer trust — the rank handles the common
case, and this is the escape hatch if it doesn't. Any `--stereotype` field. LINDDUN, PASTA, or
attack trees. Migration tooling: v0, one consumer, hand-edit.
