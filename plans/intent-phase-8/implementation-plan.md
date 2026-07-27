---
plan_slug: intent-phase-8
phase: implementation-plan
requirements_file: plans/intent-phase-8/requirements.md
status: approved
created_at: 2026-07-23T21:06:00Z
updated_at: 2026-07-27T00:00:00Z
---

# Implementation Plan: Intent Phase 8 — Distribution & CI

## Summary

Pin and ship the binary; CI enforces `validate` + `check`. One working session,
one PR. Depends on Phase 7 merged.

**Status: complete.** All four tasks are built, tested locally, and reviewed
(Standards + Spec). Awaiting a pushed PR for a green CI run before merge.

## Current System

| Area | Location | Notes |
|------|----------|-------|
| Generator | `internal/gen/` | `build` / `check` drift gate (Phase 3) |
| Version stamp | `internal/cli/root.go` | `var version = "0.0.0-dev"`, set via `-ldflags` |
| Seed fixture | `internal/model/testdata/seed.intent.yaml` | Test-only; not repo dogfood |
| CI / release | — | None yet (no `.github/`, no `goreleaser`, no `mise.toml`) |

Module: `github.com/colchuck-ai/intent`, Go 1.25, cobra + yaml.v3.

## Proposed Implementation

**Goal:** Pin and ship the binary; CI enforces validate + check.

**New files:**

```
.goreleaser.yaml
mise.toml
.github/workflows/ci.yaml
.githooks/pre-commit          # or scripts/ — pick one and document it
```

### Goreleaser

- Build `cmd/intent` for darwin/linux/windows, amd64 + arm64 where applicable.
- Embed version via
  `-ldflags "-X github.com/colchuck-ai/intent/internal/cli.version=..."`.
- Checksums + GitHub release, tag-triggered.

### mise.toml

```toml
[tools]
"go" = "1.25"
# After first release:
# "ubi:colchuck-ai/intent" = "<version>"
```

Until a release exists, CI builds from source. Document the pin upgrade path in
the PR description.

### CI workflow

```yaml
# Steps:
# - checkout
# - setup go 1.25 (or mise)
# - go test ./...
# - go build -o bin/intent ./cmd/intent
# - bin/intent validate -f <fixture>
# - bin/intent check   (once dogfood exists — Phase 9)
```

### Pre-commit

Run `intent check` when `intent.yaml` or `docs/` change. Leave `go test ./...`
to CI.

**Phase 8 / Phase 9 sequencing:** CI `check` against the real dogfood tree lands
in Phase 9, when `intent.yaml` exists. Phase 8 CI runs `go test` + `validate`
against `internal/model/testdata/seed.intent.yaml`; Phase 9 repoints it at the
root tree.

## Tasks

Config files are disjoint, so tasks 1–4 can be built in any order; verify
together at the end.

1. **`dist-goreleaser`** — `.goreleaser.yaml` + a local `--snapshot` dry run.
   - ✅ done — `goreleaser build --snapshot --clean` succeeds on all six
     darwin/linux/windows × amd64/arm64 targets; version ldflag verified.
2. **`dist-mise`** — `mise.toml` + README snippet documenting the pin.
   - ✅ done — `mise install` resolves go 1.25.12; README's Development
     section points at the post-release `ubi` pin in `mise.toml`.
3. **`dist-ci`** — `.github/workflows/ci.yaml` running test + build + validate.
   - ✅ done — valid YAML; the exact steps (`go test ./...`, `go build`,
     `intent validate -f internal/model/testdata/seed.intent.yaml`) pass
     locally. Green-on-PR still pending an actual push.
4. **`dist-precommit`** — local drift hook + a line in the README on enabling it.
   - ✅ done — `.githooks/pre-commit`; verified it no-ops on unrelated
     changes, warns-and-allows when `intent` isn't installed, and blocks
     the commit when `intent check` fails on a staged `docs/` change.

## Testing

| Scope | Verification |
|-------|----------------|
| Repo | `go test ./...` |
| Release | `goreleaser build --snapshot --clean` |
| CI | Green on the phase PR |

## Rollout

1. Build tasks 1–4 in one session; open one PR for the phase.
2. Human review + CI → merge.
3. Tag the first release after merge (optional before Phase 9 unless pinning is
   needed).

No feature flags. The phase is independently revertable via git revert of its PR.

## Open Questions

1. **Release slug for the `ubi` pin** — confirm `colchuck-ai/intent` (assumed
   from the module path).
2. **Pre-commit vs CI-only** — the hook is optional if contributors use mise; CI
   is the backstop.
