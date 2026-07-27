---
plan_slug: intent-phase-7
phase: tasks
requirements_file: plans/intent-phase-7/requirements.md
implementation_plan_file: plans/intent-phase-7/implementation-plan.md
status: done
created_at: 2026-07-23T21:11:00Z
updated_at: 2026-07-27T00:00:00Z
---

# Task Plan: Phase 7 — Skill renderer & install-skill

Five tasks, built in dependency order in one working session. All five are
done.

## Status

| # | Task | Status | Evidence |
|---|------|--------|----------|
| 1 | Adapter contract + shared render helpers | ✅ done | `f630a73` — `adapter.go`, `render.go`, `render_test.go` |
| 2 | Claude Code adapter | ✅ done | `06760b6` — `claude.go`, `claude_test.go` |
| 3 | `AGENTS.md` fallback adapter | ✅ done | `568de62` — `agentsmd.go`, `agentsmd_test.go` |
| 4 | `install-skill` CLI command | ✅ done | `568de62` — `install_skill.go`, `root.go` |
| 5 | Render + CLI integration tests | ✅ done | `568de62` — `render_test.go`, `install_skill_test.go` |

Tasks 1 and 2 were independently reviewed with no functional findings; `go
build`, `go vet`, `gofmt`, and `go test ./...` are all clean at `06760b6`.

Tasks 3–5 were built directly (the earlier external build-orchestrator attempt
referenced below was abandoned and is no longer relevant). They were also
independently reviewed on both a standards axis and a spec-fidelity axis; the
one finding (a non-idiomatic, unchecked `cmd.MarkFlagRequired` call in
`install_skill.go`) was fixed before commit. `go build`, `go vet`, `gofmt`, and
`go test ./...` are all clean at `568de62`.

## Tasks

### 1. Skill adapter contract and shared render helpers ✅

Add `internal/skill/` with an adapter interface that emits a per-agent file tree
from embedded help. Shared `SKILL.md` assembly: trigger description, Intent
paragraph, entry-point instruction, 2–3 gotchas, judgment one-liners from
`help.InPlane("judgment")` summaries. No full topic bodies — point to
`intent help <slug>`.

- Files: `internal/skill/adapter.go`, `render.go`, `render_test.go`
- Verify: `go test ./internal/skill/...`

### 2. Claude Code skill adapter ✅

Claude Code adapter writing `.claude/skills/intent/SKILL.md` under a
configurable root (default `.`), matching the `.gitignore` expectation that
installed skills are never committed.

- Files: `internal/skill/claude.go`, `claude_test.go`
- Verify: `go test ./internal/skill/...`

### 3. AGENTS.md fallback adapter ✅

Generic adapter that writes or updates `AGENTS.md` with the same core skill
content for non-Claude targets. Adapter name is `agents-md`.

- Acceptance:
  - `internal/skill/agentsmd.go` renders the `AGENTS.md` section via the shared
    helpers (`c.Markdown()`), with no Claude-style frontmatter.
  - `Name()` returns `agents-md`.
  - An existing `AGENTS.md` is merged, not clobbered; the adapter returns the
    merged bytes rather than writing them.
- Files: `internal/skill/agentsmd.go`
- Depends on: task 1
- Verify: `go test ./internal/skill/...`

### 4. install-skill CLI command ✅

Add `intent install-skill --agent <target> [--dir .]`. Register in
`internal/cli/root.go`. Support `claude-code` and `agents-md`.

- Acceptance:
  - Command writes files to disk (creating parent dirs) and prints the paths
    installed.
  - Unknown agent returns an actionable error listing supported agents.
- Files: `internal/cli/install_skill.go`, `internal/cli/root.go`
- Depends on: tasks 2, 3
- Verify: `go test ./internal/cli/...`

### 5. Skill render integration tests ✅

Integration tests proving rendered skill content is derived from embedded help
(no hand-maintained duplicate prose), plus an end-to-end CLI run in a temp dir.

- Acceptance:
  - Tests assert every judgment slug summary appears in the rendered `SKILL.md`.
  - `go test ./...` passes.
- Files: `internal/skill/render_test.go`, `internal/cli/install_skill_test.go`
- Depends on: task 4
- Verify: `go test ./...`

## Done when

`intent install-skill --agent claude-code` renders; the skill content is the
same bytes as the embedded help content; installed files are gitignored and
never committed; `go test ./...` is green.

All met as of `568de62`. Note: `AGENTS.md` (the `agents-md` target) is not
itself gitignored, since task 3 requires merging into a file a project may
already commit — the gitignore guarantee applies to the Claude Code target's
`.claude/skills/intent/` path, which is what this line was written against.
