---
plan_slug: intent-phase-9
phase: implementation-plan
requirements_file: plans/intent-phase-9/requirements.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-27T00:00:00Z
---

# Implementation Plan: Intent Phase 9 — Dogfood tree & finalize

## Summary

Author the dogfood tree via the write CLI, commit generated docs, delete
`archive/`, and rewrite `README.md`. One working session, one PR. Depends on
Phases 7 and 8 merged.

## Current System

| Area | Location | Notes |
|------|----------|-------|
| Generator | `internal/gen/` | `build` / `check` drift gate (Phase 3) |
| Mutations | `internal/mutate/`, `internal/cli/` | Full write path (Phases 4–5) |
| Help | `internal/help/` | Embedded topics + error catalog (Phase 6) |
| Seed fixture | `internal/model/testdata/seed.intent.yaml` | Test-only; not repo dogfood |
| Archive | `archive/` | Pre-rewrite Python skill, prototype, old docs |
| README | `README.md` | Still describes retired Agent Skill / markdown workflow |

Module: `github.com/colchuck-ai/intent`, Go 1.25, cobra + yaml.v3.

## Proposed Implementation

**Goal:** Self-describing repo; archive gone; README accurate.

### Authoring rules

- Use the write CLI only (`intent add`, `set`, `link`, `record`, `promote`, `mv`).
  The dogfood tree is the proof that the write path works — hand-editing
  `intent.yaml` defeats the point of the phase.
- Reference `DESIGN.md`, `BUILD_PLAN.md`, and `archive/` for salvage, but
  re-derive the prose.
- Run `intent validate` after each major subtree; `intent build` before commit.

### Target tree shape

```
intent.yaml                 # product + engineering roots
intent.config.yaml          # output_dir, paths (may exist from init scaffold)
docs/product/...            # generated
docs/engineering/...        # generated
docs/change-records/...     # generated (per gen/path.go conventions)
```

Content should cover: jobs/outcomes/risks/requirements for Intent the product;
components for the CLI subsystems (help, gen, validate, mutate, skill);
PDRs/ADRs/CRs for the major design choices (Go rewrite, embedded help, committed
docs + check gate).

### Cleanup

- Delete the entire `archive/` directory.
- Rewrite `README.md`: Go binary, mise install, `install-skill`,
  validate/build/check workflow.
- Final pass on `internal/help/content/errors/E00N.md` if dogfood exposes new
  edge cases.

## Tasks

Build in order — each step depends on the tree the prior one authored.

1. **`dogfood-scaffold`** — scaffold + product spine (jobs, outcomes, risks,
   requirements).
   - Verify: `intent validate` passes.
2. **`dogfood-engineering`** — architecture + components for the CLI subsystems.
   - Verify: `intent validate` passes; `intent trace` resolves across domains.
3. **`dogfood-records`** — PDRs/ADRs/CRs for the major decisions.
   - Verify: `intent validate` passes (E004 cross-domain rules hold).
4. **`dogfood-build-commit`** — `intent build`, commit docs, `intent check` green.
   - Verify: `intent check` exits zero; re-running `build` is byte-identical.
5. **`dogfood-finalize`** — delete `archive/`, rewrite README, stabilize the
   error catalog, and repoint CI's `check` from the seed fixture to the root
   `intent.yaml`.
   - Verify: `go test ./...`; CI green on the phase PR.

## Testing

| Scope | Verification |
|-------|----------------|
| Tree | `intent validate`, `intent check` |
| Repo | `go test ./...` |
| Docs | README accuracy read-through before merge |

## Rollout

1. Build tasks 1–5 in order in one session; open one PR for the phase.
2. Human review + CI → merge.

No feature flags. The phase is independently revertable via git revert of its PR.

## Open Questions

1. **CI switch** — task 5 must explicitly move CI's `check` from the seed
   fixture to the root `intent.yaml`; it is easy to miss.
2. **Tree scope** — decide how deep the dogfood goes. Enough to exercise every
   element type and both domains, without turning into a full product spec.
