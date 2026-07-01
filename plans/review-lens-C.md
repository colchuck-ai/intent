# Lens C — Content fidelity to JTBD/ODI

Summary

- The skill claims JTBD lineage but silently mixes Christensen-style job stories with Ulwick vocabulary; nowhere does it declare which school it follows or where it deliberately deviates.
- The job pattern's trailing "so I can [Outcome]" collides head-on with the framework's own Outcome element, seeding definitional confusion at the top of the tree.
- "Risk" is a non-ODI element inserted between outcomes and requirements; ODI would describe the same territory as underserved outcomes / opportunities. The insertion is defensible but unjustified in text.
- Outcome syntax drops ODI's contextual clarifier and narrows "performance metric" to time/number/likelihood, encouraging force-fit statements without explaining the narrowing.
- Requirements are declared "solution-aware" and template has no backlink to the outcomes/risks they serve — the outcome→requirement trace only exists on the outcome side, breaking needs-first symmetry.

Findings

FC-01 — Job pattern's "so I can [Outcome]" clause collides with the Outcome element and muddies the whole vocabulary

- Evidence: `intent/references/elements.md:47-51` — `> When [Situation], I want to [Goal], so I can [Outcome].` with "Outcome: why it matters — the benefit or result they expect"; `intent/references/elements.md:67-79` — Outcome element is a measurable success metric. Same word, two meanings on adjacent pages.
- Change: rewrite the job pattern's third clause label from "Outcome" to "Benefit" (or "Result"): `> When [Situation], I want to [Goal], so I can [Benefit].` Update the sub-bullet on line 51 to "Benefit: why it matters — the result they expect from achieving the goal."
- Effort: trivial
- Depends on: none

FC-02 — Skill silently mixes Christensen/Klement job stories with Ulwick vocabulary; declare the lineage or align

- Evidence: `intent/references/elements.md:47` uses the classic Christensen "When…I want to…so I can…" job story; `research/jtbd-outcome-driven-innovation.md:72-80` — Ulwick's ODI pattern is `[verb] + [object] + [optional contextual clarifier]` ("Pass on life lessons to children"). The skill's product-side vocabulary (solution-free, stable, functional, outcome-metric) is Ulwick's, but the job syntax is not.
- Change: rewrite `intent/SKILL.md:6-9` overview to add one sentence stating the mixture, e.g. "This framework uses Ulwick/ODI-style outcomes and Christensen-style job stories; it does not adopt ODI's job map, opportunity scoring, or job-type taxonomy." Alternatively align the job pattern to Ulwick's verb-object form.
- Effort: trivial
- Depends on: none

FC-03 — "Risk" is not an ODI element; either rename or explicitly frame as a deliberate departure

- Evidence: `intent/references/elements.md:98-127` defines Risk as a first-class element between Outcome and Requirement, measured by "likelihood × impact"; ODI has no such element — the equivalent is `underserved outcome` (importance minus satisfaction) per `research/jtbd-outcome-driven-innovation.md:229-241`.
- Change: rewrite the opening paragraph of the Risk section (`elements.md:100`) to explicitly say: "Risks are this framework's departure from ODI. Where ODI expresses unmet needs as underserved outcomes (importance − satisfaction), we express them as risks (likelihood × impact) to make the outcome→requirement mitigation trace legible to engineering." Keep the element; just declare the deviation.
- Effort: small
- Depends on: none

FC-04 — Outcome syntax drops ODI's contextual clarifier and narrows "performance metric"

