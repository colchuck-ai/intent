---
schema: gc.build.review.v1
workflow:
  id: int-9kb
  formula: build-from-convoy
methodology:
  pack: gascity
  name: build-from-convoy
producer:
  formula: build-from-review-base
  stage: review
  attempt: 1
status: blocked
trace:
  upstream:
    - path: plans/intent-phase-7/requirements.md
      hash: sha256:b98d5cd26e697400974f5b6cd3d9eea23281f8d04d08cc69d6019aefffcc098c
    - path: plans/intent-phase-7/implementation-plan.md
      hash: sha256:06d0e8b03a4942a8e18bd0afaf69da0ef9822e7a6c35a5019d95d44f8a81c984
    - path: plans/intent-phase-7/plan-review.md
      hash: sha256:2710560737203d9fcc8c4940f4a6482a5ee24eb78bd766a9bb40e97290eeb736
    - path: plans/intent-phase-7/tasks.md
      hash: sha256:dc2744f9df3496d3b42a2e30f803379695bcbb5bdfd6ea316e31a4b21b030a21
      ids: [int-sv6, int-5a3, int-e3d, int-wmv, int-k7o]
    - path: plans/intent-phase-7/build/implementation-skill-adapter-contract.md
      hash: sha256:9e3d5862a9ee17dcd8ea0416990c702157bf214cf72f42e1dd9097e63f158255
    - path: internal/skill/adapter.go
      hash: sha256:72da9c289477ba0f26dff746e6f737002f5a7c80855bcd24ce9f7d5631094bb3
    - path: internal/skill/render.go
      hash: sha256:e8c7174cc7b5ad29d0eaec45c566e91207cc7a2a3354994539aec3f0941d5595
    - path: internal/skill/render_test.go
      hash: sha256:0a1183ab3337f2a01284c92e0e877828ec49ed471a282217fe5a321a61573388
    - path: beads/int-sv6
      hash: bead:int-sv6
    - path: git:f630a73
      hash: git:f630a73a7b00a3aa75957f774d62955af0682006
  coverage:
    - id: int-sv6
      status: approved
      rationale: internal/skill/adapter.go, render.go, render_test.go independently re-verified against requirements, plan, and DESIGN §12; build/vet/gofmt/tests all pass; committed f630a73.
    - id: int-5a3
      status: blocked
      rationale: Never attempted. Drain int-7ws's skip_remaining policy skipped this item after int-sv6's wrapper (int-nin) recorded gc.outcome=fail; internal/skill/claude.go does not exist. Dependency on int-sv6 is now satisfied (closed), so this item is ready to attempt once the blocking infra defect is resolved.
    - id: int-e3d
      status: blocked
      rationale: Never attempted. Same skip_remaining cascade as int-5a3; internal/skill/agentsmd.go does not exist. Dependency on int-sv6 is now satisfied (closed), so this item is ready to attempt once the blocking infra defect is resolved.
    - id: int-wmv
      status: blocked
      rationale: Never attempted. Same skip_remaining cascade; internal/cli/install_skill.go does not exist. Also still dependency-blocked on int-5a3 and int-e3d per the decomposition.
    - id: int-k7o
      status: blocked
      rationale: Never attempted. Same skip_remaining cascade; no render or CLI integration tests were added. Also still dependency-blocked on int-wmv per the decomposition.
---

# Review: Intent Phase 7 — Skill renderer & install-skill (build-from-review-base)

## Verdict

**Blocked.** 1 of 5 convoy items (`int-sv6`, skill-adapter-contract) is
implemented, independently re-verified, and **approved** — see Findings §1.
The other 4 (`int-5a3`, `int-e3d`, `int-wmv`, `int-k7o`) were never attempted:
the same-session drain (`int-7ws`) stopped via its `skip_remaining` policy
after item 0's wrapper (`int-nin`) recorded `gc.outcome=fail` — but that
failure is a confirmed infrastructure/control-dispatch defect, not a rejection
of `int-sv6`'s code (Findings §2). Phase 7's user story — `intent
install-skill --agent claude-code` renders a thin skill file — is not
deliverable yet: `install-skill` itself, the Claude adapter, the AGENTS.md
adapter, and the integration tests proving content is help-derived all remain
unbuilt. The phase cannot merge in its current state, and the correct next
action is orchestration-layer, not a code review-fix loop (Findings §3).

