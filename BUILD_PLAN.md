# Intent — Build Plan

The plan for building Intent as specified in `DESIGN.md`. Work is split into ten phases.
Each phase is sized for a single working session, leaves the repo building and tested, and
touches one subsystem at a time so context stays manageable.

**Global writing rule (applies to all prose we author — tree, help, judgment, docs):**
concise, plain language a junior engineer can follow. When unsure whether to reuse old
wording, re-derive it fresh.

**Salvage rule (three buckets, bias toward re-deriving):**
- *Element intent* prose → the dogfood `intent.yaml` (Phase 9).
- *Methodology / judgment* prose → binary-embedded help topics (Phase 6).
- *Mechanical / factual* prose → discarded; the tool regenerates it.

**Tech choices (implementation detail):** single Go binary, `cobra` for the command tree,
`gopkg.in/yaml.v3` for canonical load/serialize with insertion order preserved.

---

## Phase 0 — Foundation & archival

**Goal:** a buildable Go project, the schema, and a rock-solid model round-trip.

**Build:**
- Move old material (`prototype/`, old `intent/`, old `docs/`, `install.sh`, `Makefile`,
  old schema) into `archive/` in one commit. Keep it greppable during the rewrite; deleted
  in Phase 9.
- Scaffold the Go module + `cobra` skeleton (`intent` builds and prints help).
- Write `intent.schema.json` (the contract) and the Go structs that mirror it.
- Canonical load → in-memory model → serialize, preserving insertion order.
- Hand-author a *tiny* seed `intent.yaml` fixture. This is the one allowed hand-edit (the
  bootstrap case); everything downstream tests against it.

**Done when:** `intent` builds; loading the seed and re-serializing is byte-identical;
the seed validates against the schema.

---

## Phase 1 — Addressing & read commands

**Goal:** find and read anything in the tree. (depends: Phase 0)

**Build:**
- Dotted-path addressing: resolve full paths, accept the shortest unambiguous suffix, list
  candidates on ambiguity, always print the full address back.
- Build the element index and edge graph in memory.
- Read commands: `tree`, `show <addr>`, `find <query>` (`--type`/`--domain`),
  `trace <addr>` (`--depth`/`--up`/`--down`/`--edges`), `affects <addr>`.

**Done when:** all read commands work against the seed; suffix resolution and ambiguity
are tested.

---

## Phase 2 — Validation

**Goal:** catch every invalid state the CLI must reject. (depends: Phase 1)

**Build:**
- `E001` dangling reference — declared edges *and* prose `{{ }}` addresses.
- `E002` inline/document containment (hand-edit backstop).
- `E003` duplicate reference in a list.
- `E004` record cross-domain `affects` (PDR → non-product / ADR → non-engineering).
- `intent validate`; every failure message ends with `→ intent help E0NN`.

**Done when:** validate passes the clean seed; each `E0NN` has a failing fixture and a test.

---

## Phase 3 — Generation (`build`) & drift gate (`check`)

**Goal:** turn `intent.yaml` into the committed markdown docs, deterministically.
(depends: Phase 1; Phase 2 recommended)

**Build:**
- Interpolation + accessors in prose: `.name` / `.path` / `.link`, `paths.*`, absolute and
  relative element addresses.
- Fixed page anatomy: banner → authored core → declared-edge section → `---` → single
  derived footer.
- Derived content (backlinks / "Referenced by", See-also, maps), address-sorted so output
  is stable and low-churn.
- Output tree (`docs/product`, `docs/engineering`, `docs/change-records`, `drs/`).
- `intent.config.yaml` loading (`output_dir`, `paths`).
- `build [--out]`; `check` (build to temp, intent-aware diff, nonzero on drift).

**Done when:** building the seed produces docs; re-running is byte-identical; `check` is
green; adding one reference touches only the pages that reference it.

*Heaviest phase. If a session runs long, split: (3a) interpolation + page anatomy,
(3b) derived footer + `check`.*

---

## Phase 4 — Write path: simple mutations

**Goal:** create and edit elements safely. (depends: Phases 2, 1)

