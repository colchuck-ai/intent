---
plan_slug: intent-phase-9
phase: implementation-plan
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
requirements_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-9/requirements.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
---

# Implementation Plan: Intent Phase 9 — Dogfood tree & finalize

## Summary

Author the dogfood tree via the write CLI, commit generated docs, delete
`archive/`, and rewrite `README.md`. One Gas City convoy (`build-from-convoy`),
one PR. Depends on Phases 7 and 8 merged.

## Current System

| Area | Location | Notes |
|------|----------|-------|
| Generator | `internal/gen/` | `build` / `check` drift gate (Phase 3) |
| Mutations | `internal/mutate/`, `internal/cli/write.go` | Full write path (Phases 4–5) |
| Seed fixture | `internal/model/testdata/seed.intent.yaml` | Test-only; not repo dogfood |
| Archive | `archive/` | Pre-rewrite Python skill, prototype, old docs |
| README | `README.md` | Still describes retired Agent Skill / markdown workflow |

Module: `github.com/colchuck-ai/intent`, Go 1.25, cobra + yaml.v3.

## Proposed Implementation

### Convoy boundary: Phase 9 — Dogfood tree & finalize

**Goal:** Self-describing repo; archive gone; README accurate.

**Authoring rules for workers:**

- Use write CLI only (`intent add`, `set`, `link`, `record`, `promote`, `mv`).
- Reference `DESIGN.md`, `BUILD_PLAN.md`, and `archive/` for salvage — re-derive prose.
- Run `intent validate` after each major subtree; `intent build` before commit.

**Target tree shape (high level):**

```
intent.yaml                 # product + engineering roots
intent.config.yaml          # output_dir, paths (may exist from init scaffold)
docs/product/...            # generated
docs/engineering/...        # generated
docs/change-records/...     # generated (per gen/path.go conventions)
```

Content should cover: jobs/outcomes/risks/requirements for Intent the product,
components for CLI subsystems (help, gen, validate, mutate), PDRs/ADRs/CRs for
major design choices (Go rewrite, embedded help, committed docs + check gate).

**Cleanup:**

- Delete entire `archive/` directory.
- Rewrite `README.md`: Go binary, mise install, install-skill, validate/build/check workflow.
- Final pass on `internal/help/content/errors/E00N.md` if dogfood exposes new edge cases.

**Drain policy for GC:** `same-session` for tree-authoring beads; final cleanup
bead depends on docs bead.

**Convoy beads (see `plans/intent-phase-9/tasks.md`):**

1. `dogfood-scaffold` — `intent init` or manual scaffold + product spine (jobs, outcomes, risks, requirements)
2. `dogfood-engineering` — architecture + components for CLI subsystems
3. `dogfood-records` — PDRs/ADRs/CRs for major decisions
4. `dogfood-build-commit` — `intent build`, commit docs, `intent check` green
5. `dogfood-finalize` — delete `archive/`, update README, stabilize error catalog, update CI to root `intent.yaml`

### Gas City execution

Artifact root: `plans/intent-phase-9/` (canonical per-plan-slug layout; build
outputs under `plans/intent-phase-9/build/`).

1. Mayor writes `plans/intent-phase-9/tasks.md` + bead payload.
2. Dry-run + create beads via `create_beads_from_tasks.py`.
3. Sling:

```bash
gc sling gc.run-operator <phase-9-convoy-id> --on build-from-convoy \
  --var artifact_root=plans/intent-phase-9/build \
  --var requirements_path=plans/intent-phase-9/requirements.md \
  --var plan_path=plans/intent-phase-9/implementation-plan.md \
  --var plan_review_path=plans/intent-phase-9/plan-review.md \
  --var decomposition_path=plans/intent-phase-9/tasks.md \
  --var interaction_mode=interactive \
  --var review_mode=agent \
  --var drain_policy=same-session \
  --var open_pr=true
```

4. Human review PR + CI → merge.

## Testing

| Phase | Verification |
|-------|----------------|
| 9 | `intent validate`; `intent check`; `go test ./...`; README accuracy review |

## Rollout

1. Approve requirements + implementation plan + plan review.
2. Decompose Phase 9 → beads → sling → merge.

No feature flags. The phase is independently revertable via git revert of its PR.

## Open Questions

1. **Phase 9 bead 5 CI update** — must explicitly move CI `check` from the seed
   fixture to the root `intent.yaml`.
2. **Rig-scoped gascity roles** — verify worker routing before the sling (`gc rig status intent`).
