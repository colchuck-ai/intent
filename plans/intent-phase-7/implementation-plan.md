---
plan_slug: intent-phase-7
phase: implementation-plan
requirements_file: plans/intent-phase-7/requirements.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-27T00:00:00Z
---

# Implementation Plan: Intent Phase 7 — Skill renderer & install-skill

## Summary

Add `intent install-skill --agent <target>` plus skill adapters that render
skill files from `internal/help`. One working session, one PR. Phase 6
(`internal/help/`, `intent help`, footers) is merged and is the sole source for
skill rendering.

**Status: complete.** All five tasks are built, committed, and verified. See
`plans/intent-phase-7/tasks.md` for the per-task checklist and the commits
that landed.

## Current System

| Area | Location | Notes |
|------|----------|-------|
| CLI root | `internal/cli/root.go` | Cobra tree; commands registered in `AddCommand` |
| Embedded help | `internal/help/` | `go:embed content/`; `All()`, `Resolve()`, `InPlane()` |
| Judgment topics | `internal/help/content/judgment/` | 10 topics; source of the skill one-liners |
| Skill package | `internal/skill/` | Adapter contract, shared render, Claude adapter (built) |
| Footers | `internal/cli/footer.go` | Judgment slugs resolve to embedded topics |
| Gitignore | `.gitignore` | Already ignores `.claude/skills/intent/` |

Module: `github.com/colchuck-ai/intent`, Go 1.25, cobra + yaml.v3.

## Shipped contract — build against this, not the original sketch

`internal/skill/adapter.go` and `render.go` are merged. The interface differs
from this plan's first draft (`Content` is passed in explicitly, and the file
type is `File`, not `OutputFile`). The remaining tasks must implement against
the real shape:

```go
// internal/skill/adapter.go
type File struct {
    Path    string // relative to install root
    Content []byte
}

type Adapter interface {
    Name() string                                 // --agent flag value
    Files(root string, c Content) ([]File, error) // must NOT write to disk
}

// internal/skill/render.go
type Content struct {
    Description string            // trigger; frontmatter only
    Intent      string            // one paragraph
    EntryPoint  string
    Gotchas     []string
    Judgments   []JudgmentPointer // {Slug, Summary}
}

func Render() Content               // assembles Content from embedded help
func (c Content) Markdown() string  // shared body; EXCLUDES Description
```

Key invariants:

- `Render()` reads only `help.InPlane("judgment")` and fixed template strings.
  No topic bodies are duplicated in code.
- `help.InPlane` returns slug-sorted topics, so output order is deterministic.
- Adapters return files; the caller writes them. An adapter that merges into an
  existing file may *read* from `root`, but must return the merged result.

`internal/skill/claude.go` is the reference implementation to mirror.

## Remaining work

### Task 3 — `AGENTS.md` fallback adapter

`internal/skill/agentsmd.go`. `AGENTSAdapter` implementing `Adapter`, with
`Name() == "agents-md"`. Emits `AGENTS.md` using `c.Markdown()` — the same
shared body, without Claude's YAML frontmatter. If `AGENTS.md` already exists at
`root`, read it and return the merged result with the Intent section replaced
between stable markers; never write from inside the adapter.

### Task 4 — `install-skill` CLI command

`internal/cli/install_skill.go`, registered in `root.go`'s `AddCommand(...)`.

```
intent install-skill --agent claude-code [--dir .]
intent install-skill --agent agents-md   [--dir .]
```

Selects the adapter by `Name()`, calls `skill.Render()`, writes each returned
`File` under `--dir` (creating parent directories), and prints the installed
paths. An unknown `--agent` returns an actionable error listing supported
agents.

### Task 5 — Render and CLI integration tests

`internal/skill/render_test.go` (extend) and `internal/cli/install_skill_test.go`.

- Assert every judgment slug summary appears in the rendered `SKILL.md`.
- Assert rendered content is help-derived — no hand-maintained duplicate prose.
- Run the CLI end-to-end in a temp dir; verify file placement matches the
  gitignored path.

## Testing

| Scope | Verification |
|-------|----------------|
| Package | `go test ./internal/skill/... ./internal/cli/...` |
| Repo | `go test ./...` |
| Manual | `intent install-skill --agent claude-code` in a scratch dir |

## Rollout

Build tasks 3 → 4 → 5 in order (each depends on the prior), then open one PR for
the phase. Human review + CI → merge → proceed to Phase 8.

No feature flags. The phase is independently revertable via git revert of its PR.

## Open Questions

1. **`description` wording** — confirm whether `render.go`'s shortened
   `description` constant is intentional relative to the retired
   `archive/intent/SKILL.md` wording, before the phase merges. Non-blocking.
2. **`AGENTS.md` merge markers** — decide the exact marker syntax for replacing
   an existing Intent section on reinstall.