| ID | Status |
| --- | --- |
| int-sv6 | approved |
| int-5a3 | blocked |
| int-e3d | blocked |
| int-wmv | blocked |
| int-k7o | blocked |

## Findings

1. **`int-sv6` code review: approved, no findings.** Independently re-ran (all
   read-only): `go build ./...` clean; `go vet ./internal/skill/...` clean;
   `gofmt -l internal/skill/` clean; `go test ./internal/skill/... -v` 5/5
   pass (`TestRenderJudgmentsMatchHelp`, `TestRenderJudgmentsStableOrder`,
   `TestContentMarkdownStructure`, `TestContentMarkdownExcludesFullTopicBodies`,
   `TestContentMarkdownExcludesDescription`); full-repo `go test ./...` green
   across every package. Recomputed sha256 for all three changed files from
   the working tree — identical, byte-for-byte, to the hashes recorded in the
   implementation summary at commit time: zero drift between what was
   implemented/tested and what's committed now. Cross-checked
   `internal/help/help.go`: `InPlane`, `Resolve`, and
   `Topic{Slug,Plane,Title,Summary,Body}` match exactly what `render.go` and
   `render_test.go` call — no invented API surface. Cross-checked DESIGN.md
   §12 directly: the shipped `Content` shape (trigger `Description`, one
   `Intent` paragraph, the identical 4-clause entry-point instruction, 3
   gotchas — within the "2–3" band, one `JudgmentPointer` per `judgment:*`
   topic sourced live via `help.InPlane("judgment")`) matches verbatim, and
   `.gitignore` already excludes `.claude/skills/intent/` as
   `requirements.md` assumed. Every acceptance criterion for
   `skill-adapter-contract` in `tasks.md` is satisfied.

2. **The `int-nin`=fail vs. `int-sv6`-evidence=approved discrepancy is
   adjudicated: infrastructure failure, not a code rejection.** Traced
   `int-nin` → its dependency `int-8sa` ("Implement shared-drain item"):
   `int-8sa`'s metadata records `gc.controller_error`: "int-8sa: running
   check: int-8sa: resolving check path: resolving gate condition path: lstat
   .../.gc/scripts: no such file or directory",
   `gc.failure_reason=control_dispatch_error`,
   `gc.final_disposition=control_quarantined`. This is the check-runner
   failing to locate the artifact validator script itself — not the
   validator rejecting the artifact's contents. Confirmed directly: `find`
   across the `intent` rig, the `gas-city-test` rig, and the `knowhere` city
   root shows `.gc/scripts/checks/build-artifact-valid.sh` does not exist
   anywhere in this Gas City workspace (only an unrelated
   `knowhere/.gc/scripts/gc-beads-bd.sh` exists); `gc doctor` reports no
   related structural warning, so this looks like a genuine, still-open setup
   gap rather than a transient blip. This same missing path gates every
   `ralph`-checked build step in this workflow, including this review stage's
   own closing gate and the downstream `finalize` stage
   (`gc.build.final-report.v1`) — so a `control_quarantined` mark on this
   step or the next should be read the same way: infra gap, not a verdict on
   this artifact's content.

