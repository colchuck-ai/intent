---
plan_slug: intent-phase-8
phase: plan-review
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
requirements_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-8/requirements.md
implementation_plan_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-8/implementation-plan.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
review_verdict: approved
---

# Plan Review: Intent Phase 8 — Distribution & CI

## Scope reviewed

- Requirements: `plans/intent-phase-8/requirements.md`
- Implementation plan: `plans/intent-phase-8/implementation-plan.md`
- Ground truth: `BUILD_PLAN.md`, `DESIGN.md` §11

## Findings

### Approved aspects

1. **CI sequencing is realistic** — running `intent check` against the seed
   fixture until the Phase 9 dogfood exists avoids a chicken-and-egg failure.

2. **Convoy sizing matches BUILD_PLAN** — one phase, one PR, human merge gate.

3. **Drain policy** — `separate` for the disjoint config files is correct.

### Risks (accepted with mitigations)

| Risk | Mitigation |
|------|------------|
| goreleaser needs GitHub token / first tag | Bead includes snapshot dry-run; tag after merge |
| GC workers may not route without rig roles | Verify `gc rig status intent` before first sling |

### Minor adjustments recommended during execution

- Confirm the release slug (`colchuck-ai/intent`) before wiring the `ubi` pin.

## Verdict

**Approved.** No blocking issues. Plan is concrete enough for Phase 8 bead
decomposition once Phase 7 has merged.

## Next checkpoint

After Phase 7 merges:

1. Mayor drafts `plans/intent-phase-8/tasks.md`
2. Create beads (dry-run → real)
3. Sling `build-from-convoy` with `drain_policy=separate`
