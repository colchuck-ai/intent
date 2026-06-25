---
name: intent
description: Build and maintain a traceable docs/product and docs/engineering documentation tree that links jobs-to-be-done and outcomes through risks, requirements, architecture, and components, with change records (CR), product decision records (PDR), and architectural decision records (ADR). Use when creating, editing, tracing, or placing these documents and their IDs/paths/tiers.
---

# Intent

A framework for structuring product and engineering intent as a traceable hierarchy — from jobs and outcomes through risks and requirements down to architecture and components — with decision and change records that keep the whole narrative coherent over time.

## Global Guidance

- Don't drift into implementation: jobs, outcomes, and risks must name no product, feature, or technology.
- Keep outcomes measurable and from the customer's perspective (what "better" looks like), not feature requests.
- State requirements as what the product must do, not how it's built; tie each to the risk(s) it mitigates.
- Keep the trace coherent: nothing downstream may contradict what an ancestor claims.
- Be candid about unknowns: "I don't know" / "we need more information" are valid, and mark assumptions as assumptions.
- Write plainly so a junior PM or engineer gets it on the first pass.

## Gotchas

Non-obvious traps in this framework. Read before diving into a reference file.

- Tier prerequisites cascade: a Tier-1 requirement file needs a Tier-2 *nested* outcome folder first, and Tier-1 outcome/requirement/component change records need a Tier-2 nested parent (product- and engineering-scoped CRs do not). Promote the parent before creating the deeper path. See [references/structure.md](references/structure.md).
- Creating a brand-new element produces no Change Record — there is no prior state. CRs are only for updates and deletes. See [references/workflows.md](references/workflows.md).
- A Change Record is not a decision record: a CR logs *what changed and why*; PDRs/ADRs capture *a decision and its alternatives*. Don't conflate them. See [references/elements.md](references/elements.md).
- Slug vs logical ID: filenames/folders (slugs) use only the local ID segment (e.g. `O001-<name>`); cross-tree references in document bodies use fully-qualified logical IDs (e.g. `O001-R001`). Don't repeat ancestor IDs in slugs. See [references/structure.md](references/structure.md).

## Elements

Load when you need to define, judge, or explain what something *is* — jobs, outcomes, risks, requirements, architecture, components, or records (CR, PDR, ADR). Use the minimal patterns and good/bad examples as the rubric for drafting and review.

See [references/elements.md](references/elements.md).

## Structure

Load when you need to place content on disk — choosing a complexity tier, naming a file or folder, resolving a logical ID, or finding the path and template for an element or record.

See [references/structure.md](references/structure.md).

## Workflows

Load when you are reading, creating, updating, or deleting intent documents. Covers context tracing (what to load before acting), the read workflow, and the three-session modification process (draft, judge, review).

See [references/workflows.md](references/workflows.md).