3. **Remaining 4/5 items are unattempted, not failing, and a code fix-loop
   does not apply.** `int-5a3`, `int-e3d`, `int-wmv`, `int-k7o` were skipped
   by `int-7ws`'s `skip_remaining` policy after item 0's spurious failure —
   confirmed by its own `gc.drain_manifest.v1`
   (`failure_reason: previous_item_failed` for all four rows). No files exist
   for any of them (`internal/skill/claude.go`, `internal/skill/agentsmd.go`,
   `internal/cli/install_skill.go` all absent). `int-sv6` closing means
   `int-5a3` and `int-e3d` are now dependency-satisfied and ready to run in
   parallel; `int-wmv` still needs both of those; `int-k7o` still needs
   `int-wmv`. `review_mode=agent`'s structured fix handoff is a code-fix
   mechanism dispatched via `review_fix_formula` (`fix-loop-base`); there is
   no code to fix for these four, so that dispatch does not apply (see Fix
   Handoff below for what does).

4. **Minor plan drift, non-blocking.** The plan sketched
   `Adapter.Files(root string) ([]OutputFile, error)`; it shipped as
   `Files(root string, c Content) ([]File, error)` — `Content` is passed
   explicitly and the type is named `File`, not `OutputFile`. The
   implementation summary documents the rationale (adapters stay pure,
   stateless functions of `(root, Content)`). No consumer exists yet to
   contradict this shape — `int-5a3`/`int-e3d`, the two concrete adapters,
   are unimplemented — so nothing is broken, but `int-5a3`/`int-e3d`/`int-wmv`
   must implement against this actual shipped signature, not the plan's
   original sketch.

## Fix Handoff

Not a `review_fix_formula` (`fix-loop-base`) dispatch — there is no code to
fix. Structured handoff for the caller / next continuation:

- **Blocking defect (orchestration, not code):**
  `.gc/scripts/checks/build-artifact-valid.sh` does not exist anywhere in
  this Gas City workspace. Any `ralph`-gated build step that reaches its
  artifact-validation check (implement-item, this review step, the finalize
  step) will hard-fail with `control_dispatch_error` / `control_quarantined`
  until this script exists and is resolvable from the rig's
  `.gc/scripts/checks/`.
- Once that's fixed, restart the implementation drain over the remaining
  convoy items in dependency order: `int-5a3` and `int-e3d` can run in
  parallel now; `int-wmv` next; `int-k7o` last.
- No rework needed on `int-sv6` / `internal/skill/adapter.go` / `render.go` /
  `render_test.go` — approved as committed in `f630a73`.

## Verification

- Read `int-c8l` (this step), `int-9kb` (workflow root), `int-ui6`
  (prepare-review, closed), `int-7ws` (drain control), `int-nin`
  (do-work-item wrapper), `int-54d` (its finalize), `int-8sa` (implement-item
  — the actual quarantine site), `int-sv6` and its four sibling convoy beads
  `int-5a3`/`int-e3d`/`int-wmv`/`int-k7o`, and `int-zui` (convoy).
- Confirmed `.gc/scripts/checks/build-artifact-valid.sh` absent via
  filesystem search across the `intent` rig, the `gas-city-test` rig, and the
  `knowhere` city root; cross-checked against `gc doctor`, which reports no
  related structural warning.
- Independently reran (read-only): `go build ./...`, `go vet
  ./internal/skill/...`, `gofmt -l internal/skill/`, `go test
  ./internal/skill/... -v`, `go test ./...` — all pass, matching the
  implementation summary's self-reported results.
- Recomputed sha256 for `internal/skill/{adapter,render,render_test}.go` from
  the current working tree; identical to the hashes recorded in
  `plans/intent-phase-7/build/implementation-skill-adapter-contract.md` at
  implementation time — no drift.
- Cross-read `internal/help/help.go` to confirm `InPlane`, `Resolve`, and
  `Topic{Slug,Plane,Title,Summary,Body}` match exactly what `render.go` /
  `render_test.go` call.
- Cross-read `DESIGN.md` §12 (SKILL.md shape) and `.gitignore`
  (`.claude/skills/intent/` already present) against the implementation and
  `requirements.md`'s acceptance criteria.
- `git log --oneline -8` / `git show f630a73 --stat` confirm a clean,
  isolated 3-file / 259-line commit matching the summary; `git status
  --short` shows no uncommitted Go source changes.
