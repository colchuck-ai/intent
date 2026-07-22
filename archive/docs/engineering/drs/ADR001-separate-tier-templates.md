# ADR001 - Separate Tier Templates

Requirement, outcome, and component each ship two template files — one per tier (Tier-1 own-document, Tier-2 own-directory) — instead of a single merged template with conditional sections.

## Context

The Tier-1 and Tier-2 templates for the same element type share most of their body (same headings, same minimal pattern, same record backlink sections); only link depth and one added subsection differ. Keeping them as separate files risks the two copies drifting out of sync as one is edited and the other isn't.

## Options

- Merge each pair into a single template file with conditional/optional blocks marked for Tier-2-only content: removes the duplication, but adds a second markup system (conditional blocks) on top of the tier system Intent already teaches.
- Keep both files, and add a header comment on the Tier-2 file stating exactly how it differs from its Tier-1 sibling, with a note to keep the shared body in sync: duplication remains, but each file stays a flat, copy-paste-ready skeleton.

## Decision

Kept both files with the delta called out in a header comment on the Tier-2 variant. A merged template would need its own notation for "this block is Tier-2 only," which is one more thing for an assistant to interpret correctly on top of the tier system itself; two flat files stay simpler to load and fill in directly.

## Consequences

- The two files can still drift — the header comment is a documentation aid, not an enforced constraint; nothing currently checks that the shared body of a Tier-1/Tier-2 pair stays in sync.
- Authors get a complete, ready-to-fill file per tier with no conditional markup to interpret.
