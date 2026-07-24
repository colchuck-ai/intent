---
plan_slug: intent-phase-7
phase: plan-review
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
requirements_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-7/requirements.md
implementation_plan_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-7/implementation-plan.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
review_verdict: approved
---

# Plan Review: Intent Phase 7 — Skill renderer & install-skill

## Scope reviewed

- Requirements: `plans/intent-phase-7/requirements.md`
- Implementation plan: `plans/intent-phase-7/implementation-plan.md`
- Ground truth: `BUILD_PLAN.md`, `DESIGN.md` §12, merged Phase 6 (`internal/help/`)

## Findings

### Approved aspects

1. **Convoy sizing matches BUILD_PLAN** — one phase, one PR, human merge gate.
   Aligns with Gas City worktree isolation (one run cannot see unmerged prior work).

2. **Grounded in existing API** — `internal/help.All()`, `InPlane()`, topic
   frontmatter summaries are sufficient for thin SKILL.md assembly without
   duplicating bodies.

3. **Drain policy** — `same-session` for the coupled code chain is correct.

### Risks (accepted with mitigations)

| Risk | Mitigation |
|------|------------|
| Agent review is not the backstop | Human merge gate + CI on every PR |
| GC workers may not route without rig roles | Verify `gc rig status intent` before first sling |

### Minor adjustments recommended during execution

- If `install-skill` tests are flaky on path layout, pin golden files under
  `internal/skill/testdata/`.

## Verdict

**Approved.** No blocking issues. Plan is concrete enough for Phase 7 bead
decomposition and the calibration sling.

## Next checkpoint

After approval:

1. Mayor drafts `plans/intent-phase-7/tasks.md`
2. Create beads (dry-run → real)
3. Sling `build-from-convoy` with `drain_policy=same-session`, `max_iterations=4`
