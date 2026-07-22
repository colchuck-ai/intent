# Intent — Structured-Source Design

**Status:** working draft (in progress), revised after a design-review pass. This is the shared
design doc for pivoting Intent from a hand-maintained multi-file markdown tree to a **single
schema'd `intent.yaml` operated by a read/write `intent` CLI that generates the docs on demand**.
This is a full reset of the skill.

## Guiding principle

> **Declare intent explicitly; derive only the mechanical duals of what's declared.**

Nothing infers *intent*. Not tier (was: presence of a `body` → now: explicit `type`). Not
link text (was: auto from name → now: written or an explicit accessor). Not dependencies
(never from prose → only declared edges). Not image-vs-link (was: extension sniffing → now:
explicit markdown at the point of use).

What we *do* still derive are the mechanical duals of declared state: backlinks, See-Also
blocks, and maps from declared edges; directory structure from declared `type`s; link hrefs
and file paths from the tree layout.

## Operating model (who does what)

- **Writes are agent-driven** through the `intent` CLI — the happy path. Hand-editing
  `intent.yaml` stays *possible* but is **not a design driver**: we do not optimize the raw file
  for human authoring.
- **Humans review the MR diff.** The reviewed artifacts are the **`intent.yaml` diff** *and* the
  **generated-docs diff** (docs are committed — §8). **Clean, low-noise diffs are therefore an
  explicit design goal**, not a nice-to-have — everything downstream (canonical re-serialize,
  deterministic generation, the derived-content layout) exists to serve reviewability.
- **The CLI guarantees *validity*; judgment guarantees *quality*.** The CLI is deterministic —
  it never refuses a valid-but-mediocre element. Judgment (is this really a requirement? the
  right altitude?) lives in editorial content that *ships with the binary* and is surfaced to the
  operator (§12).

---

## 1. File & root structure

**One `intent.yaml`**, validated by one schema, with three top-level roots. (Earlier drafts
sharded into three files for selective loading; the read/write CLI — §11 — made that moot: the
CLI is the read interface, so nothing loads the raw file, and it mediates writes surgically, so
single-file merges stay small.)

```
intent.yaml           roots: product, engineering, change_records
intent.schema.json    the contract
intent.config.yaml    optional generator/CLI config (see §7)
```

- `product` → `name`, `summary`, `jobs{}`, `decision_records{}` (the PDRs)
- `engineering` → `name`, `summary`, `principles[]`, `constraints[]`, `components{}`, `decision_records{}` (the ADRs)
- `change_records{}` → top-level, cross-cutting

**Decision records nest by domain, and are domain-pure.** A PDR is always product, an ADR always
engineering, so they live with the thing they decide about — and their `affects` targets are
constrained to their own domain (`product.*` for PDRs, `engineering.*` for ADRs), linter-enforced
(§4, §9). **Change records stay top-level** because a CR is the one record type that legitimately
crosses domains — one CR can `affect` a requirement *and* a component at once.

> **Why the asymmetry (decisions domain-pure, changes cross-cutting):** *reasoning* decomposes by
> domain, *change-events* don't. What looks like a cross-domain decision is almost always a product
> decision (the *what*) plus an engineering decision (the *how*), joined by a `fulfills` edge between
> elements — the trace graph carries the cross-domain story, not the record. A change, by contrast,
> is a single atomic modification event that can't be split without lying about it. See §4.

Nesting:

```
product.jobs.<job>.outcomes.<outcome>.risks.<risk>
product.jobs.<job>.outcomes.<outcome>.requirements.<requirement>
engineering.components.<component>
```

---

## 2. Element model

Every element is a dict entry. Its **snake_case key is its identity** — also its path segment
and slug basis. There is no `id` field (a key can't disagree with itself).

Common fields:

- `name` — **required.** Display / heading text. Free to diverge from the key; nothing
  resolves on it. **Never required to be unique** — links resolve by computed path, so names
  may collide freely.
- `statement` (outcomes/risks/requirements) or `story` (jobs) — the prose. **The prose-field
  name is intentional semantic typing** (a job's `story` is a JTBD narrative, an outcome/risk/
  requirement `statement` is a declarative sentence, a component's `responsibility` is a one-line
  charter, a root/DR `summary` is an elevator pitch). The name prompts the author toward the right
  shape; it is fixed per element type, never chosen. Shapes are pinned in `ref:*` +
  `judgment:writing-statements` (§12). We keep the distinct names deliberately.
- `type` — `inline | document` (see §5). Only on **promotable** types: outcome, requirement,
  component. Jobs and records are always their own document; risks are always inline — none
  carry `type`.
- `detail` — optional freeform markdown. Available on **every** element type. The universal
  extended-content field; replaces the old per-type `notes`. **It is the one place the
  "invalid is impossible" guarantee can't reach** — freeform can't be validated for "should this
  have been structured?" Its discipline is *soft*: guarded by a judgment topic ("`detail` is for
  narrative with no structured home; never a substitute for a structured field") plus a write-time
  nudge from the mutating verbs, not by a lint rule (§14).
- type-specific fields:
  - **requirement** → `acceptance_criteria` (list)
  - **component** → `responsibility` (required) + optional structured `data_model`,
    `interfaces`, `behavior`, `edge_cases`, `success_criteria` (flattened, no `body:` wrapper)
- edges — `mitigates`, `fulfills`, `affects`, `dependsOn`, `relationships` (see §4).

`principles` / `constraints` on the engineering root stay **bare string lists** — they are axioms
(givens), not choices-among-options, so they don't fit the ADR shape. (Caveat: as bare strings
they aren't graph-referenceable; promote them to keyed elements only if a real need to reference
them appears.) There is **no `technology_choices` field** — a technology choice *is* an ADR
(a choice among options, with consequences), so it lives in `engineering.decision_records`; any
"stack at a glance" view is **derived** from the tech-choice ADRs.