**Build:**
- Canonical re-serialize on every write; run `validate` before committing any write.
- `add <type> <parent> <key> [--field]`, `set <addr> <field> <value>`,
  `rm <addr>` (refuse if it would dangle a reference), `link`/`unlink` (writes the storage
  shape so nobody hand-builds an edge list).
- The `detail` write-time nudge: if a written `detail` contains a resolvable `{{ }}`
  address with no matching edge, footer a suggestion to declare the edge. No lint, nothing
  in CI.
- Scaffold command-output footers (next command + judgment topic) — slugs filled in Phase 6.

**Done when:** mutations produce canonical, valid files with local diffs; the nudge fires;
guarded `rm` is tested.

---

## Phase 5 — Write path: structural mutations

**Goal:** the moves that reshape the tree, valid by construction. (depends: Phase 4)

**Build:**
- `promote` / `set type document` with ancestor cascade (like `mkdir -p`).
- Demotion guarded; `--cascade` demotes a subtree explicitly.
- `mv` (reorder `--before`/`--after`, `--rename`, re-parent) with reference cascade across
  the whole tree.
- `record <cr|pdr|adr>`.

**Done when:** cascades keep containment valid without any error path; `mv` rewrites every
referring edge; results round-trip cleanly.

---

## Phase 6 — Help & judgment content system

**Goal:** the embedded, on-demand instruction layer. (depends: Phase 5)

**Build:**
- Embedded content mechanism (Go `embed`) keyed by namespaced slugs.
- `intent help` (task-first index), `--list`, `--json`, `intent help <slug>`,
  `intent help E0NN`.
- Wire the Phase 4/5 footers to real topics.
- Author the launch set (the bold topics in DESIGN §12): `concept:<element>`, all `ref:*`,
  the bold `judgment:*`, bold `guide:*`, the full `errors` catalog. Salvage bucket 2, bias
  to re-derive, plain/junior-friendly.

**Done when:** every footer and error points to a real topic; `--json` output is stable;
the content reads plainly.

*Content-heavy — this is more writing than coding. Keep it its own session.*

---

## Phase 7 — Skill renderer & `install-skill`

**Goal:** generate the agent skill from embedded content. (depends: Phase 6)

**Build:**
- Adapter contract: emit a per-agent file tree from the embedded content.
- Claude Code adapter (`.claude/skills/…`) + generic `AGENTS.md` fallback.
- Thin generated `SKILL.md`: trigger description, one paragraph on what Intent is, the
  entry-point instruction, 2–3 gotchas, judgment one-liners.

**Done when:** `intent install-skill --agent claude-code` renders; the skill is the same
bytes as the embedded content; installed files are gitignored (never committed).

---

## Phase 8 — Distribution & CI

**Goal:** ship and pin the binary. (depends: Phase 3 for `check`)

**Build:**
- `goreleaser` → per-platform binaries + checksums → GitHub Releases.
- `mise.toml` (`ubi` backend) pins per project; `mise exec -- intent …` for CI/agents.
- CI runs `intent validate` + `intent check` (the drift gate); pre-commit hook.

**Done when:** the release artifact builds; `mise install` fetches the pinned version;
CI is green on this repo.

---

## Phase 9 — Dogfood tree & finalize

**Goal:** Intent describes itself, in its own format. (depends: Phases 5, 3, 6)

**Build:**
- Author Intent's own `intent.yaml` describing this design, using the write CLI. Salvage
  bucket 1, bias to re-derive, plain language.
- Generate and commit the docs; wire `check` into CI.
- Delete `archive/`; update `README.md`; finalize/stabilize the `E0NN` catalog in `errors`.

**Done when:** the dogfood tree validates and builds; docs are committed; `archive/` is
gone; the README describes the new world.

---

## Dependency map

```
P0 ─┬─ P1 ─┬─ P2 ─┬─ P4 ── P5 ─┬─ P6 ── P7
    │      │      │            │
    │      └─ P3 ─┴─ P8        └──────────── P9
    │            (P3 also → P9)
    └ (P3, P9 consume P0's model)
```

Suggested order: **0 → 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8 → 9.**
Phase 8 can move earlier if you want CI wired sooner; everything else follows the arrows.
