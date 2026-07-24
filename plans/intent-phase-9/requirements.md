---
plan_slug: intent-phase-9
phase: requirements
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
---

# Requirements: Intent Phase 9 — Dogfood tree & finalize

> Part of the Intent Phases 7–9 build (BUILD_PLAN). Split from the former
> `intent-build` plan into one plan per phase. Runs last, after Phase 8 merges.

## Problem Statement

Intent does not yet describe itself in its own format. The repo still ships a
pre-rewrite `archive/` and a stale `README.md`, so it does not demonstrate the
tool it ships.

## Solution

Complete BUILD_PLAN Phase 9 as one merge-sized PR via the Gas City
`build-from-convoy` workflow: author the dogfood `intent.yaml` via the write CLI,
commit generated docs, delete `archive/`, and rewrite `README.md`. Depends on
Phases 7 and 8 being merged.

## User Stories

### Phase 9 — Dogfood tree & finalize

**As a** new adopter reading this repo,
**I want** Intent to describe itself in `intent.yaml` with generated docs committed,
**so that** the repo demonstrates the tool it ships.

Acceptance criteria:

- Root `intent.yaml` (and `intent.config.yaml` if needed) authored **only** via
  write CLI commands (`add`, `set`, `link`, `record`, etc.) — no hand-editing.
- Prose is re-derived from `DESIGN.md` / `BUILD_PLAN.md` with plain language;
  salvage bucket 1 from `archive/` is reference only.
- `intent validate` passes; `intent build` produces committed docs under
  `docs/product`, `docs/engineering`, and change-record paths per generator config.
- `intent check` is green with committed docs.
- `archive/` is deleted.
- `README.md` describes the Go CLI world (not the retired Agent Skill / npx flow).
- Error catalog topics in `internal/help/content/errors/` are finalized and stable.
- `go test ./...` passes.

## Out Of Scope

- Skill renderer / `install-skill` — Phase 7.
- Distribution/CI tooling — Phase 8 (this phase only updates CI from the seed
  fixture to the root `intent.yaml`).
- Additional help topics deferred in DESIGN §12 unless a real gap appears during
  dogfood.

## Other Notes

- **Global writing rule** (BUILD_PLAN): concise, plain language a junior engineer
  can follow; re-derive rather than paste old archive wording.
- **Salvage rule:** bucket 1 (element intent) → dogfood tree, reference only;
  bucket 2 (judgment) already landed in Phase 6; bucket 3 (mechanical) discarded.
- **GC cadence:** one formula run → one PR → human review + merge.
  See `plans/intent-phase-9/implementation-plan.md` for convoy boundaries and
  drain policy.
