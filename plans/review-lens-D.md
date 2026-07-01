# Lens D — Internal DRY, clarity, and cross-reference integrity

Summary

- The "records live in four central directories" rule is stated in five places; structure.md must own it and every other site should compress to a pointer.
- The Records subsections in elements.md repeat the same coherence-check boilerplate verbatim three times (CR, PDR, ADR).
- workflows.md opens with philosophy and re-signposts context tracing that context-tracing.md already owns; both preambles can be cut.
- structure.md's "Child rendering" and "Form rules" say the "one line of substance + See link" rule twice in adjacent paragraphs, and the standalone "Change records" section restates the record-elements tier table.
- SKILL.md gotchas duplicate content that lives in structure.md/workflows.md; gotchas should be one-liners that point rather than restate.

Findings

FD-01 — "Records live in four central directories" is restated in five files.

- Evidence:
  - `intent/SKILL.md:23` — "Records (CR, PDR, ADR) live in four central directories under `docs/product/` and `docs/engineering/`, never in appendices and never in per-element record folders."
  - `intent/references/structure.md:25` — "Every CR, PDR, and ADR lives in one of four central directories — `docs/product/crs/`, `docs/product/drs/`, `docs/engineering/crs/`, `docs/engineering/drs/` — never inline, never in per-element folders, never in appendices."
  - `intent/references/structure.md:242` — "Records live in exactly four central directories under `docs/product/` and `docs/engineering/`."
  - `intent/references/context-tracing.md:13` — "Records (CRs, PDRs, ADRs) live in the four central directories — `docs/product/crs/`, `docs/product/drs/`, `docs/engineering/crs/`, `docs/engineering/drs/` — regardless of the element they scope."
  - `intent/references/elements.md:225` — "Every captured CR lives in its own file or folder in one of the two central CR directories — `docs/product/crs/` … or `docs/engineering/crs/` …"
- Change: Canonical home = `structure.md:23-32` (Record elements). Delete the restatement at `structure.md:242`. Rewrite `SKILL.md:23` gotcha to: "Records (CR, PDR, ADR) live in the four central record directories, never inline or in per-element folders. See [references/structure.md](references/structure.md#record-elements)." Rewrite `context-tracing.md:11-13` to: "For record locations and backlink conventions, see [structure.md](structure.md#record-elements)." Rewrite the CR paragraph at `elements.md:225` to drop the directory recitation and keep only "See [structure.md](structure.md#change-records) for where CRs live and how the outcome/component backlinks."
- Effort: small
- Depends on: none

FD-02 — CR/PDR/ADR subsections in elements.md repeat identical coherence-check boilerplate.

- Evidence:
  - `intent/references/elements.md:221` (Records intro) — "Per the coherence check, records at related scope must not contradict each other on the same fact without an explicit superseding record, and where two records touch the same interface or decision surface they agree."
  - `intent/references/elements.md:229` (CR) — near-verbatim: "other change records that bear on the same fact or touch the same interface must not tell conflicting stories, and where they affect shared boundaries they agree on facts."
  - `intent/references/elements.md:251` (PDR) — "product decision records at overlapping or related scope must not contradict without an explicit superseding record, and where decisions overlap they agree on facts and intent."
  - `intent/references/elements.md:273` (ADR) — same sentence with "architectural decision records" swapped in.
- Change: Canonical home = the Records intro at `elements.md:221`. Delete the three per-subsection coherence paragraphs (lines 229, 251, 273). The generic rule at line 221 already covers all three record types.
- Effort: trivial
- Depends on: none

FD-03 — workflows.md preamble is philosophy, not procedure.

- Evidence: `intent/references/workflows.md:3-5` — "The intent documentation forms a tree … The discipline is familiar from engineering practice; the gap was usually on the intent side … Assistants make it practical …" then "The verbs create, read, update, and delete describe what you do to elements and files. Read is the workflow for loading context …"
- Change: Cut lines 3-5 entirely. The SKILL.md description already frames purpose; this section should open on procedure. Replace with a single line: "Focus for any workflow is one or more elements defined in [elements.md](elements.md); paths and tiers are in [structure.md](structure.md)."
- Effort: trivial
- Depends on: none

FD-04 — workflows.md `## Context Tracing` section duplicates its target and the "For any workflow" list restates what the Read/Modification sections already say.