Records (`decision_records`, `change_records`) carry **no `kind`** — location implies the type
(product decision / arch decision / change), the same way position implies identity.

```yaml
product:
  name: Intent
  summary: "..."
  jobs:
    understand_the_rationale_behind_an_element:
      name: "Understand The Rationale Behind An Element"
      story: "When I need to understand why ... so I can be confident ..."
      outcomes:
        fast_rationale_lookup:
          name: "Fast Rationale Lookup"
          statement: "Minimize the time to trace ..."
          risks:
            ambiguous_or_premature_reference:
              name: "Ambiguous Or Premature Reference"
              statement: "An element referenced before its file exists ..."
          requirements:
            stable_logical_ids:
              name: "Stable Logical IDs"
              statement: "Intent must assign every element a stable logical ID ..."
              mitigates:
                - product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.risks.ambiguous_or_premature_reference
  decision_records:
    risk_over_odi_scoring:
      name: "Risk Over ODI Opportunity Scoring"
      affects: [product]                 # PDR → product-domain targets only
      summary: "..."
      context: "..."
      options: ["...", "..."]
      decision: "..."
      consequences: ["...", "..."]
```

> **Appendix A** (end of doc) gives the *maximal* expansion of every element type — every legal
> field populated. Real elements default to minimal: `name` + `statement`/`story` + required edges.

---

## 3. Identity & ordering

- **Key = identity = path segment = slug.** snake_case, format-linted.
- Uniqueness of keys is **structural** — dict keys can't repeat within a parent. Keys are
  only *locally* unique; position qualifies them.
- **Presentation order = insertion order** in the file. No `order:` field. Reordering is a
  block move. (Since the CLI is the sole writer and canonically re-serializes — §11 — insertion
  order is preserved and diffs stay local.)

---

## 4. References & the trace graph

**One address encoding — the dotted path — everywhere.** Every reference (`fulfills`,
`dependsOn`, `mitigates`, `affects`) is a **flat list of dotted-path strings**; the same dotted
path is what prose interpolation uses (§6). A whole-domain target is the bare root literal
(`product` / `engineering`). This collapses the earlier hybrid of nested maps + bare lists +
shaped coordinates into a single form.

- Why flat: nobody hand-authors the file (§Operating model), so the vertical compactness of nested
  maps bought nothing; a flat list gives **one self-contained, greppable line per edge**, so
  add/remove/re-parent is a one-line diff (no structural reshuffle), the schema and linter shrink,
  and dedup is a trivial set check. **Ancestor-grouping is a *derived view*** rendered by
  `intent show`/`trace`, not a storage shape — grouping is a mechanical dual, not declared intent.

```yaml
# a requirement
alignment_check_before_finalizing:
  mitigates:                              # same-outcome risk(s)
    - product.jobs.change_documentation_without_drift.outcomes.coherent_change.risks.linked_elements_fall_out_of_sync
  dependsOn:                              # requirement leaves, anywhere
    - product.jobs.change_documentation_without_drift.outcomes.coherent_change.requirements.mechanical_linter_first_pass
    - product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.stable_logical_ids
```

- **`relationships`** stays a **keyed map** (`<dotted-path>: note`) — it carries a *payload* (the
  note), so it can't collapse to a bare address list; the dotted-path key keeps addressing uniform
  and dedup structural:
  ```yaml
  relationships:
    engineering.components.reference_guides: "encodes the rules defined there as automated checks."
  ```
- **`affects`** (records) is a flat dotted-path list, root literals allowed:
  ```yaml
  affects:
    - engineering.components.bundle_linter
    - product.jobs.change_documentation_without_drift.outcomes.coherent_change
    - product
  ```