- Evidence: `intent/references/elements.md:75-79` — `[Verb] [Unit of Measure] [Object]` with Verb restricted to minimize/maximize and Unit restricted to time/number/likelihood; `research/jtbd-outcome-driven-innovation.md:155-164` — ODI pattern is `Direction of improvement + Performance metric + Object of control + Contextual clarifier`, e.g. "Minimize the time it takes to get the songs **in the desired order for listening**."
- Change: rewrite the pattern in `elements.md:75-79` to `[Direction] [Metric] [Object] [Context]` with sub-bullets: "Direction: minimize/maximize. Metric: performance measure (time, count, likelihood, rate, cost, effort…). Object: the specific thing being measured. Context: the situation the outcome applies to." Update the good examples on 85-87 to include an explicit context clause.
- Effort: small
- Depends on: none

FC-05 — Job "good" example blends emotional job into a core-job frame at wrong altitude — the exact anti-pattern ODI warns against

- Evidence: `intent/references/elements.md:57` — "When I'm comparing options before a big purchase, I want to feel confident I'm not missing a better alternative, so I can buy without second-guessing." Compare `research/jtbd-outcome-driven-innovation.md:82-92` — core functional job should be defined at correct altitude, and `research/jtbd-outcome-driven-innovation.md:186-192` — emotional/social jobs are a distinct category, not to be conflated with the core functional job. "Feel confident" is emotional; "comparing options" is too narrow (a step, not the job).
- Change: rewrite the good example to a core functional job at correct altitude, e.g. "When I'm about to make a significant purchase, I want to choose the option that best fits my needs, so I don't regret the decision later." Update the "Why it's good" prose on line 59 to note it names the core functional job (not an emotional job) and stays above any particular step in the shopping process.
- Effort: trivial
- Depends on: none

FC-06 — No distinction between job executor, buyer, and lifecycle roles; the first-person "I want" implicitly assumes one person

- Evidence: `intent/references/elements.md:41-51` — job is defined as "a goal someone is trying to achieve"; the pattern uses "I want to". `research/jtbd-outcome-driven-innovation.md:176-197` — ODI's needs framework splits **job executor**, **product lifecycle support team**, and **buyer**, each with distinct outcomes; in B2B they are different people.
- Change: rewrite `elements.md:41` to state: "A job is what the **job executor** — the person using the solution to make the progress — is trying to achieve. Buyer and lifecycle roles have distinct jobs that are out of scope for this framework unless separately captured." Retain the first-person pattern only for the executor's job.
- Effort: trivial
- Depends on: none

FC-07 — No job map / process steps and no core/related/emotional/consumption/purchase job taxonomy; declare the scope

- Evidence: entire `intent/references/elements.md` treats Job as a flat single element; `research/jtbd-outcome-driven-innovation.md:107-126` — ODI centers a `Define/Locate/Prepare/Confirm/Execute/Monitor/Modify/Conclude` job map used to organize outcomes; `research/jtbd-outcome-driven-innovation.md:184-197` — five job categories.
- Change: rewrite the Job section preamble (`elements.md:41-43`) to add: "Scope: this framework captures **core functional jobs** only. It does not model job steps, related jobs, emotional/social jobs, consumption-chain jobs, or purchase-decision jobs; outcomes that belong to those categories are either restated as core-job outcomes or omitted." No new sections — just declare the boundary.
- Effort: trivial
- Depends on: FC-06

FC-08 — No importance / satisfaction / opportunity prioritization mechanism; declare the omission

- Evidence: `intent/references/elements.md:67-97` — outcomes are defined without any priority or unmet-need signal; `research/jtbd-outcome-driven-innovation.md:219-241` — ODI's core mechanism is `opportunity = importance − satisfaction`, and needs-first sequencing rests on it.
- Change: add one line to the Outcome section (`elements.md:69`): "Prioritization by importance/satisfaction (ODI's opportunity score) is out of scope for this framework; capture outcomes as targets, not ranked opportunities. Rank them elsewhere if you need to." Keeps the omission deliberate.
- Effort: trivial
- Depends on: none

FC-09 — Outcome stability across solutions is not asserted, though the skill claims it for jobs

