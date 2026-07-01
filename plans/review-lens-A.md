# Lens A — Skill craft

Summary

- The skill body itself is compact (45 lines) and well-organized around three reference files, but leaks rubric content that belongs downstream.
- Progressive disclosure has one clear violation: `references/context-tracing.md` is only reachable through `references/workflows.md`, a two-level chain the spec explicitly warns against.
- Frontmatter is compliant; description is well under 1024 chars and includes a "Use when" trigger, though its keyword list can be tightened.
- Several sections restate the same rule two or three times (record placement, context-tracing prelude, alignment guidance), inflating tokens for no clarity gain.
- Templates directory (14 files) matches the element × tier matrix; no drift there.

Findings

FA-01 — `references/context-tracing.md` is a nested reference chain reachable only via `workflows.md`.

- Evidence: `intent/SKILL.md:40-44` links only to `workflows.md`; `intent/references/workflows.md:11,16` is the sole entry point to `context-tracing.md`; spec says "Keep file references one level deep from SKILL.md. Avoid deeply nested reference chains." (`research/agentskills-spec.md:263`).
- Change: move — link `context-tracing.md` directly from SKILL.md's own section pointer block (add a fourth "Context Tracing" pointer sibling to Elements/Structure/Workflows), and drop the wrapper "## Context Tracing" section in `workflows.md`.
- Effort: small
- Depends on: FA-04

FA-02 — SKILL.md "Global Guidance" duplicates rubric content that already lives in `elements.md`.

- Evidence: `intent/SKILL.md:10-17` bullets ("Don't drift into implementation", "Keep outcomes measurable…", "State requirements as what the product must do…", "Keep the trace coherent…") restate rubrics in `intent/references/elements.md:69-71,100-102,130-132` and are formalized as the shared Alignment check at `intent/references/elements.md:7-24`.
- Change: cut — remove the "Global Guidance" section entirely. The two rules not covered by rubrics ("Be candid about unknowns" and "Write plainly so a junior PM or engineer gets it on the first pass") merge into Gotchas as one bullet.
- Effort: trivial
- Depends on: none

FA-03 — SKILL.md intro paragraph paraphrases the frontmatter description.

- Evidence: `intent/SKILL.md:8` ("A framework for structuring product and engineering intent as a traceable hierarchy — from jobs and outcomes through risks and requirements down to architecture and components…") restates the frontmatter at `intent/SKILL.md:3`. The agent has already loaded both.
- Change: cut — delete line 8. Let the section pointers speak.
- Effort: trivial
- Depends on: none

FA-04 — `workflows.md` states the context-tracing prelude twice in adjacent paragraphs.

- Evidence: `intent/references/workflows.md:11` ("Trace context by element before read or any modification. Resolve tiers and paths from structure.md, then trace vertically and horizontally per context-tracing.md. Do not load record types…") is restated as a numbered list at `intent/references/workflows.md:13-17` ("1. Load this skill … 2. Trace context per context-tracing.md … 3. Load Records only when you need…").
- Change: merge — keep only the numbered list under a single `## Context tracing` heading; drop the prose paragraph. If FA-01 lands, remove this section entirely and rely on the SKILL.md pointer.
- Effort: trivial
- Depends on: none

FA-05 — `workflows.md` "Read" is a three-line stub whose only content is "trace context and don't emit records".

- Evidence: `intent/references/workflows.md:19-21` — the section repeats the tracing rule and adds one negative claim about records.
- Change: merge into `## Modification sessions` — add one sentence at the top: "Read is context tracing without modification; no CR, PDR, or ADR is emitted." Delete the standalone heading.
- Effort: trivial
- Depends on: FA-04

FA-06 — Three-session workflow prose is dense and mixes solo/shared paths inline.

- Evidence: `intent/references/workflows.md:25` runs on for six lines mixing "shared/team repo" flow with solo-collapse instructions; `intent/references/workflows.md:27-35` then re-numbers the same steps with the solo caveat re-embedded in step 2.iv. A junior reader has to read twice to extract the actual sequence.
- Change: rewrite: split into two visibly labeled subsections. Example: "`### Shared repo` — Draft → Judge → Review (independent). `### Solo repo` — Draft → Judge (self). Both loops use the same binary checklist: subsection rubric + Alignment check + Records rubric when relevant." Then a compact numbered list per path with no inline "unless you are solo" caveats.
- Effort: small
- Depends on: none

FA-07 — Record-placement rule is stated three times in `structure.md`.

