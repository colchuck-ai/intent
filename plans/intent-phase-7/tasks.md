---
plan_slug: intent-phase-7
phase: tasks
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
requirements_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-7/requirements.md
implementation_plan_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-7/implementation-plan.md
status: created
created_at: '2026-07-23T21:11:00Z'
updated_at: '2026-07-24T16:59:00Z'
created_beads_at: '2026-07-23T21:17:14Z'
---

# Task Plan: Phase 7 — Skill renderer & install-skill

Calibration convoy for Gas City `build-from-convoy`. One PR; drain policy
`same-session`.

## Work beads

1. **skill-adapter-contract** — shared adapter interface + SKILL.md render helpers
2. **skill-claude-adapter** — Claude Code file tree (`.claude/skills/intent/SKILL.md`)
3. **skill-agentsmd-adapter** — generic `AGENTS.md` fallback adapter
4. **skill-install-cmd** — `intent install-skill --agent …` cobra command
5. **skill-tests** — render + CLI integration tests

## Bead Creation Payload

```yaml
target_rig: intent
labels:
  - plan:intent-phase-7
  - phase:7
convoys:
  - key: phase-7
    title: "Phase 7: Skill renderer & install-skill"
    description: |
      Generate agent skill files from embedded help content. Implements
      BUILD_PLAN Phase 7 and DESIGN §12 install-skill path.
    metadata:
      gc.plan.phase: "7"
      gc.plan.slug: intent-phase-7
    beads:
      - key: skill-adapter-contract
        title: Skill adapter contract and shared render helpers
        type: feature
        priority: 2
        description: |
          Add `internal/skill/` with an adapter interface that emits a per-agent
          file tree from embedded help. Implement shared SKILL.md assembly:
          trigger description, Intent paragraph, entry-point instruction, 2–3
          gotchas, judgment one-liners from help.InPlane("judgment") summaries.
          Do not duplicate full topic bodies — point to intent help slug.
        acceptance_criteria:
          - internal/skill/adapter.go and render.go exist with documented contract.
          - Render logic reads only from internal/help package APIs.
          - Unit tests cover judgment one-liner extraction and stable output order.
        files:
          - internal/skill/adapter.go
          - internal/skill/render.go
          - internal/skill/render_test.go
        verification:
          - go test ./internal/skill/...

      - key: skill-claude-adapter
        title: Claude Code skill adapter
        type: feature
        priority: 2
        description: |
          Implement Claude Code adapter writing .claude/skills/intent/SKILL.md
          under a configurable root (default .). Match .gitignore expectation
          that installed skills are never committed.
        acceptance_criteria:
          - internal/skill/claude.go renders SKILL.md via shared render helpers.
          - Output path is .claude/skills/intent/SKILL.md relative to install root.
        files:
          - internal/skill/claude.go
        dependencies:
          - skill-adapter-contract
        verification:
          - go test ./internal/skill/...

      - key: skill-agentsmd-adapter
        title: AGENTS.md fallback adapter
        type: feature
        priority: 2
        description: |
          Implement generic adapter that writes or updates AGENTS.md with the
          same core skill content for non-Claude targets.
        acceptance_criteria:
          - internal/skill/agentsmd.go renders AGENTS.md section via shared helpers.
          - Adapter name is agents-md for the --agent flag.
        files:
          - internal/skill/agentsmd.go
        dependencies:
          - skill-adapter-contract
        verification:
          - go test ./internal/skill/...

      - key: skill-install-cmd
        title: install-skill CLI command
        type: feature
        priority: 2
        description: |
          Add intent install-skill --agent target [--dir .] command. Register
          in internal/cli/root.go. Support agents claude-code and agents-md.
        acceptance_criteria:
          - Command writes files to disk and prints paths installed.
          - Unknown agent returns actionable error listing supported agents.
        files:
          - internal/cli/install_skill.go
          - internal/cli/root.go
        dependencies:
          - skill-claude-adapter
          - skill-agentsmd-adapter
        verification:
          - go test ./internal/cli/...

      - key: skill-tests
        title: Skill render integration tests
        type: chore
        priority: 2
        description: |
          Add integration tests proving rendered skill content is derived from
          embedded help (no hand-maintained duplicate prose). Test CLI end-to-end
          in temp directory.
        acceptance_criteria:
          - Tests assert every judgment slug summary appears in rendered SKILL.md.
          - go test ./... passes.
        files:
          - internal/skill/render_test.go
          - internal/cli/install_skill_test.go
        dependencies:
          - skill-install-cmd
        verification:
          - go test ./...
```

## Created Beads

| Key | Kind | Bead ID | Title |
|---|---|---|---|
| phase-7 | convoy | int-zui | Phase 7: Skill renderer & install-skill |
| skill-adapter-contract | bead | int-sv6 | Skill adapter contract and shared render helpers |
| skill-claude-adapter | bead | int-5a3 | Claude Code skill adapter |
| skill-agentsmd-adapter | bead | int-e3d | AGENTS.md fallback adapter |
| skill-install-cmd | bead | int-wmv | install-skill CLI command |
| skill-tests | bead | int-k7o | Skill render integration tests |
