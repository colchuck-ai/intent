# Workflows

The intent documentation forms a tree: product and engineering elements, plus records. The discipline is familiar from engineering practice; the gap was usually on the intent side, where it is easy to drift straight into implementation. Assistants make it practical to draft, check, and review intent documents the same way you would code.

The verbs create, read, update, and delete describe what you do to elements and files. Read is the workflow for loading context and building a shared picture without changing documents. Create, update, and delete are modifications and follow the three-session workflow below (unless you are only exploring, in which case use read).

Focus for any workflow is one or more elements among product (including jobs), outcomes (including risks), requirements, architecture, and components — as defined in [elements.md](elements.md). Naming, paths, and tiers are in [structure.md](structure.md).

## Context Tracing

Trace context by **element** before read or any modification. Resolve tiers and paths from [structure.md](structure.md), then trace vertically and horizontally per [context-tracing.md](context-tracing.md). **Do not** load record types (CRs, PDRs, ADRs) until you need that history.

For any workflow below:

1. Load this skill (`intent/SKILL.md`).
2. Trace context per [context-tracing.md](context-tracing.md).
3. Load [Records](elements.md#records) (CRs, PDRs, ADRs) only when you need that history or those decisions for the task.

## Read

Use read when you want an assistant to build an understanding of an element without changing documents. Trace context as in [Context Tracing](#context-tracing). Read does not emit change records or decision records.

## Modification sessions

The three-session flow below is for a shared/team repo with team review: Draft, then Judge (validation loop), then commit/push/open a PR or MR, then independent Review. A solo author, or any repo without a second reviewer, collapses it to draft, then self-run the same Judge validation loop, then finalize/commit — no PR/MR or separate reviewer required, but the same binary checklist (element rubric + [Alignment check](elements.md#alignment-check) + [Records](elements.md#records) when relevant) must still pass before finalizing.

For create, update, or delete, work across these sessions with clear handoffs:

1. Drafting session — Trace context, then work with an assistant to draft the change (new or edited documents, and any new record files you intend to add).
2. Judge session — Trace context again, then run a validation loop on the draft that already exists from the drafting session:
   1. Validate the draft against a binary checklist — every check is pass/fail. The checklist is the good/bad rubric plus the Minimal Pattern from the subsection for the relevant element type under [elements.md](elements.md), plus the shared [Alignment check](elements.md#alignment-check). When the draft includes or implies decision or history artifacts, also check it against [Records](elements.md#records) for what each type is for: [Change Record (CR)](elements.md#change-record-cr), [Product Decision Record (PDR)](elements.md#product-decision-record-pdr), and [Architectural Decision Record (ADR)](elements.md#architectural-decision-record-adr).
   2. Fix the failures.
   3. Repeat the validate-and-fix cycle until every check passes.
   4. Only then finalize: commit (and, if the repo uses pull/merge requests, push and open a PR or MR). A solo author stops here.
3. Review session (shared repo / second person only) — An independent reviewer either reviews that PR/MR in the host, or checks out the branch locally for an agent-assisted review. They trace context, run the same validation loop (the same binary checklist: the subsection rubric, [Alignment check](elements.md#alignment-check), and [Records](elements.md#records) when relevant), then combine the assistant's feedback with their own judgment and leave suggestions for the author.

For update or delete, consider impact on all descendants of the changed element in the trace, and adjust or verify downstream documents so the tree stays coherent. For create, there is no prior element state, so there is no change in the sense a CR describes; create does not produce a Change Record. For update and delete, add a [Change Record (CR)](elements.md#change-record-cr) scoped to the modification when the change is material to the success of outcomes; trivial edits do not need one. [PDRs](elements.md#product-decision-record-pdr) and [ADRs](elements.md#architectural-decision-record-adr) may be added during any modification when a product or engineering decision should be captured; they are optional and depend on whether a decision worth recording was made. See [structure.md](structure.md#choosing-a-tier) for record tier selection.