- Evidence: `intent/references/structure.md:25` ("Every CR, PDR, and ADR lives in one of four central directories…never inline, never in per-element folders, never in appendices."), restated at `intent/references/structure.md:85-86`, and again per-record-type at `intent/references/structure.md:171,187,202-208`. The gotcha at `intent/SKILL.md:23` also asserts it.
- Change: cut — keep the single authoritative statement in the `### Record elements` tier table area (line 25 block). In the per-record-type Paths sections, drop the "All PDRs live in `docs/product/drs/` regardless of what they scope" / "All ADRs live in `docs/engineering/drs/` regardless of what they scope" repeats and let the path itself do the work.
- Effort: trivial
- Depends on: none

FA-08 — `elements.md` "Run the Alignment check" and sibling-scoping caveats repeat under every spine element.

- Evidence: `intent/references/elements.md:43,71,102,132,162,170,190` — nearly identical "Run the Alignment check" reminders on Job, Outcome, Risk, Requirement, Engineering, Architecture, Component; sibling-scoping restated inline at `intent/references/elements.md:71,102,132`.
- Change: merge — state the pattern once in the `## Alignment check` intro: "Every element below runs this check with parent/descendants read along the spine (product → job → outcome → risk → requirement → architecture → component)." Then remove per-element "Run the Alignment check" bullets. Keep only per-element deviations (e.g. component sibling-only note).
- Effort: small
- Depends on: none

FA-09 — Frontmatter description trigger list is narrow and repetitive.

- Evidence: `intent/SKILL.md:3` ends with "Use when creating, editing, tracing, or placing these documents and their IDs/paths/tiers." Best practice (`research/agentskills-best-practices.md:278`) says the description "triggers on the right prompts" — this list misses common phrasings like "product docs", "requirements", "ADR/PDR/CR", "job story", "outcome". "IDs/paths/tiers" is skill-internal jargon unlikely to appear in a user prompt.
- Change: rewrite: `description: Structure product and engineering docs as a traceable tree of jobs, outcomes, risks, requirements, architecture, and components, with change records (CR), product decision records (PDR), and architectural decision records (ADR). Use when writing, reviewing, tracing, or placing any of these — including "product docs", "requirements", "ADR", "PDR", "CR", "job story", or "outcome".`
- Effort: trivial
- Depends on: none

FA-10 — Templates directory is referenced from two places with divergent framing.

- Evidence: `intent/references/structure.md:48` names four parent templates as "canonical form" for child rendering; individual per-record-type paths at `intent/references/structure.md:97,109,115,129,147,159,165,177,183,193,199,217,226` link the same templates in a different context. There is no single index of templates in the skill.
- Change: merge — add a compact `## Templates` table at the bottom of `structure.md` listing each of the 14 templates with its element type and tier. Replace per-path "Template: [name.md](../assets/templates/name.md)" links with a table reference or drop the redundant one-liner in the child-rendering prose.
- Effort: small
- Depends on: none

FA-11 — `structure.md` "Complexity Tiers" mixes tier definition, per-element promotion notes, and product profiles in one flow.

- Evidence: `intent/references/structure.md:6-39` — the Spine table (lines 15-20), then a paragraph of per-element promotion exceptions (line 21), then the Record table (lines 27-31), then Product scale profiles (lines 34-39). A reader looking up "when do I promote a requirement to Tier 2?" has to scan three shapes.
- Change: rewrite: keep the two tier tables tight (spine + record). Move per-element promotion notes into the per-element Paths sections (`### Requirements` etc.), so the answer to "when do I promote X?" lives next to X's paths. Keep Product scale profiles where they are; they're audience-level not element-level.
- Effort: small
- Depends on: none

FA-12 — Gotchas #1 and #3 partially overlap with content the reader will hit as soon as they load `elements.md` or `structure.md`.

- Evidence: `intent/SKILL.md:23,25` — the four-directory rule is restated in `intent/references/structure.md:25`; the CR-vs-decision distinction is restated in the CR/PDR/ADR subsections of `intent/references/elements.md:223-287`.
- Change: cut — trim Gotcha #1 to just "Capturing a record never promotes the element it concerns" (the non-obvious half); the "four central directories" fact is load-bearing in structure and doesn't need a preload. Trim Gotcha #3 to "A CR logs what changed and why; PDRs/ADRs capture the decision and its alternatives — don't merge them into one record." (drop the pointer to elements.md — the reader is about to open it).
- Effort: trivial
- Depends on: none
