# Workflows

Focus for any workflow is one or more elements defined in [elements.md](elements.md); paths and tiers are in [structure.md](structure.md). Before any read or modification, trace context per [context-tracing.md](context-tracing.md).

## Modification sessions

Read is context tracing without modification; no CR, PDR, or ADR is emitted.

For create, update, or delete, work across sessions with clear handoffs.

### Shared repo

1. Drafting session — Trace context per [context-tracing.md](context-tracing.md), then work with an assistant to draft the change (new or edited documents, and any new record files you intend to add).
2. Judge session — Trace context again, then run a validation loop on the draft that already exists from the drafting session:
   1. Validate the draft against a binary checklist — every check is pass/fail. The checklist is the good/bad rubric plus the Minimal Pattern from the subsection for the relevant element type under [elements.md](elements.md), plus the shared [Alignment check](elements.md#alignment-check). When the draft touches decisions or history, also validate against [Records](elements.md#records).
   2. Fix the failures.
   3. Repeat the validate-and-fix cycle until every check passes.
   4. Only then finalize: commit, push, and open a PR or MR.
3. Review session — An independent reviewer either reviews the PR/MR in the host, or checks out the branch locally for an agent-assisted review. They trace context, run the same validation loop, then combine the assistant's feedback with their own judgment and leave suggestions for the author.

### Solo repo

Draft and Judge only, using the same binary checklist. Finalize and commit at the end of Judge; no PR/MR and no separate reviewer.

For update or delete, consider impact on all descendants of the changed element in the trace, and adjust or verify downstream documents so the tree stays coherent. For create, there is no prior element state, so there is no change in the sense a CR describes; create does not produce a Change Record. For update and delete, add a [Change Record (CR)](elements.md#change-record-cr) scoped to the modification when the change is material to the success of outcomes; trivial edits do not need one. [PDRs](elements.md#product-decision-record-pdr) and [ADRs](elements.md#architectural-decision-record-adr) may be added during any modification when a product or engineering decision should be captured; they are optional and depend on whether a decision worth recording was made. See [structure.md](structure.md#record-elements) for record tier selection.
