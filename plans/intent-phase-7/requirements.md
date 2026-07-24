---
plan_slug: intent-phase-7
phase: requirements
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
---

# Requirements: Intent Phase 7 — Skill renderer & install-skill

> Part of the Intent Phases 7–9 build (BUILD_PLAN). Split from the former
> `intent-build` plan into one plan per phase. Execution order across the split
> plans: **Phase 7 → Phase 8 → Phase 9**, human merge gate between each.

## Problem Statement

Phases 0–6 of the Intent CLI are merged to `main`. The tool can read, validate,
build, mutate, and serve embedded help — but it cannot yet render an agent skill
from that embedded content (`install-skill`). Without it, always-loaded agent
context and `intent help` can silently diverge.

## Solution

Complete BUILD_PLAN Phase 7 as one merge-sized PR via the Gas City
`build-from-convoy` workflow. Phase 6 is treated as done (merged commit
`5ceb499`) and is the sole source for skill rendering. Phase 7 is the
calibration slice — crisp test gates, `max_iterations=4` on the first sling.

## User Stories

### Phase 7 — Skill renderer

**As an** agent user with the Intent CLI installed,
**I want** `intent install-skill --agent claude-code` to render a thin skill file
from the embedded help registry,
**so that** always-loaded agent context and `intent help` cannot diverge.

Acceptance criteria:

- An adapter contract maps embedded topics → a per-agent file tree.
- Claude Code adapter writes `.claude/skills/intent/SKILL.md` (already gitignored).
- A generic `AGENTS.md` fallback adapter exists for non-Claude targets.
- Generated `SKILL.md` contains: trigger description, one Intent paragraph,
  entry-point instruction, 2–3 gotchas, judgment one-liners (full deep-dives
  remain on-demand via `intent help`).
- Tests prove rendered skill content is derived from `internal/help` (no
  hand-maintained duplicate prose).
- `go test ./...` passes.

## Out Of Scope

- Distribution/CI tooling (`goreleaser`, `mise`, GitHub Actions) — Phase 8.
- Dogfood `intent.yaml` and `archive/` retirement — Phase 9.
- Additional help topics deferred in DESIGN §12 unless a real gap appears.
- Auto-merging PRs from Gas City workflows.
- Rewriting Phases 0–6 or changing the Phase 6 help content launch set.

## Other Notes

- **Global writing rule** (BUILD_PLAN): concise, plain language a junior engineer
  can follow; re-derive rather than paste old archive wording.
- **Salvage rule:** bucket 2 (judgment) already landed in Phase 6; bucket 3
  (mechanical) discarded. Bucket 1 (element intent) is Phase 9 only.
- **GC cadence:** one formula run → one PR → human review + merge → next phase.
  See `plans/intent-phase-7/implementation-plan.md` for convoy boundaries and
  drain policy.