- Evidence:
  - `intent/references/workflows.md:9-17` — a Context Tracing heading whose body is "Trace context by element before read or any modification. Resolve tiers and paths from structure.md, then trace vertically and horizontally per context-tracing.md" plus a 3-step "Load this skill / Trace context / Load records only when needed" list.
  - `intent/references/context-tracing.md:1-19` already owns the tracing rules and record-loading guidance.
  - The 3-step list is then re-invoked at `workflows.md:21` ("Trace context as in Context Tracing") and inside Modification sessions at `workflows.md:29-30` ("Trace context, then work with an assistant…"; "Trace context again, then run…").
- Change: Delete `workflows.md:9-17` outright. Move the "Load Records only when you need that history" rule into the top of `context-tracing.md` (canonical home for what to trace and when). In workflows.md, replace every "Trace context (as in Context Tracing)" reference with a bare link to `context-tracing.md`.
- Effort: small
- Depends on: FD-03

FD-05 — structure.md states the "one line of substance + See link" rule twice.

- Evidence:
  - `intent/references/structure.md:46` — "**Tier 1 or Tier 2 child — reference.** The parent keeps the child's heading or bullet, one line of substance (the child's minimal pattern, or the Responsibility line for components), and a `See [<Logical ID> - <Name>](path).` link to the canonical document."
  - `intent/references/structure.md:54` (first bullet of Form rules) — "A reference-form child shows exactly one line of substance — the element's minimal pattern, or the Responsibility line for components — followed by the See link. No paragraphs, no sub-lists…"
- Change: Merge into a single statement at line 46. Delete the duplicate bullet at line 54; keep only the "no paragraphs, no sub-lists" clarifier as an appended sentence at line 46.
- Effort: trivial
- Depends on: none

FD-06 — structure.md's `### Change records` section restates the Record elements tier table and adds only the split-directory rule.

- Evidence:
  - `intent/references/structure.md:23-32` — Record elements table already defines Tier 1/2 shapes for all record types.
  - `intent/references/structure.md:201-226` — "Change records" restates: "Tier 1 — own document / Tier 2 — own directory" with the same shape, plus the two-directory split.
  - The PDR and ADR sections (lines 169-199) have the same shape as well but at least add the "single directory" fact.
- Change: In the Change records section keep only (a) the two-directory selector, (b) the double-touch rule at line 208 ("recorded once on the side most central to the change, and the other side's affected element(s) still backlink"), and (c) the two paths. Delete the Tier 1/Tier 2 sub-headings; a one-line "Templates: cr-document.md (Tier 1), cr-directory.md (Tier 2). See [Record elements](#record-elements) for tier selection." is enough. Apply the same trim to PDR/ADR sections.
- Effort: small
- Depends on: none

FD-07 — Modification-sessions preamble sentence duplicates the numbered steps that follow it.

- Evidence: `intent/references/workflows.md:25` — "The three-session flow below is for a shared/team repo with team review: Draft, then Judge (validation loop), then commit/push/open a PR or MR, then independent Review. A solo author, or any repo without a second reviewer, collapses it to draft, then self-run the same Judge validation loop, then finalize/commit — no PR/MR or separate reviewer required, but the same binary checklist … must still pass before finalizing." Immediately followed by numbered steps 1-3 (`workflows.md:27-35`) that re-tell the same sequence.
- Change: Delete the preamble at line 25. Move the "solo author collapses steps 1–2 and skips step 3" note to a single sentence at the end of the numbered list, e.g. "Solo authors run Draft and Judge, then commit; skip Review."
- Effort: trivial
- Depends on: none

FD-08 — "Creating a brand-new element produces no CR" stated in two places.

- Evidence:
  - `intent/SKILL.md:24` (gotcha) — "Creating a brand-new element produces no Change Record — there is no prior state. CRs are only for updates and deletes."
  - `intent/references/workflows.md:37` — "For create, there is no prior element state, so there is no change in the sense a CR describes; create does not produce a Change Record."
- Change: Canonical home = `workflows.md:37`. Rewrite `SKILL.md:24` to: "Create produces no Change Record — CRs cover only updates and deletes. See [references/workflows.md](references/workflows.md)."
- Effort: trivial
- Depends on: none

FD-09 — "Capturing a record never promotes the element" stated in two places.

- Evidence:
  - `intent/SKILL.md:23` (gotcha) — "Capturing a record never promotes the element it concerns."
  - `intent/references/structure.md:25` — "**Capturing a record never promotes the element it concerns.**"
