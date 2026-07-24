---
plan_slug: intent-phase-7
phase: implementation-plan
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
requirements_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-7/requirements.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
---

# Implementation Plan: Intent Phase 7 — Skill renderer & install-skill

## Summary

Add `intent install-skill --agent <target>` plus skill adapters that render
skill files from `internal/help`. One Gas City convoy (`build-from-convoy`), one
PR. Phase 6 (`internal/help/`, `intent help`, footers) is merged and is the sole
source for skill rendering.

## Current System

| Area | Location | Notes |
|------|----------|-------|
| CLI root | `internal/cli/root.go` | Cobra tree; `help` via `SetHelpCommand` |
| Embedded help | `internal/help/` | `go:embed content/`; 40 topics; `All()`, `Resolve()`, `InPlane()` |
| Footers | `internal/cli/footer.go` | Judgment slugs resolve to embedded topics |
| Gitignore | `.gitignore` | Already ignores `.claude/skills/intent/` |

Module: `github.com/colchuck-ai/intent`, Go 1.25, cobra + yaml.v3.

## Proposed Implementation

### Convoy boundary: Phase 7 — Skill renderer

**Goal:** `intent install-skill --agent <target>` renders skill files from embedded help.

**New packages / files:**

```
internal/skill/
  adapter.go          # Adapter interface: Render(outDir string) ([]string, error)
  render.go           # Shared SKILL.md assembly from help.All() / judgment summaries
  claude.go           # Claude Code adapter → .claude/skills/intent/SKILL.md
  agentsmd.go         # AGENTS.md fallback adapter
  render_test.go      # Byte-identity / content-source tests
internal/cli/
  install_skill.go    # install-skill command, --agent flag, --dir override
  install_skill_test.go
```

**Adapter contract:**

```go
type Adapter interface {
    Name() string
    Files(root string) ([]OutputFile, error) // relative path + content
}
```

- `Render` reads only from `help.All()`, `help.InPlane("judgment")`, and fixed
  template strings for gotchas / entry-point (no duplicated topic bodies in code).
- Claude adapter: write `.claude/skills/intent/SKILL.md` under `--dir` (default `.`).
- `agents-md` adapter: append or write `AGENTS.md` section with same core content.

**SKILL.md structure** (DESIGN §12):

1. YAML frontmatter with trigger `description` (when to load the skill).
2. One paragraph: what Intent is.
3. Entry-point: drive CLI; `intent help`; `intent.yaml` canonical; never hand-edit generated docs.
4. 2–3 gotchas (suffix addressing, validate-before-write, check drift gate).
5. Judgment one-liners: one line per `judgment:*` topic summary from embedded frontmatter.

Deep topic bodies are **not** copied into SKILL.md — point to `intent help <slug>`.

**CLI:**

```
intent install-skill --agent claude-code [--dir .]
intent install-skill --agent agents-md [--dir .]
```

Register in `root.go` alongside other commands.

**Tests:**

- Golden or snapshot test: rendered SKILL.md contains every judgment summary slug.
- Test that rendering does not read filesystem help files (only `internal/help` package).
- Integration test: run CLI, verify files land under temp dir, gitignore paths match.

**Drain policy for GC:** `same-session` (sequential build-on-prior).

**Convoy beads (see `plans/intent-phase-7/tasks.md`):**

1. `skill-adapter-contract` — interface + shared render helpers
2. `skill-claude-adapter` — Claude Code file tree
3. `skill-agentsmd-adapter` — AGENTS.md fallback
4. `skill-install-cmd` — cobra command + wiring
5. `skill-tests` — render tests + CLI integration tests

### Gas City execution

Artifact root: `plans/intent-phase-7/` (canonical per-plan-slug layout; build
outputs under `plans/intent-phase-7/build/`).

1. Mayor writes `plans/intent-phase-7/tasks.md` + bead payload.
2. Dry-run + create beads via `create_beads_from_tasks.py`.
3. Sling:

```bash
gc sling gc.run-operator <phase-7-convoy-id> --on build-from-convoy \
  --var artifact_root=plans/intent-phase-7/build \
  --var requirements_path=plans/intent-phase-7/requirements.md \
  --var plan_path=plans/intent-phase-7/implementation-plan.md \
  --var plan_review_path=plans/intent-phase-7/plan-review.md \
  --var decomposition_path=plans/intent-phase-7/tasks.md \
  --var interaction_mode=interactive \
  --var review_mode=agent \
  --var drain_policy=same-session \
  --var max_iterations=4 \
  --var open_pr=true
```

4. Human review PR + CI → merge → proceed to Phase 8.

## Testing

| Phase | Verification |
|-------|----------------|
| 7 | `go test ./internal/skill/... ./internal/cli/...`; manual `intent install-skill --agent claude-code` |

## Rollout

1. Approve requirements + implementation plan + plan review (done).
2. Decompose Phase 7 → beads → sling → merge (calibration slice).

No feature flags. The phase is independently revertable via git revert of its PR.

## Open Questions

1. **Rig-scoped gascity roles** — verify worker routing before first sling (`gc rig status intent`).
2. **install-skill test flakiness** — if path layout is flaky, pin golden files under `internal/skill/testdata/`.
