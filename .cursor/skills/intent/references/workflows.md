# Workflows

The intent documentation forms a tree: product and engineering elements, plus records. The discipline is familiar from engineering practice; the gap was usually on the intent side, where it is easy to drift straight into implementation. Assistants make it practical to draft, check, and review intent documents the same way you would code.

The verbs create, read, update, and delete describe what you do to elements and files. Read is the workflow for loading context and building a shared picture without changing documents. Create, update, and delete are modifications and follow the three-session workflow below (unless you are only exploring, in which case use read).

Focus for any workflow is one or more elements among product (including jobs), outcomes (including risks), requirements, architecture, and components — as defined in [Elements](../SKILL.md#elements) and located on disk in [Files](../SKILL.md#files) and [Structure](structure.md).

## Loading context

For any workflow below, loading appropriate context means:

1. Load this skill (`intent/SKILL.md`).
2. For the element under discussion, trace context vertically and horizontally as described in [Context Tracing](context-tracing.md).

Load [Records](../SKILL.md#records) (CRs, PDRs, ADRs) only when you need that history or those decisions for the task.

## Read

Use read when you want an assistant to build an understanding of an element without changing documents. Load context as above. Read does not emit change records or decision records.

## Modification sessions

For create, update, or delete, work across three sessions with clear handoffs:

1. Drafting session — Load context, then work with an assistant to draft the change (new or edited documents, and any new record files you intend to add).
2. Judge session — Load context again, then prompt an assistant to judge the draft using the subsection for the relevant element type under [Elements](../SKILL.md#elements) as the rubric. When the draft includes or implies decision or history artifacts, also use [Records](../SKILL.md#records) as the reference for what each type is for: [Change Record (CR)](../SKILL.md#change-record-cr), [Product Decision Record (PDR)](../SKILL.md#product-decision-record-pdr), and [Architectural Decision Record (ADR)](../SKILL.md#architectural-decision-record-adr). The author incorporates feedback, then commits, pushes, and opens a pull or merge request.
3. Review session — An independent reviewer either reviews that PR/MR in the host, or checks out the branch locally for an agent-assisted review. They load context, run the same judge prompt against that subsection (and [Records](../SKILL.md#records) when relevant), then combine the assistant's feedback with their own judgment and leave suggestions for the author.

For update or delete, consider impact on all descendants of the changed element in the trace, and adjust or verify downstream documents so the tree stays coherent. For create, there is no prior element state, so there is no change in the sense a CR describes; create does not produce a Change Record. For update and delete, add a [Change Record (CR)](../SKILL.md#change-record-cr) scoped to the modification when the change is material to the success of outcomes; trivial edits do not need one. [PDRs](../SKILL.md#product-decision-record-pdr) and [ADRs](../SKILL.md#architectural-decision-record-adr) may be added during any modification when a product or engineering decision should be captured; they are optional and depend on whether a decision worth recording was made.