- Change: Canonical home = `structure.md:25` (already bolded there). Delete the clause from the SKILL.md gotcha; the FD-01 rewrite of that gotcha already points to `structure.md#record-elements`, which carries the promotion rule.
- Effort: trivial
- Depends on: FD-01

FD-10 — Slug-vs-logical-ID gotcha restates the Naming section.

- Evidence:
  - `intent/SKILL.md:26` — "Slug vs logical ID: filenames/folders (slugs) use only the local ID segment (e.g. `O001-<name>`); cross-tree references in document bodies use fully-qualified logical IDs (e.g. `O001-R001`). Don't repeat ancestor IDs in slugs."
  - `intent/references/structure.md:63-83` — Naming section with the same rules and a full slug/logical-ID table.
- Change: Rewrite `SKILL.md:26` to a pointer: "Slugs use only the local ID segment; document bodies use fully-qualified logical IDs. See [references/structure.md](references/structure.md#naming)."
- Effort: trivial
- Depends on: none

FD-11 — context-tracing.md uses a label that doesn't match the templates.

- Evidence:
  - `intent/references/context-tracing.md:13` — "follow the links in the scoped element's `## See Also` (Tier 1+) or its inline `**Decision Records**` / `**Change Records**` lists (Tier 0)."
  - Templates use `**Product Decision Records**` and `**Architectural Decision Records**` (never bare `**Decision Records**`): `intent/assets/templates/product.md:37`, `intent/assets/templates/outcome-document.md:24`, `intent/assets/templates/architecture.md:38, 62`.
- Change: Rewrite the offending clause to: "follow the links in the scoped element's `## See Also` (Tier 1+) or its inline `**Product/Architectural Decision Records**` and `**Change Records**` lists (Tier 0)."
- Effort: trivial
- Depends on: none

FD-12 — context-tracing.md opens with an empty pointer header.

- Evidence: `intent/references/context-tracing.md:1-3` — "# Context Tracing / Tracing rules for [workflows.md](workflows.md). Tiers, paths, and logical IDs are defined in [structure.md](structure.md)."
- Change: Delete "Tracing rules for [workflows.md]." — the reader arrived here from workflows and needs no reminder. Keep only the "Tiers, paths, and logical IDs are defined in [structure.md](structure.md)." sentence.
- Effort: trivial
- Depends on: none

FD-13 — Judge step 1 embeds a long parenthetical that duplicates the Records intro.

- Evidence: `intent/references/workflows.md:31` — "When the draft includes or implies decision or history artifacts, also check it against [Records](elements.md#records) for what each type is for: [Change Record (CR)](elements.md#change-record-cr), [Product Decision Record (PDR)](elements.md#product-decision-record-pdr), and [Architectural Decision Record (ADR)](elements.md#architectural-decision-record-adr)."
- Change: Compress to "When the draft touches decisions or history, also validate against [Records](elements.md#records)." The three per-type deep links add nothing the Records section doesn't already surface.
- Effort: trivial
- Depends on: none

FD-14 — Sibling-scoping riders in elements.md restate the Alignment check for each element.

- Evidence: `intent/references/elements.md:71` (outcome), `:102` (risk), `:132` (requirement) — each has a "Run the Alignment check … For sibling scoping, [restatement of minimal-overlap/agree-on-contact/cover-the-parent with element-specific words]" rider that repeats the three sibling checks defined at `elements.md:16-20`.
- Change: Keep only the rider clauses that add net-new element-specific facts (e.g. requirement's "cousins traced through shared architecture stay compatible"). Delete the redescriptions of minimal-overlap/agree/cover — the reader already followed the link.
- Effort: small
- Depends on: none

FD-15 — SKILL.md's three "Load when …" sections are ceremony around a link.

- Evidence: `intent/SKILL.md:28-44` — three sections (Elements, Structure, Workflows) each with a one-paragraph "Load when …" gloss and a bare "See [references/x.md]." Same shape three times.
- Change: Collapse into a single "References" section with a three-row table: `File | Load when`. Rows: elements — defining/judging what something is; structure — placing content on disk, tiers, naming; workflows — reading, creating, updating, deleting docs including context tracing. Removes ~15 lines and gives the reader a single scannable index.
- Effort: small
- Depends on: none
