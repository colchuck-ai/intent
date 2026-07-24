---
plan_slug: intent-phase-9
phase: plan-review
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
requirements_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-9/requirements.md
implementation_plan_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-9/implementation-plan.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
review_verdict: approved
---

# Plan Review: Intent Phase 9 — Dogfood tree & finalize

## Scope reviewed

- Requirements: `plans/intent-phase-9/requirements.md`
- Implementation plan: `plans/intent-phase-9/implementation-plan.md`
- Ground truth: `BUILD_PLAN.md`, `DESIGN.md` §11–§12

## Findings

### Approved aspects

1. **Authoring constraint** — write-CLI-only dogfood tree is enforceable and
   testable; matches the BUILD_PLAN salvage rule for bucket 1.

2. **Convoy sizing matches BUILD_PLAN** — one phase, one PR, human merge gate.

3. **Drain policy** — `same-session` for the coupled code/content chain is
   correct, with the final cleanup bead depending on the docs bead.

### Risks (accepted with mitigations)

| Risk | Mitigation |
|------|------------|
| Dogfood prose quality | Interactive mode; human reads PR diff |
| CI still points at seed fixture | Bead 5 explicitly updates CI to root `intent.yaml` |
| GC workers may not route without rig roles | Verify `gc rig status intent` before the sling |

### Minor adjustments recommended during execution

- Phase 9 bead 5 must call out the CI switch from seed fixture to root
  `intent.yaml` (already in the implementation plan).

## Verdict

**Approved.** No blocking issues. Plan is concrete enough for Phase 9 bead
decomposition once Phases 7 and 8 have merged.

## Next checkpoint

After Phase 8 merges:

1. Mayor drafts `plans/intent-phase-9/tasks.md`
2. Create beads (dry-run → real)
3. Sling `build-from-convoy` with `drain_policy=same-session`
