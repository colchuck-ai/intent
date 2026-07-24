---
plan_slug: intent-phase-8
phase: requirements
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
---

# Requirements: Intent Phase 8 — Distribution & CI

> Part of the Intent Phases 7–9 build (BUILD_PLAN). Split from the former
> `intent-build` plan into one plan per phase. Runs after Phase 7 merges.

## Problem Statement

The Intent CLI has no release or drift-enforcement loop: no `.github/`, no
`goreleaser`, no `mise.toml`. CI and agents cannot invoke a known, pinned Intent
version, and doc drift is not rejected automatically.

## Solution

Complete BUILD_PLAN Phase 8 as one merge-sized PR via the Gas City
`build-from-convoy` workflow: pinned releases plus automated `validate` + `check`
drift gates. Depends on Phase 7 being merged.

## User Stories

### Phase 8 — Distribution & CI

**As a** project maintainer,
**I want** pinned releases and automated drift checks,
**so that** CI and agents invoke a known Intent version and reject doc drift.

Acceptance criteria:

- `.goreleaser.yaml` builds per-platform binaries + checksums for GitHub Releases.
- `mise.toml` pins the binary via the `ubi` backend (`DESIGN.md` §11).
- GitHub Actions workflow runs `intent validate` and `intent check` on PR/push.
- A pre-commit hook runs the drift gate locally.
- CI is green on `main` using `mise exec -- intent …` (or equivalent pinned invoke).
- `go test ./...` still passes.

## Out Of Scope

- Skill renderer / `install-skill` — Phase 7.
- Dogfood `intent.yaml` and `archive/` retirement — Phase 9.
- npm/uv/Nix distribution adapters (DESIGN §15 open question).
- Auto-merging PRs from Gas City workflows.

## Other Notes

- **Global writing rule** (BUILD_PLAN): concise, plain language a junior engineer
  can follow; re-derive rather than paste old archive wording.
- **GC cadence:** one formula run → one PR → human review + merge → next phase.
  See `plans/intent-phase-8/implementation-plan.md` for convoy boundaries and
  drain policy.
- Phase 8 could have run before Phase 7 per BUILD_PLAN; we keep 7 → 8 → 9 so
  Phase 7 stays the calibration slice with crisp test gates.
