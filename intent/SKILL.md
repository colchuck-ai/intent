---
name: intent
description: Use when creating, editing, tracing, placing, or reviewing product or engineering intent documents under `docs/product/` or `docs/engineering/` — jobs-to-be-done, outcomes, risks, requirements, architecture, components, and their Change Records (CR), Product Decision Records (PDR), and Architectural Decision Records (ADR). Fire even when the user doesn't name the framework: "write a product spec", "capture why we chose X", "log this change", "record this decision", "draft an ADR/PDR/CR", "define the customer's job", "trace a requirement to a component", "check these docs for coherence", "where does this doc go / what should I name it". Also use for the draft/judge/review workflow and alignment checks. Skip for code changes, READMEs, marketing copy, or ad-hoc notes that don't link into the intent tree.
---

# Intent

## Principles

- Be candid about unknowns: "I don't know" / "we need more information" are valid, and mark assumptions as assumptions.
- Write plainly so a junior PM or engineer gets it on the first pass.
- Product-side vocabulary follows Ulwick/ODI (jobs, outcomes, needs-first sequencing); job syntax follows the Christensen job-story form. Deviations from ODI — Risk as a first-class element, no job map, no importance/satisfaction scoring, no job-type taxonomy, executor-only scope — are deliberate and documented in `references/elements.md`.

## Gotchas

Non-obvious traps in this framework. Read before diving into a reference file.

- Records (CR, PDR, ADR) live in the four central record directories, never inline or in per-element folders. See [references/structure.md](references/structure.md#record-elements).
- Create produces no Change Record — CRs cover only updates and deletes. See [references/workflows.md](references/workflows.md).
- A CR logs what changed and why; PDRs/ADRs capture the decision and its alternatives — don't merge them into one record.
- Slugs use only the local ID segment; document bodies use fully-qualified logical IDs. See [references/structure.md](references/structure.md#naming).

## References

| File | Load when |
|---|---|
| [references/elements.md](references/elements.md) | Defining, judging, or explaining what something *is* — jobs, outcomes, risks, requirements, architecture, components, or records (CR, PDR, ADR). Use the minimal patterns and good/bad examples as the rubric for drafting and review. |
| [references/structure.md](references/structure.md) | Placing content on disk — choosing a complexity tier, naming a file or folder, resolving a logical ID, or finding the path and template for an element or record. |
| [references/workflows.md](references/workflows.md) | Reading, creating, updating, or deleting intent documents. Covers the read workflow and the draft/judge/review modification process. |
| [references/context-tracing.md](references/context-tracing.md) | Tracing what to load along the spine and laterally before read or any modification — vertical and horizontal traversal rules, and when to load records. |