**Domain-scope rule (records).** A **PDR**'s `affects` targets must all resolve under `product.*`;
an **ADR**'s under `engineering.*` (linter-enforced, §9). A **CR**'s `affects` may mix domains and
levels freely — it is the cross-cutting record. A decision's cross-domain *impact* that must be
tracked is expressed either through the trace graph (a `fulfills`/`dependsOn` edge between elements)
or, if it's an actual modification to the other-domain element, as a **CR** — never as a
cross-domain `affects` on the decision itself.

**The trace graph is only what you declare here.** Prose links (§6) never create edges. Every
declared reference must resolve; a dangling target is `E001`.

---

## 5. Rendering: `type`, not tier

- **`type: inline | document`** — authored, binary. Default `inline` (minimal-first).
- **Directories are emergent:** a `document` with `document` children becomes a folder with an
  index file. Never authored.
- **Containment invariant — maintained by construction, not merely linted.** Document-ness must be
  connected from the root down: an `inline` element's entire subtree must be `inline` (an inline
  parent has no directory to host a document child's file). Rather than let you create the invalid
  state and error, the CLI keeps it valid by construction (consistent with §11's "make it
  impossible, not discouraged"):
  - **Promotion cascades up.** `intent promote <addr>` / `set <addr> type document` auto-promotes
    every inline ancestor in one call (like `mkdir -p`). The agent states the intent once.
  - **Demotion is the guarded direction.** `document → inline` is refused if the element has
    document descendants (it would orphan their files); `--cascade` demotes the subtree explicitly.
  - **`E002` survives only as a hand-edit backstop** — unreachable via the CLI, exactly like the
    slug/format/backlink checks the generator can't emit wrong (§9).

---

## 6. Interpolation, accessors & links

**One interpolation syntax, `{{ }}`,** disambiguated by prefix:

| Expression prefix | Kind | Address of |
|---|---|---|
| `paths.*` | config path variable | a directory/file (§7) |
| `product.` / `engineering.` / `change_records.` | absolute element address | an element |
| `.` (leading dot) | relative element address | an element in the current container |

**Accessors** — explicit, reserved as the terminal token. No positional inference:

| Accessor | Yields |
|---|---|
| `.name` | the element's `name` (text) |
| `.path` | the computed destination path (raw string) |
| `.link` | a **rendered markdown hyperlink** to the target; text = `name` |

`.link` is what most content uses — it re-derives the correct relative markdown href from the tree
layout, so links never break when the tree reorganizes (`mv`). Raw forms `.name` / `.path` exist
for hand-composed markdown (e.g. custom link text: `[the linter]({{.bundle_linter.path}})`).

```yaml
behavior: |
  Runs the ID/backlink pass, then the semantic pass; both encode the rules from
  {{.reference_guides.link}}. See the flow diagram: ![Linter flow]({{paths.assets}}/linter-flow.png)
  Cross-outcome checks touch
  {{product.jobs.change_documentation_without_drift.outcomes.coherent_change.requirements.alignment_check_before_finalizing.link}}.
```

- **Output is markdown-only.** There is no `link_style` switch and no `wiki` output — that was
  speculative generality (the committed docs are GitHub-flavored markdown reviewed in MRs). Wiki
  can return later as an *output adapter* if a real consumer appears; it is not carried now.
- **No `.image_link` accessor.** Accessors attach to element addresses, and no element is an image
  (assets never live in the model — §7). Embed assets with **plain markdown + path interpolation**:
  `![alt]({{paths.assets}}/x.png)`. (This dissolves the old "accessor on composed asset paths"
  open question — its premise is gone.)
- **No `[[ ]]` sugar.**
- **Links are cosmetic** — they render correctly and create **no** trace edges.
- Dangling prose addresses fail validation exactly like a dangling declared edge (`E001`).

---

## 7. Assets & config

- **Assets never live in the model.** Reference them by path variable: `{{paths.assets}}/x.png`.
- **`intent.config.yaml`** — optional; sane defaults apply when omitted:

  ```yaml
  output_dir: docs                 # where generated markdown is written
  paths:
    assets: docs/assets            # {{paths.assets}} resolves here
    # diagrams: design/diagrams    # user-defined named vars welcome
  ```

- `paths` is a **map of named path variables** resolved by **pure literal string
  substitution**. Deliberately **not a template language** — no conditionals, no logic, no
  arbitrary lookups, no injection surface.

---

## 8. Output tree & page anatomy

```
docs/product/          jobs / outcomes / … , drs/ (PDRs)
docs/engineering/      components , drs/ (ADRs)
docs/change-records/   top-level CRs
```

Decision records render under their domain's `drs/`; change records render to one top-level
directory.

**Generated docs are committed and CI-enforced.** The committed markdown *is* the human-facing
impact view (see §Operating model, §9). For that to work the generator must be **deterministic and
churn-free** — all derived content (backlinks, See-Also, maps) is emitted in a stable order sorted
by **canonical dotted address** (never discovery/insertion order), so re-running on unchanged intent
is byte-identical and adding element X touches only the pages that actually reference X.

**Fixed page anatomy** (top → bottom), so authored vs derived is visually separable in every diff:

1. **Generated banner** — `<!-- generated by intent from intent.yaml — do not edit -->`.
2. **Authored core** — the prose (`statement`/`story`/`responsibility`/`summary`), then
   type-specific structured fields (`acceptance_criteria`; a component's `behavior`/`edge_cases`/…),
   then `detail`. *This is where the author's/agent's diff lands.*
3. **Declared-edge section** — outgoing edges (Fulfills / Depends on / Mitigates / Relationships)
   rendered as `.link`s.
4. **`---` then a single derived footer** — "Referenced by" (backlinks/incoming), "See also", and
   maps. **Pure-derived, address-sorted, contiguous, at the bottom.** Section headings are
   canonical/fixed so their presence never itself churns.

The principle: **all pure-derived content lives in one trailing zone** → diff triage is instant
(body = intent, footer = ripple) and the "don't hand-edit" boundary is unambiguous.

---

## 9. Validation (the linter shrinks)

Most of the old 486-line linter evaporates:

- format / uniqueness / no-reuse → dict keys + snake_case regex
- "requirement without a risk" → `minItems: 1` on `mitigates`
- slug / path / tier / backlink correctness → the generator can't emit them wrong
- containment → maintained by construction via cascading promotion (§5)

**What the linter still does:**

1. **Reference resolution** — every reference resolves: structured edges *and* prose
   `{{ }}` / link addresses.
2. **Record domain-scope** — a PDR's `affects` targets resolve under `product.*`, an ADR's under
   `engineering.*` (CRs exempt). §4.
3. **No duplicate references** — within a leaf list, `mitigates`, or `affects`, the same target
   may not appear twice. (`relationships` keys dedup structurally.)
4. **Containment backstop** — the `inline`/`document` rule from §5, guarding against hand-edits
   (unreachable via the CLI).

Error codes (to be enumerated/stabilized as the validator is built — §15):

| Code | Check |
|---|---|
| `E001` | dangling reference (declared edge or prose address) |
| `E002` | inline/document containment (hand-edit backstop) |
| `E003` | duplicate reference in a list |
| `E004` | record cross-domain `affects` (PDR → non-product / ADR → non-engineering) |

(The old wiki name-collision code is gone with wiki — §6.) All of this runs as `intent validate`
(and automatically before every write). Failure messages end with `→ intent help E0NN` so the fix
is one hop away.

---

## 10. Declared vs derived

| Declared (intent) | Derived (mechanical dual) |
|---|---|
| elements, `name`, `statement`, `type` | backlinks, See-Also, maps |
| edges (`fulfills` / `mitigates` / `affects`) | directory structure, ancestor-grouped views |
| prose, explicit links | link hrefs, file paths, doc slugs, stack-summary view |

---

## 11. Interaction layer — the `intent` CLI

The intent tree is operated through a **read/write CLI (`intent`)**, not by hand-editing. It is the
read interface (progressive disclosure), the write interface (surgical, validated mutations), the
generator, and the validator — and it carries its own documentation and judgment content (§12).

**Delivery.** The CLI is a **single Go binary** — zero runtime prerequisite, fast cold start (it's
invoked many times per agent turn), trivial cross-platform release artifacts. Canonical re-serialize
on write (no comments in the model) via `gopkg.in/yaml.v3` with insertion order preserved
(= presentation order); since the CLI is the *sole* writer and always emits canonical form,
consecutive serializations diff cleanly by construction. *(This retires the earlier Python /
`uv` / `pyyaml` / PEP 723 choice.)*

**Distribution & pinning** (terraform-tooling model): `goreleaser` publishes per-platform binaries +
checksums to **GitHub Releases** (the single source of truth). Projects pin per-project via a
committed **`mise.toml`** using the **`ubi`** backend (`["ubi:you/intent"] = "1.4.2"`); `mise install`
fetches exactly that version, and multiple versions coexist in mise's store so switching projects
switches versions automatically (per-project selection, like `pyenv`/`nvm` — not a global
single-install). CI and pre-commit invoke the pinned binary deterministically via **`mise exec --
intent …`**. Additional ecosystem adapters (npm/uv/go tool, Nix flake) are thin wrappers over the
same release artifact, added **only when a real adopter needs one** (minimal-first). Nix/flake was
evaluated and dropped as the canonical channel — a heavier prerequisite than mise for no gain here.

**Addressing.** The canonical address is the dotted path; the CLI accepts the **shortest
unambiguous suffix** (`intent show stable_logical_ids`), lists candidates on ambiguity, and always
prints full addresses back. (Suffixes are input ergonomics only; anything persisted — prose,
scripts — uses the full path.)

| Group | Commands |
|---|---|
| Explore | `find <query>` (search name/statement/story/detail; `--type`/`--domain`; returns addresses + snippet) · `tree` (outline) · `show <addr>` (element + outgoing edges) · `trace <addr>` (neighborhood; `--depth N` / `--up` / `--down` / `--edges`) · `affects <addr>` (impact set / backlinks) |
| Write | `add <type> <parent> <key> [--field …]` · `set <addr> <field> <value>` · `link <kind> <from> <to>` / `unlink` (writes the storage shape for you) · `promote <addr>` / `set type` (**cascades ancestors** — §5) · `mv <addr> <new-parent \| --rename \| --before/--after>` (reorder / rename / re-parent, cascading references) · `rm <addr>` (guarded against dangling refs) · `record <cr \| pdr \| adr> …` |
| Validate / build | `validate` (§9; auto-runs pre-write) · `build [--out docs/]` (generate the markdown view) · `check` (build to temp + intent-aware diff vs on-disk docs; nonzero on drift — the CI/pre-commit drift gate) |
| Bootstrap / distribute | `init` (scaffold a new tree + config) · `install-skill --agent <target>` (render the skill from embedded content — §12) |

`show` / `affects` / `trace` are a deliberate **cost/scope gradient** (cheapest-sufficient wins):
`show` = "what is this?" (outgoing, cheap); `affects` = "what breaks if I change this?" (incoming);
`trace` = "give me the context" (full neighborhood, priced by depth).

Every write is validity-guaranteed, and edges are stated as relationships with the CLI writing the
storage shape — so nobody hand-builds an edge list or an ancestor grouping.

**Why CLI+Skills, not MCP or bare CLI.** Intent operates on *local files in a git repo*, not a
shared authenticated external system, so MCP's centralized auth/governance solves a problem we don't
have while taxing every turn with up-front schemas for the whole surface. Bare CLI is too loose (the
agent improvises malformed edges). CLI+Skills fits — and we harden it past the usual pattern:
reliability is pushed *into* the execution layer (`validate` rejects invalid writes; `link`/`mv`/
`promote` compute the hard shapes, so malformed calls are impossible, not merely discouraged), and
the *mechanics-coupled* instruction layer lives *inside* the binary (`intent help`, loaded on
demand), keeping token cost near bare-CLI and preventing skill-goes-stale drift. MCP could earn its
keep *later* only as a **read-side** query layer if intent trees become an org-wide, many-agent
*queried* asset — a complement over the same engine, not a replacement.

---

## 12. The skill, the help system & judgment

**The cleave: validity (CLI) vs. judgment (editorial content).** The litmus for where any piece of
content lives: *would it go stale when the schema/commands change, or only when the methodology
changes?*

- **Schema/command-coupled → the binary.** Execution logic (`validate`/`build`/mutations) plus the
  factual help that must track the tool: `ref:*` (field/edge/command/addressing tables), the
  `errors:E0NN` catalog, and command usage. Staleness here is fatal, so it's version-bound to the
  binary and served by `intent help`.
- **Methodology-coupled → judgment content.** requirement-vs-task, risk-vs-feature, altitude,
  materiality, writing-statements. This is editorial; it evolves on its own cadence.

**Judgment ships *with* the binary but the CLI never *makes* judgments.** The judgment content is
**embedded in the binary** (single source, co-versioned by construction, drift impossible) and
surfaced two ways over that one source:

1. **`intent help judgment:*`** — universal, on-demand; works for a human, CI, or a skill-less
   agent (guidance exists even with no skill attached).
2. **The skill is *generated* from that embedded content** via `intent install-skill` — so the
   always-loaded agent context and the CLI help are the same bytes, and can't diverge. Installed
   skill files are **derived, never hand-edited**.

The CLI remains a deterministic *carrier and server* of judgment, not a decider — `validate` still
passes anything well-formed. Its mutating verbs **proactively nudge**: `intent add requirement …`
footers the requirement-vs-task test (a pointer the loaded skill expands), so the agent meets
judgment at the moment of action — closing the "agent won't fetch what it doesn't know it needs" gap.

**The skill is a thin, generated wrapper.** `SKILL.md` shrinks to: the trigger `description`, one
paragraph on what Intent is, the entry-point instruction (*drive the `intent` CLI; run `intent
help`; `intent.yaml` is canonical; never hand-edit generated docs*), the 2–3 must-know gotchas, and
the concise judgment *tests* (one-liner per distinction) — full good/bad-pair deep-dives are pulled
on demand. `references/` and `templates/` as hand-maintained assets are gone: mechanics → the
generator; factual help → `intent help`; judgment → embedded + rendered.

**Install.** `mise` installs the CLI → `intent install-skill --agent <target>` renders the skill.
mise is the single entry point (no "skill installs the CLI" chicken-and-egg), and `install-skill` is
the **primary/only** skill-install path (the old `npx skills` route is retired — it can't guarantee
fresh content). The adapter contract emits a **per-agent file tree** (not a hardcoded single file),
so a nested `SKILL.md` + `references/` + sub-skills is possible if ever needed — though the default
stays a thin `SKILL.md` that is ~a pointer into `intent help`. Adapters: Claude Code's
`.claude/skills/` + a generic `AGENTS.md` fallback at launch; others follow.

**Help is plane-based** — organized by the *question being asked*. Topics use namespaced slugs.

| Plane | Purpose | Topics (launch set in **bold**) |
|---|---|---|
| `concept:*` | what things are & why | **`concept:<element>`** · `concept:model` · `concept:philosophy` |
| `ref:*` | factual / structural, tabular | **all of `ref:*`** (cheap, schema-derived): `ref:<element>` · `ref:edges` · `ref:addressing` · `ref:syntax` · `ref:config` · `ref:commands` |
| `judgment:*` | good, not just valid | **`risk-vs-feature`, `coherence`, `materiality`, `writing-statements`, `which-domain-owns-this-decision`** · then `outcome-vs-solution`, `requirement-vs-task`, `job-vs-activity`, `altitude`, `split-or-merge`, `restraint`, `worth-recording`, `naming`, `red-flags` |
| `guide:*` | task-routed recipes (goal→commands) | **`edit`, `trace`, `add-requirement`, `record-decision`, `draft-judge-review`** · then `new-tree`, `move`, `promote`, `build` |
| `errors` | recovery keystone | **the full `E0NN` catalog** — cause + exact fix; reached inline from failures |

Every `judgment:*` topic uses a fixed shape: *question → one-line test → good/bad pairs → failure
mode.* Worked examples are drawn from Intent's own dogfood tree.

**Discoverability layer** (kept lean — the consumer is an agent that *retrieves*, not browses):
`intent help` index task-first; `intent help --list` (flat, greppable) and `--json`; guessable
namespaced slugs; command-output footers that suggest the next command and the relevant judgment
topic. Heavier human-navigation polish (curated cross-links) is deferred until a real need appears.

**Launch scope (minimal-first):** ship the bold topics — the high-traffic core plus the full, cheap
`ref:*` — and add the rest as real gaps appear.

---

## 13. Reset impact

- **Survives (as content):** the element rubric and the draft/judge/review discipline — now
  `concept:<element>` + `judgment:*` and `guide:draft-judge-review`, **embedded in the binary** and
  rendered into the skill. Judgment stays human/LLM.
- **Absorbed into the binary:** factual `references/` → `intent help`; `context-tracing.md` →
  `intent trace`; judgment content → embedded + rendered to the skill; most of
  `lint_intent_tree.py` → `intent validate`; templates → the generator.
- **Gone:** hand-maintained `templates/` and `structure.md` — the generator owns paths / tiers /
  naming / trees / child-rendering. Wiki output, the `.image_link` accessor, `technology_choices`,
  and the CR `affects_narrative`/`summary` fields.
- **Newly stale (dogfood):** Intent's *own* `intent.yaml` still describes the old design and must be
  re-authored to describe this one.

---

## 14. Costs accepted with eyes open

- The CLI (engine + generator + validator + embedded help/judgment) must exist and run — no
  graceful hand-write-a-markdown-file degradation.
- **`mise` is a prerequisite** for consumers/CI (a single small binary, no daemon — far lighter than
  Nix, which we evaluated and dropped). The CLI itself is a zero-runtime Go binary.
- Authoring the CLI in Go is more verbose than Python — a one-time cost bought back by a permanent
  runtime + distribution win.
- Dotted-path edges are long/repetitive in the raw file — accepted, since nobody hand-authors it and
  each edge is a clean one-line diff.
- Committed generated docs double the PR surface (intent + docs) — accepted, because the docs diff
  *is* the impact view; determinism keeps it low-noise.
- **`detail` discipline is a soft guarantee** — the one place "invalid is impossible" can't reach,
  guarded by judgment + write-time nudge only (§2).

---

## 15. Open questions (to work together)

- **Error-code catalog.** Enumerate and stabilize the `E0NN` codes (§9) as the validator is built.
- **Committed vs. installed skill files.** Installed per-agent skill files are derived and
  regenerated (§12); whether they're also committed for browsability (parallel to the committed-docs
  choice) is open.
- **Additional distribution adapters.** Which ecosystem wrappers (npm/uv/go tool/Nix) get added, and
  when — driven by real adopter demand.

*(Resolved this pass: committed-vs-ephemeral docs → **committed + CI `intent check`**; CLI
distribution channel → **Go binary + mise `ubi`**; accessor-on-asset-path → **dissolved** with
`.image_link`; skill install targets → **`intent install-skill`**, adapter emits a file tree.)*

---

## 16. Decision log

1. **Three files, one schema** — sharded on root boundaries. *(reversed by #16)*
2. **Roots** = `product`, `engineering`, `change_records`.
3. **Decision records nest by domain; change records are top-level.** *(refined by #31 — decisions
   are domain-pure/linter-scoped)*
4. **Key = identity**, snake_case; no `id`; `name` required; location implies record type.
5. **Insertion order** is presentation order.
6. **Hybrid references** (nested maps + bare lists + shaped coordinates). *(reversed by #23)*
7. **`type: inline | document`**, directories emergent. *(refined by #30 — containment by
   construction)*
8. **`body` flattened** onto the element.
9. **`{{ }}` interpolation** with `paths.*` / absolute / relative addresses.
10. **Four accessors** + no image sniffing. *(reversed by #29 — three accessors, `.image_link` cut)*
11. **Links are cosmetic** — the graph is declared edges only.
12. **`link_style: markdown | wiki`**. *(reversed by #28 — markdown-only)*
13. **Assets externalized**; `intent.config.yaml`; named path vars; not a template language.
14. **Linter** = reference resolution + containment + name uniqueness (wiki). *(name-uniqueness
    dropped by #28; recomposed in #31)*
15. **Universal `detail`** field (replaces `notes`); **requirement** adds `acceptance_criteria`;
    **component** keeps structured fields. *(discipline clarified as soft — #36)*
16. **Single `intent.yaml`** (reverses #1).
17. **Read/write `intent` CLI** as the interaction layer; canonical re-serialize; shortest-suffix
    addressing; `mv` cascades. *(the Python / `uv` / `pyyaml` delivery is reversed by #25)*
18. **Thin skill wrapper**; `references/`+`templates/` absorbed. *(refined by #26/#27 — judgment
    stays editorial but ships embedded in the binary and renders into the skill)*
19. **Plane-based help** + a lean discoverability layer; authored minimal-first.
20. **Interface = CLI+Skills** (not MCP / bare-CLI); reliability pushed into the CLI; MCP reserved as
    a possible future read-side layer. *(refined by #26)*
21. **Committed, self-bootstrapping repo bundle, two install paths.** *(superseded by #27 —
    generated by `install-skill`; `npx skills` retired; mise bootstraps the CLI)*
22. **Operating model** — agent writes via the CLI; humans review the MR diff; hand-authoring is not
    a design driver; **clean MR diffs are an explicit design goal**.
23. **One address encoding** — all edges + prose are flat **dotted-path** lists (reverses #6);
    `relationships` stays a keyed `path: note` map; ancestor-grouping is a derived render view.
24. **Generated docs committed + CI `intent check`** drift gate; deterministic, address-sorted,
    churn-free generation; fixed page anatomy with a single derived footer zone.
25. **Distribution** — single **Go binary**, `goreleaser` → GitHub Releases, pinned per-project via
    **`mise` (`ubi`)**, `mise exec` for CI/agents (retires #17's Python/`uv`/`pyyaml`); Nix/flake
    evaluated and dropped as canonical; other adapters added on demand.
26. **Validity/judgment cleave** — CLI owns deterministic execution + schema-coupled help; judgment
    is editorial content **embedded in the binary**, served universally by `intent help`, and the
    CLI carries but never *enforces* judgment (nudges via footers).
27. **Skill generated by `intent install-skill`** from embedded content (primary/only path; `npx
    skills` retired); adapter emits a per-agent file tree; installed files derived.
28. **Markdown-only output** — cut `link_style`/wiki; names never require uniqueness; wiki-collision
    lint removed. Wiki returns only as a future output adapter.
29. **Three accessors** (`.name` / `.path` / `.link`) — cut `.image_link`; assets embedded via raw
    markdown; asset-path open question dissolved.
30. **Containment by construction** — `promote`/`set type` cascades ancestors; demotion guarded
    (`--cascade`); `E002` becomes a hand-edit backstop.
31. **Records domain-dogmatic** — PDR `affects` ⊆ `product.*`, ADR `affects` ⊆ `engineering.*`
    (linter-enforced, new `E004`); cross-domain intent via the trace graph; cross-domain effects
    routed to CRs.
32. **No `technology_choices` field** — tech choices are ADRs; stack summary derived;
    `principles`/`constraints` stay bare lists.
33. **CR trimmed** — drop `affects_narrative` and `summary` (`change` + `rationale` are the required
    what/why); DR keeps required `summary`; list blurbs derived.
34. **`intent find`** added (search → addresses); `show`/`affects`/`trace` kept as a documented
    cost/scope gradient (`trace` gains `--depth`/`--up`/`--down`/`--edges`).
35. **Distinct prose-field names kept** (`summary`/`story`/`statement`/`responsibility`) as semantic
    typing / author prompt; shapes pinned in `ref:*` + `judgment:writing-statements`.
36. **`detail` stays universal**, guarded by judgment + write-time nudge; its discipline is
    explicitly a *soft* (non-structural) guarantee.

---

## Appendix A — Maximal expansion of every element type

Every legal field populated, `required` / `optional` marked. Real elements default to minimal
(`name` + `statement`/`story` + required edges; `type` defaults to `inline`). Any prose field
may contain `{{ }}` interpolation (§6). All edges are flat dotted-path lists (§4).

### Product side

```yaml
# ─── product (root) ───────────────────────────────────────────────
product:
  name: "Intent"                     # required
  summary: "..."                     # required
  detail: "..."                      # optional — freeform markdown
  jobs: { … }                        # map
  decision_records: { … }            # map (PDRs)

# ─── job ──────────────────────────────────────────────────────────
change_documentation_without_drift:
  name: "Change Documentation Without Drift"          # required
  story: "When I'm editing docs I want ... so ...."   # required
  detail: "..."                      # optional
  outcomes: { … }                    # map
  # no `type`

# ─── outcome ──────────────────────────────────────────────────────
coherent_change:
  name: "Coherent Change"            # required
  type: document                     # optional — inline | document (default inline)
  statement: "Minimize the likelihood that ..."   # required
  detail: "..."                      # optional
  risks: { … }                       # map
  requirements: { … }                # map

# ─── risk ─────────────────────────────────────────────────────────
linked_elements_fall_out_of_sync:
  name: "Linked Elements Fall Out Of Sync"   # required
  statement: "Frequent AI-assisted edits ..."   # required
  detail: "..."                      # optional
  # no `type`, no edges

# ─── requirement ──────────────────────────────────────────────────
alignment_check_before_finalizing:
  name: "Alignment Check Before Finalizing"   # required
  type: document                     # optional
  statement: "Intent must validate every change ..."   # required
  mitigates:                         # required, ≥1 — flat dotted-path list (same-outcome risks)
    - product.jobs.change_documentation_without_drift.outcomes.coherent_change.risks.linked_elements_fall_out_of_sync
  dependsOn:                         # optional — flat dotted-path list (requirement leaves, anywhere)
    - product.jobs.change_documentation_without_drift.outcomes.coherent_change.requirements.mechanical_linter_first_pass
    - product.jobs.understand_the_rationale_behind_an_element.outcomes.fast_rationale_lookup.requirements.stable_logical_ids
  acceptance_criteria:               # optional — list
    - "..."
  detail: "..."                      # optional
```

### Engineering side

```yaml
# ─── engineering (root) ───────────────────────────────────────────
engineering:
  name: "Intent"                     # required
  summary: "..."                     # required
  detail: "..."                      # optional
  principles: ["..."]                # optional — bare list (axioms)
  constraints: ["..."]               # optional — bare list (axioms)
  components: { … }                  # map
  decision_records: { … }            # map (ADRs — technology choices live here)

# ─── component ────────────────────────────────────────────────────
bundle_linter:
  name: "Bundle Linter"              # required
  type: document                     # optional
  responsibility: "Runs a mechanical first pass over a docs/ tree."   # required
  fulfills:                          # required, ≥1 — flat dotted-path list
    - product.jobs.change_documentation_without_drift.outcomes.coherent_change.requirements.mechanical_linter_first_pass
  relationships:                     # optional — keyed map, <dotted-path>: note
    engineering.components.reference_guides: "encodes the rules defined there as automated checks."
  data_model: "..."                  # optional — structured detail
  interfaces: "..."                  # optional
  behavior: "Runs the ID pass, then the semantic pass. See ![flow]({{paths.assets}}/linter-flow.png) and {{.reference_guides.link}}"   # optional
  edge_cases: ["..."]                # optional — list
  success_criteria: ["..."]          # optional — list
  detail: "..."                      # optional — freeform (was `notes`)
```

### Records

```yaml
# ─── decision record (PDR in product / ADR in engineering) ────────
risk_over_odi_scoring:
  name: "Risk Over ODI Opportunity Scoring"   # required
  affects:                           # required, ≥1 — flat dotted-path list, DOMAIN-PURE
    - product                        #   PDR → product.* only (ADR → engineering.* only)
  summary: "..."                     # required — the DR's core statement
  context: "..."                     # optional
  options: ["...", "..."]            # optional — list
  decision: "..."                    # optional
  consequences: ["...", "..."]       # optional — list
  # no `kind`, no `type`, no `detail`

# ─── change record (top-level change_records) ─────────────────────
add_relationships_to_components:
  name: "Add Relationships To Components"   # required
  affects:                           # required, ≥1 — flat dotted-path list, MIXED domains/levels
    - engineering.components.bundle_linter
    - product.jobs.change_documentation_without_drift.outcomes.coherent_change.requirements.alignment_check_before_finalizing
    - product                                     # a root literal
  change: "What changed."            # required
  rationale: "Why it changed."       # required
  # no `kind`, no `type`, no `summary`, no `affects_narrative`
```
