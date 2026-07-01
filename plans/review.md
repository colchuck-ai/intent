# Intent skill review — implementation plan

Progress

- **Done**: T01, T03, T04, T05, T06, T07 (session 2026-07-01); T02, T08, T09 (session 2026-07-01, cont.). All nine tasks complete.
- **Pending**: none.
- **Deviation from plan**: T05's rewrite of the outcome minimal pattern (`[Verb] [Unit of Measure] [Object]` → `[Direction] [Metric] [Object] [Context]`) and the job pattern's third clause (`[Outcome]` → `[Benefit]`) was propagated into three templates that quote those patterns verbatim (`assets/templates/product.md`, `outcome-document.md`, `outcome-directory.md`). The plan for T05 only named `elements.md`; the templates followed as a necessary consistency fix so drafters aren't handed the old vocabulary. No task in the DAG covered this, so it's called out here.
- **Residual finding not addressed**: Lens D's FD-01 flagged five sites of the record-directory rule; T02 (SKILL.md), T04 (elements.md), and T07 (structure.md) collapsed three. The recitation at `context-tracing.md:13` remains — it wasn't in T09's spec (which handled only the label mismatch and preamble drop). Not a bug in the plan, but a candidate for a small follow-up if further DRY is wanted.

Executive summary

- **Description frontmatter is the highest-ROI single edit.** It's descriptive ("Build and maintain…") not imperative, uses skill-internal jargon ("IDs/paths/tiers"), and has no negative boundary; a one-line rewrite improves trigger rate and disambiguation for a trivial cost. (T01)
- **SKILL.md is doing double duty** as intro, rubric, and pointer table. Slim it to a purpose sentence, a lineage note, gotchas-as-pointers, and a single References table; also declare the ODI/Christensen mixture that today is left implicit. (T02, T03)
- **elements.md repeats the Alignment check and record-directory rules under every element and every record type.** Consolidate to one canonical statement per section top and delete the per-element and per-record-type restatements. (T04)
- **Product-side content silently blends Ulwick outcomes with Christensen job stories,** has a definitional collision on "Outcome" (job pattern's third clause vs. the Outcome element), and doesn't declare its deliberate departures from ODI (Risk as a first-class element, no job map, no importance/satisfaction scoring). A targeted set of elements.md edits aligns vocabulary and makes deviations explicit. (T05)
- **Requirement traceability is one-directional.** Outcomes link forward to risks/requirements; requirement templates never backlink to the risks or outcomes they serve. Add `## Mitigates` to both requirement templates. (T06)
- **structure.md restates the record-directory rule three times and mixes tier definition with per-element promotion notes.** Consolidate the rule, move promotion into per-element Paths sections, dedupe the child-rendering rule, and add a compact Templates index. (T07)
- **workflows.md hides `context-tracing.md` behind a nested reference chain** (SKILL.md → workflows.md → context-tracing.md), violating the spec's one-level-deep rule; promote context-tracing.md to a first-class pointer from SKILL.md, cut the workflows.md philosophy preamble and its wrapper Context Tracing section, and split the modification-sessions prose into visible shared/solo subsections. (T08, T09)

Tasks

T01 — Rewrite the SKILL.md description frontmatter  [DONE]

- Files: `intent/SKILL.md`
- Findings: FB-01, FB-02, FB-03, FB-04, FB-05, FB-06; also FA-09
- Change: Replace `intent/SKILL.md:3` with Lens B's Alt 1 (recommended default; ~790 chars, ≤ 1024). It leads imperative ("Use when…"), keeps the two `docs/` scope signals, spells out all seven elements and three record acronyms, lists six paraphrase triggers ("write a product spec", "capture why we chose X", "log this change", "record this decision", "draft an ADR/PDR/CR", "define the customer's job"), covers workflow-term triggers (draft/judge/review, alignment check), and closes with a negative boundary ("Skip for code changes, READMEs, marketing copy, or ad-hoc notes that don't link into the intent tree"). Retain Alt 2 and Alt 3 in review-lens-B.md as evaluation candidates.
- Effort: trivial
- Depends on: none

T02 — Collapse SKILL.md gotchas into thin pointers  [DONE]

- Files: `intent/SKILL.md`
- Findings: FA-12, FD-01 (SKILL.md side), FD-08, FD-09, FD-10
- Change: Rewrite each gotcha as a one-liner + link. (1) "Records (CR, PDR, ADR) live in the four central record directories, never inline or in per-element folders. See [references/structure.md](references/structure.md#record-elements)." — the "Capturing a record never promotes the element" clause is deleted from SKILL.md since it's canonically bolded at `structure.md:25`. (2) "Create produces no Change Record — CRs cover only updates and deletes. See [references/workflows.md](references/workflows.md)." (3) "A CR logs what changed and why; PDRs/ADRs capture the decision and its alternatives — don't merge them into one record." (drop the elements.md pointer; the reader is about to open it via T03). (4) Slug gotcha becomes: "Slugs use only the local ID segment; document bodies use fully-qualified logical IDs. See [references/structure.md](references/structure.md#naming)."
- Effort: trivial
- Depends on: T03 (same file — sequence to avoid churn), T07 (structure canonical sections must exist)

T03 — Slim the SKILL.md body: intro, guidance → lineage, References table  [DONE]

- Files: `intent/SKILL.md`
- Findings: FA-02, FA-03, FC-02, FC-13, FD-15
- Change: Cut the intro paragraph at `SKILL.md:8` (it paraphrases the frontmatter). Delete the "Global Guidance" bulleted section: its rubric bullets duplicate `elements.md` rubrics; replace with a compact "Principles" block containing only (a) the two orphans — "Be candid about unknowns" and "Write plainly so a junior PM or engineer gets it on the first pass" — and (b) a new lineage sentence: "Product-side vocabulary follows Ulwick/ODI (jobs, outcomes, needs-first sequencing); job syntax follows the Christensen job-story form. Deviations from ODI — Risk as a first-class element, no job map, no importance/satisfaction scoring, no job-type taxonomy, executor-only scope — are deliberate and documented in `references/elements.md`." Collapse the three "Load when …" sections into a single "References" table with rows for elements, structure, workflows, and — as a fourth first-class pointer — context-tracing (this promotion is what unlocks T08).
- Effort: small
- Depends on: T01

T04 — elements.md: dedupe Alignment / coherence / record-directory repetition  [DONE]

- Files: `intent/references/elements.md`
- Findings: FA-08, FD-02, FD-14; also FD-01 (elements.md side)
- Change: In the `## Alignment check` intro (currently around lines 7-24), add a single line: "Every element below runs this check with parent/descendants read along the spine (product → job → outcome → risk → requirement → architecture → component); per-element sections note only element-specific deviations." Then delete every "Run the Alignment check" reminder at `elements.md:43,71,102,132,162,170,190`, and delete the redescribed sibling-scoping riders at `:71,:102,:132` (keep only rider clauses that add net-new element-specific facts, e.g. requirement's "cousins traced through shared architecture stay compatible"). In `## Records`, keep the coherence statement at `elements.md:221` as canonical; delete the near-verbatim coherence sentences at `:229` (CR), `:251` (PDR), `:273` (ADR). Rewrite the CR paragraph at `elements.md:225` to drop the directory recitation and keep only "See [structure.md](structure.md#change-records) for where CRs live and the outcome/component backlink rules."
- Effort: small
- Depends on: none

T05 — elements.md: JTBD/ODI content alignment and declared deviations  [DONE — see Progress note on template propagation]

- Files: `intent/references/elements.md`
- Findings: FC-01, FC-03, FC-04, FC-05, FC-06, FC-07, FC-08, FC-09, FC-10, FC-12
- Change: (1) Rename job pattern's third clause from "Outcome" to "Benefit" at `:47-51` (removes name collision with the Outcome element). (2) Rewrite the Job section preamble at `:41-43` to name the **job executor** as the person the job belongs to, and declare scope: "core functional jobs only; no job map, no related/emotional/consumption/purchase jobs." (3) Replace the current job "good" example at `:57-59` with a core-functional-job example at correct altitude (not the emotional "feel confident" wording). (4) Extend the Outcome pattern at `:75-79` from `[Verb] [Unit] [Object]` to `[Direction] [Metric] [Object] [Context]`, broaden the metric list beyond time/number/likelihood (rate, cost, effort, etc.), and update good examples at `:85-87` to include a context clause. (5) Append to `:69`: "Like the job it serves, an outcome is stable across changes in the solutions used to get the job done. Prioritization by importance/satisfaction (ODI's opportunity score) is out of scope for this framework." (6) Rewrite the Risk section opening at `:100` to state the ODI deviation explicitly: "Risks are this framework's departure from ODI. Where ODI expresses unmet needs as underserved outcomes (importance − satisfaction), we express them as risks (likelihood × impact) to make the outcome→requirement mitigation trace legible to engineering." (7) Rewrite the Requirement definition at `:130-131` to require risk-backlink: "A requirement is what the product must do or respect to mitigate one or more risks. A requirement that cannot be traced to a named risk is a feature request, not a requirement." (8) Sharpen the Cover-the-parent alignment bullet at `:20` (and its invocations at `:71,:132`): "…do the siblings together *completely enumerate* what better looks like for the parent — every customer-defined way of judging success?"
- Effort: medium
- Depends on: T04 (same file — sequence to avoid stepping on the same sections)

T06 — Add `## Mitigates` to both requirement templates  [DONE]

- Files: `intent/assets/templates/requirement-document.md`, `intent/assets/templates/requirement-directory.md`
- Findings: FC-11
- Change: Immediately after the title/pattern in each template, add a `## Mitigates` section that lists the risks (with logical IDs) and, through them, the outcomes served — e.g. `- **O<NNN>-RSK<NNN>** - <Risk Name>  (serves **O<NNN>** - <Outcome Name>)`. Makes requirement-side traceability symmetric with the outcome-side backlinks that already exist in `outcome-document.md` and `outcome-directory.md`.
- Effort: trivial
- Depends on: T05 (requirement definition rewrite establishes the semantics)

T07 — structure.md: consolidate record-directory rule, tier tables, child-rendering rule, add Templates index  [DONE]

- Files: `intent/references/structure.md`
- Findings: FA-07, FA-10, FA-11, FD-05, FD-06; also FD-01 (structure.md side)
- Change: Keep the canonical statement of the four-directory rule at `structure.md:23-32` (Record elements). Delete the restatement at `structure.md:242`. In the per-record-type Paths sections at `:171,:187,:202-208` drop the "All PDRs live in…/All ADRs live in…" restatements — the path itself carries the fact. Merge the two "one line of substance + See link" statements at `:46` and `:54` into a single statement at `:46`, appending the "no paragraphs, no sub-lists" clarifier. In the `### Change records` section (`:201-226`), drop the Tier 1/Tier 2 sub-headings and replace with a one-line "Templates: `cr-document.md` (Tier 1), `cr-directory.md` (Tier 2). See [Record elements](#record-elements) for tier selection." Apply the same trim to the PDR and ADR sections at `:169-199`. Move per-element promotion notes (currently a paragraph at `:21`) into the per-element Paths sections so the answer to "when do I promote X?" lives next to X. Add a `## Templates` table at the bottom listing all 14 templates with element type and tier; replace the per-path "Template: [name.md]" one-liners with a reference to that table.
- Effort: medium
- Depends on: none

T08 — workflows.md: cut philosophy, remove wrapper Context Tracing section, split modification sessions, compress judge step  [DONE]

- Files: `intent/references/workflows.md`
- Findings: FA-01, FA-04, FA-05, FA-06, FD-03, FD-04, FD-07, FD-13
- Change: Delete the philosophy preamble at `:3-5`; open on procedure with "Focus for any workflow is one or more elements defined in [elements.md](elements.md); paths and tiers are in [structure.md](structure.md)." Delete the entire wrapper `## Context Tracing` section at `:9-17` (the direct SKILL.md pointer to context-tracing.md added in T03 replaces it; broken references become bare links to `context-tracing.md`). Merge the three-line `## Read` section at `:19-21` into the top of `## Modification sessions` as one sentence: "Read is context tracing without modification; no CR, PDR, or ADR is emitted." Delete the mixed shared/solo preamble at `:25`. Split modification sessions into two explicitly labeled subsections — `### Shared repo` (Draft → Judge → Review) and `### Solo repo` (Draft → Judge) — with the same binary checklist referenced once; move the solo-collapse note to a single trailing sentence rather than embedding it in step 2.iv. Compress Judge step 1's parenthetical at `:31` to "When the draft touches decisions or history, also validate against [Records](elements.md#records)."
- Effort: medium
- Depends on: T03 (SKILL.md's References table must already promote context-tracing.md to a direct pointer)

T09 — context-tracing.md: fix label mismatch, drop empty preamble, absorb records-loading rule  [DONE]

- Files: `intent/references/context-tracing.md`
- Findings: FD-11, FD-12; also FD-04 (context-tracing side)
- Change: Delete "Tracing rules for [workflows.md]." at `:1-3`; keep only "Tiers, paths, and logical IDs are defined in [structure.md](structure.md)." Rewrite the label at `:13` from `**Decision Records**` / `**Change Records**` to `**Product/Architectural Decision Records**` and `**Change Records**`, matching the templates (`product.md:37`, `outcome-document.md:24`, `architecture.md:38,62`). Move the "Load Records only when you need that history" rule (deleted from workflows.md in T08) into the top of context-tracing.md as the canonical statement of "what to trace and when."
- Effort: trivial
- Depends on: T08 (T08 removes the rule from workflows.md; T09 accepts it here)

DAG

```
┌─────┐        ┌─────┐        ┌─────┐
│     │        │     │        │     │
│ T01 │        │ T04 │        │ T07 │
│     │        │     │        │     │
└──┬──┘        └──┬──┘        └──┬──┘
   │              │              │
   ▼              ▼              │
┌─────┐        ┌─────┐           │
│     │        │     │           │
│ T03 │        │ T05 │           │
│     │        │     │           │
└──┬──┘        └──┬──┘           │
   │              │              │
   ├──────┐       ▼              │
   │      │    ┌─────┐           │
   │      │    │     │           │
   │      │    │ T06 │           │
   │      │    │     │           │
   │      │    └─────┘           │
   │      │                      │
   ▼      ▼                      ▼
┌─────┐  ┌───────────────────────────┐
│     │  │                           │
│ T08 │  │           T02             │
│     │  │                           │
└──┬──┘  └───────────────────────────┘
   │
   ▼
┌─────┐
│     │
│ T09 │
│     │
└─────┘
```

Edges (for source-view clarity):

- T01 → T03 (same file, sequence to avoid churn)
- T03 → T02 (same file), T03 → T08 (T03 introduces the direct context-tracing pointer in SKILL.md)
- T07 → T02 (T02's gotcha pointers reference canonical sections T07 establishes)
- T04 → T05 (same file, sequence)
- T05 → T06 (T05's requirement definition sets the semantics T06's `## Mitigates` implements)
- T08 → T09 (T08 hands off the records-loading rule that T09 receives)

Deferred / rejected findings

- None. Every finding across the four lenses is adopted into a task. Overlapping findings are consolidated: FA-09 folds into T01 (Lens B has stronger evidence and the concrete alternatives); the record-directory duplication finding FD-01 is a cross-file finding present in T02 (SKILL.md side), T04 (elements.md side), and T07 (structure.md side, where the canonical statement lives); FD-04 spans T08 (source of the rule move) and T09 (destination).
