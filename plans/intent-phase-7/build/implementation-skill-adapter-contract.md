---
schema: gc.build.implementation-summary.v1
workflow:
  id: int-nin
  formula: do-work-item
methodology:
  pack: gascity
  name: build-basic
producer:
  formula: do-work-item
  stage: implement-item
  attempt: 1
status: approved
trace:
  upstream:
    - path: beads/int-sv6
      hash: bead:int-sv6
      ids: [int-sv6]
    - path: internal/skill/adapter.go
      hash: sha256:72da9c289477ba0f26dff746e6f737002f5a7c80855bcd24ce9f7d5631094bb3
    - path: internal/skill/render.go
      hash: sha256:e8c7174cc7b5ad29d0eaec45c566e91207cc7a2a3354994539aec3f0941d5595
    - path: internal/skill/render_test.go
      hash: sha256:0a1183ab3337f2a01284c92e0e877828ec49ed471a282217fe5a321a61573388
  coverage:
    - id: int-sv6
      status: covered
---

# Implementation summary: skill-adapter-contract (int-sv6)

## Summary

Added `internal/skill/` with the per-agent `Adapter` interface (`adapter.go`)
and the shared SKILL.md render helpers (`render.go`): a `Content` struct
assembled by `Render()` and rendered to Markdown by `Content.Markdown()`,
sourcing judgment one-liners from `help.InPlane("judgment")` per DESIGN §12.

## Intended Behavior

- `Adapter` is the contract concrete per-agent adapters (Claude Code, the
  AGENTS.md fallback — separate drain items in this same convoy) implement:
  `Name() string` for the `--agent` flag value, and
  `Files(root string, c Content) ([]File, error)` to emit a file tree without
  touching disk themselves, so adapters stay pure and testable and the
  install-skill CLI (also a separate drain item) owns all disk I/O.
- `Render()` assembles `Content` from embedded help: a fixed trigger
  `Description`, one `Intent` paragraph, the `EntryPoint` instruction, 2-3
  `Gotchas`, and one `JudgmentPointer{Slug, Summary}` per `judgment:*` topic,
  in the slug order `help.InPlane` already guarantees.
- `Content.Markdown()` renders the part of the body every adapter shares
  (Intent paragraph, entry point, gotchas, judgment pointers as
  `<summary> — \`intent help <slug>\`` lines) so no adapter hand-duplicates a
  judgment topic's full question/test/good/bad/failure-mode body. `Description`
  is intentionally excluded from `Markdown()` — it's Claude-style frontmatter
  trigger text, and adapters without a frontmatter concept (e.g. AGENTS.md)
  need to decide independently whether/where to fold it in.
- `Description`, the `Intent` paragraph, `EntryPoint`, and `Gotchas` are
  authored constants in `render.go`, not derived from `help`, because DESIGN §12
  specifies their content/shape but there is no corresponding help plane to
  source them from (only the judgment one-liners are dynamic, per the
  acceptance criterion that render logic reads only `internal/help` APIs).
  `archive/intent/SKILL.md` (the retired hand-maintained skill) was read for
  tone/precedent on the trigger `description` only; the rest of that file's
  `references/`+`scripts/` structure is exactly what DESIGN §12 replaces and
  was not reused.

## Changed Files

- `internal/skill/adapter.go` (new) — `File`, `Adapter` interface, package doc.
- `internal/skill/render.go` (new) — `Content`, `JudgmentPointer`, `Render()`,
  `Content.Markdown()`.
- `internal/skill/render_test.go` (new) — coverage below.

## Verification

- First verification command, run right after writing the implementation:
  `go test ./internal/skill/... -v` → **PASS**, 5/5 tests
  (`TestRenderJudgmentsMatchHelp`, `TestRenderJudgmentsStableOrder`,
  `TestContentMarkdownStructure`, `TestContentMarkdownExcludesFullTopicBodies`,
  `TestContentMarkdownExcludesDescription`).
- Final proof command, run after committing: `go build ./... && go test ./...`
  → **PASS** — build clean; every package passes, including
  `github.com/colchuck-ai/intent/internal/skill  0.229s`. `gofmt -l
  internal/skill/` and `go vet ./internal/skill/...` are both clean (no
  output).
- `TestRenderJudgmentsMatchHelp` and `TestRenderJudgmentsStableOrder` are the
  acceptance-criteria tests for "judgment one-liner extraction and stable
  output order": they assert every `help.InPlane("judgment")` topic is
  represented exactly once with its real slug/summary, in slug order, and
  identically across repeated `Render()` calls.
  `TestContentMarkdownExcludesFullTopicBodies` asserts a good/bad-example
  sentence that only exists in a judgment topic's full body
  (`judgment:risk-vs-feature`) is absent from the rendered Markdown, enforcing
  "do not duplicate full topic bodies — point to intent help slug."

## Remaining Risks

- The generic step instructions reference "the authoritative worktree
  recorded on the source anchor" and `cd "$WORKTREE"`. No worktree/`work_dir`
  metadata exists anywhere in this drain's bead graph (checked the source
  anchor `int-sv6`, the synthetic drain-unit convoy `int-8t2`, the drain
  control `int-7ws`, the `do-work-item` root `int-nin`, the `build-from-convoy`
  root `int-9kb`, and the parent convoy `int-zui`), and this formula's graph
  has no `prepare-worktree` step (`int-nin`'s only real dependency is the
  `workflow-finalize` control bead `int-54d`). `git worktree list` shows a
  single worktree at the rig root, on `main`. Given
  `drain_policy=same-session`, I implemented directly in that existing
  checkout after verifying `pwd -P` equals the rig root, rather than failing
  the step for a worktree contract this formula never provisions. This is
  consistent with the sibling `plans/intent-phase-7/build/review.md` from the
  prior (failed, pre-implementation) attempt, whose own verification
  methodology checked the rig root directly for `internal/skill/`, `git log`,
  and `git status` rather than any separate worktree path.
- `Description`, `Intent`, `EntryPoint`, and `Gotchas` string content is a
  reasonable authored rendering of DESIGN §12's prose, not a literal quote of
  a spec value — DESIGN §12 prescribes the shape and the entry-point
  instruction's four clauses verbatim, but leaves the exact `Intent` paragraph
  and the 2-3 gotchas open. A human may want to tune wording; the structure
  (fields, order, "point at the slug" rule) is the load-bearing contract these
  tests pin, not the exact prose.
- `Adapter` is exercised only indirectly (no concrete implementation exists
  yet in this checkout) since the Claude Code adapter (`int-5a3`) and AGENTS.md
  adapter (`int-e3d`) are separate drain items in this same convoy; the
  interface shape is my best design against both known concrete uses
  described in `plans/intent-phase-7/tasks.md`, not yet proven by a second
  implementation.

## Coverage

| ID | Status |
| --- | --- |
| int-sv6 | covered |
