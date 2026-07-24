---
plan_slug: intent-phase-8
phase: implementation-plan
rig: intent
rig_root: /Users/max.dunn/dev/personal/colchuck-ai/intent
artifact_root: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans
requirements_file: /Users/max.dunn/dev/personal/colchuck-ai/intent/plans/intent-phase-8/requirements.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-24T16:59:00Z
---

# Implementation Plan: Intent Phase 8 — Distribution & CI

## Summary

Pin and ship the binary; CI enforces `validate` + `check`. One Gas City convoy
(`build-from-convoy`), one PR. Depends on Phase 7 merged.

## Current System

| Area | Location | Notes |
|------|----------|-------|
| Generator | `internal/gen/` | `build` / `check` drift gate (Phase 3) |
| Seed fixture | `internal/model/testdata/seed.intent.yaml` | Test-only; not repo dogfood |
| CI / release | — | None yet (no `.github/`, no `goreleaser`, no `mise.toml`) |

Module: `github.com/colchuck-ai/intent`, Go 1.25, cobra + yaml.v3.

## Proposed Implementation

### Convoy boundary: Phase 8 — Distribution & CI

**Goal:** Pin and ship the binary; CI enforces validate + check.

**New files:**

```
.goreleaser.yaml
mise.toml
.github/workflows/ci.yaml
scripts/pre-commit-intent-check   # or .githooks/pre-commit — repo convention TBD
```

**Goreleaser:**

- Build `cmd/intent` for darwin/linux/windows amd64 + arm64 where applicable.
- Embed version via `-ldflags "-X github.com/colchuck-ai/intent/internal/cli.version=..."`.
- Checksums + GitHub release (workflow_dispatch or tag-triggered — match org conventions).

**mise.toml:**

```toml
[tools]
"go" = "1.25"
# After first release:
# "ubi:colchuck-ai/intent" = "<version>"
```

Until a release exists, CI can `go build` + `go install` locally; document the
pin upgrade path in commit message / final report.

**CI workflow (`.github/workflows/ci.yaml`):**

```yaml
# Pseudocode steps:
# - checkout
# - setup go 1.25 (or mise)
# - go test ./...
# - go build -o bin/intent ./cmd/intent
# - bin/intent validate -f <fixture or future intent.yaml>
# - bin/intent check   (once dogfood exists; Phase 8 may use seed fixture until P9)
```

**Pre-commit:** run `intent check` when `intent.yaml` or `docs/` change; run
`go test ./...` optionally or defer to CI only.

**Phase 8 / Phase 9 sequencing note:** CI `check` against real dogfood lands in
Phase 9 when `intent.yaml` exists. Phase 8 CI should run `go test` + `validate`
against `internal/model/testdata/seed.intent.yaml` (or a minimal committed
fixture); Phase 9 updates CI to the root dogfood tree.

**Drain policy for GC:** `separate` (disjoint config files).

**Convoy beads (see `plans/intent-phase-8/tasks.md`):**

1. `dist-goreleaser` — release config + local dry-run docs
2. `dist-mise` — mise.toml + README snippet for pinning
3. `dist-ci` — GitHub Actions workflow
4. `dist-precommit` — local drift hook

### Gas City execution

Artifact root: `plans/intent-phase-8/` (canonical per-plan-slug layout; build
outputs under `plans/intent-phase-8/build/`).

1. Mayor writes `plans/intent-phase-8/tasks.md` + bead payload.
2. Dry-run + create beads via `create_beads_from_tasks.py`.
3. Sling:

```bash
gc sling gc.run-operator <phase-8-convoy-id> --on build-from-convoy \
  --var artifact_root=plans/intent-phase-8/build \
  --var requirements_path=plans/intent-phase-8/requirements.md \
  --var plan_path=plans/intent-phase-8/implementation-plan.md \
  --var plan_review_path=plans/intent-phase-8/plan-review.md \
  --var decomposition_path=plans/intent-phase-8/tasks.md \
  --var interaction_mode=interactive \
  --var review_mode=agent \
  --var drain_policy=separate \
  --var open_pr=true
```

4. Human review PR + CI → merge → proceed to Phase 9.

## Testing

| Phase | Verification |
|-------|----------------|
| 8 | `go test ./...`; goreleaser `--snapshot`; CI green on PR |

## Rollout

1. Approve requirements + implementation plan + plan review.
2. Decompose Phase 8 → beads → sling → merge.
3. Tag first release after Phase 8 merge (optional before Phase 9 if pinning needed).

No feature flags. The phase is independently revertable via git revert of its PR.

## Open Questions

1. **GitHub org/repo for goreleaser `ubi` pin** — confirm release slug (`colchuck-ai/intent` assumed from module path).
2. **Phase 8 CI before dogfood** — use seed fixture for `check` until Phase 9 lands root tree (recommended above).
3. **Pre-commit vs CI-only** — hook optional if contributors use mise; CI is the backstop.