- Evidence: `intent/references/elements.md:41` — "A job… is functional, situational, and stable across changes to tools and features." No parallel claim for outcomes; `research/jtbd-outcome-driven-innovation.md:140-146` — ODI emphasizes outcomes are "stable over time (the job's success metrics persist; solutions change)."
- Change: append to `elements.md:69` (Outcome definition sentence): "Like the job it serves, an outcome should be stable across changes in the solutions used to get the job done."
- Effort: trivial
- Depends on: none

FC-10 — Requirements are declared "solution-aware" without stating the ODI deviation; force outcome-traceability to compensate

- Evidence: `intent/references/elements.md:149` — good requirements are "solution-aware but not implementation-specific"; `research/jtbd-outcome-driven-innovation.md:140-146,267-273` — ODI insists needs be "devoid of solutions." The skill's Requirement is the engineering-facing bridge, but the "solution-aware" phrase invites premature solution-shaping without linking back to a specific underserved outcome.
- Change: rewrite `elements.md:130-131` to: "A requirement is what the product must do or respect to mitigate one or more risks. It is the engineering-facing bridge from outcomes to architecture, so unlike outcomes it may name product capabilities — but it must name the risk(s) and, through them, the outcome(s) it serves. A requirement that cannot be traced to a named risk is a feature request, not a requirement."
- Effort: small
- Depends on: FC-03

FC-11 — Requirement template has no backlink to the outcomes/risks it serves — traceability is one-directional

- Evidence: `intent/assets/templates/requirement-document.md:1-30` — sections are Detail, Edge Cases, Examples, Acceptance Criteria, Dependencies (on other requirements only), See Also (records only). No section names the risks or outcomes this requirement mitigates. The reverse-direction map lives only on the outcome (`intent/assets/templates/outcome-document.md:15-18`, `outcome-directory.md:22-25`).
- Change: add a `## Mitigates` section to both `requirement-document.md` and `requirement-directory.md` immediately after the title/pattern, listing the risks (with logical IDs) and, through them, the outcomes served. Example: `- **O<NNN>-RSK<NNN>** - <Risk Name>  (serves **O<NNN>** - <Outcome Name>)`.
- Effort: trivial
- Depends on: FC-10

FC-12 — "Cover the parent" alignment check understates ODI's completeness demand; reword to enumerate customer-defined success metrics

- Evidence: `intent/references/elements.md:20` — "Cover the parent: do the siblings together cover the parent's scope without gaps?"; `research/jtbd-outcome-driven-innovation.md:148-152` — ODI expects ~50–150 outcomes per core job, framed as a **complete** enumeration of customer-defined success metrics, not a "no gaps" heuristic.
- Change: rewrite the Cover-the-parent bullet on `elements.md:20` to: "Cover the parent: do the siblings together *completely enumerate* what better looks like for the parent — every customer-defined way of judging success? 'No' means part of the parent's scope has no measurable expression." Apply the sharper phrasing wherever the coverage check is invoked (Outcome:71, Requirement:132).
- Effort: trivial
- Depends on: none

FC-13 — SKILL.md never declares its JTBD lineage or deliberate deviations; readers cannot tell which school they're inheriting

- Evidence: `intent/SKILL.md:1-45` — global guidance uses ODI-flavored phrasing ("solution-free", "measurable", "from the customer's perspective") but never names JTBD, ODI, Ulwick, or Christensen; no statement of what is deliberately excluded (job map, job types, opportunity scoring, executor/buyer split).
- Change: add a one-line "Lineage" bullet to `intent/SKILL.md` Global Guidance section: "Product-side vocabulary follows Ulwick/ODI (jobs, outcomes, needs-first sequencing); job syntax follows the Christensen job-story form. Deviations from ODI — Risk as a first-class element, no job map, no importance/satisfaction scoring, no job-type taxonomy, executor-only scope — are deliberate and covered in `references/elements.md`."
- Effort: trivial
- Depends on: FC-02, FC-03, FC-06, FC-07, FC-08
