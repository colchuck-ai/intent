# Lens B — Description triggering

Summary

- Current description leads with what the skill *does* ("Build and maintain…"), not with when to *use* it — violates the baseline's imperative-phrasing rule.
- Trigger list ("creating, editing, tracing, or placing") is abstract; it misses paraphrased asks like "capture a decision," "write a product spec," "log why we chose X."
- Record-type keywords (CR/PDR/ADR) are present but only after a long list of elements; won't fire quickly on "add an ADR" style asks.
- Trailing "IDs/paths/tiers" is internal vocabulary users don't say — dead weight for matching.
- No negative boundary — risks over-triggering on any `docs/` work, README edits, or generic spec writing.

Findings

FB-01 — Description leads with "Build and maintain…" instead of "Use when…"

- Evidence: `intent/SKILL.md:3` — "Build and maintain a traceable docs/product and docs/engineering documentation tree…"
- Change: rewrite so the first clause is imperative "Use when…". The baseline (`research/agentskills-optimizing-skill-descriptions.md:26`) says: "Frame the description as an instruction to the agent: 'Use this skill when...' rather than 'This skill does...'". The current lead is descriptive, weakening trigger weight versus skills whose descriptions start with a task match.
- Effort: trivial
- Depends on: FB-06 (adopt one of the alternatives)

FB-02 — Trigger verbs are abstract; missing paraphrased asks

- Evidence: `intent/SKILL.md:3` — "Use when creating, editing, tracing, or placing these documents and their IDs/paths/tiers." Cross-check `intent/references/elements.md:39-51,67-97,98-127,128-157,186-213,223-287` and `intent/references/workflows.md:19-37` which name concrete user-facing acts (define the job, capture a decision, log a change, run alignment check, draft/judge/review).
- Change: expand the trigger list to include paraphrases that don't name the framework: `"write/draft a product spec"`, `"capture why we chose X"`, `"log this change"`, `"record this decision"`, `"define the customer's job"`, `"trace a requirement to a component"`, `"review docs for coherence"`, `"add an ADR/PDR/CR"`. Baseline (`research/agentskills-optimizing-skill-descriptions.md:28`) explicitly pushes for coverage "even if they don't explicitly mention 'CSV' or 'analysis.'"
- Effort: trivial
- Depends on: FB-06

FB-03 — "IDs/paths/tiers" is internal jargon users don't say

- Evidence: `intent/SKILL.md:3` — "…and their IDs/paths/tiers."
- Change: cut the phrase or replace with the paraphrase users actually say ("naming, placing, or picking a folder/slug for a doc"). Users don't ask for "tiers"; they ask "where does this file go?" Keeping raw framework vocab wastes trigger surface.
- Effort: trivial
- Depends on: FB-06

FB-04 — No negative boundary; risks over-triggering on any `docs/` work

- Evidence: `intent/SKILL.md:3` — no "skip for…" clause; whole scope is opt-in.
- Change: add one short negative clause: `"Skip for code changes, READMEs, marketing copy, or ad-hoc notes that don't link into the intent tree."` Baseline (`research/agentskills-optimizing-skill-descriptions.md:56-67`) treats near-misses (shared keywords, different intent) as the highest-value trigger discriminators; without a boundary clause the skill will fire on "help me write documentation."
- Effort: trivial
- Depends on: FB-06

FB-05 — Structural paths (`docs/product/`, `docs/engineering/`) buried mid-sentence

- Evidence: `intent/SKILL.md:3` — path fragments appear inside the "Build and maintain…" clause.
- Change: keep the two paths (they're strong scope signals that cheaply disambiguate from adjacent skills) but promote them near the front so trigger scans hit them quickly, e.g. "…intent documents under `docs/product/` or `docs/engineering/`…". Do not drop them.
- Effort: trivial
- Depends on: FB-06

FB-06 — Replace the description with one of three alternatives

- Evidence: `intent/SKILL.md:3` — current 377-char description exhibits FB-01…FB-05.
- Change: adopt one of the three alternatives below. Each is ≤ 1024 chars, imperative-first, keyword-dense on record types (CR/PDR/ADR) and element types (jobs, outcomes, risks, requirements, architecture, components), includes paraphrase triggers, and ends with a negative boundary.

  **Alt 1 — Pushy, paraphrase-heavy (~790 chars). Recommended default.**

  > Use when creating, editing, tracing, placing, or reviewing product or engineering intent documents under `docs/product/` or `docs/engineering/` — jobs-to-be-done, outcomes, risks, requirements, architecture, components, and their Change Records (CR), Product Decision Records (PDR), and Architectural Decision Records (ADR). Fire even when the user doesn't name the framework: "write a product spec", "capture why we chose X", "log this change", "record this decision", "draft an ADR/PDR/CR", "define the customer's job", "trace a requirement to a component", "check these docs for coherence", "where does this doc go / what should I name it". Also use for the draft/judge/review workflow and alignment checks. Skip for code changes, READMEs, marketing copy, or ad-hoc notes that don't link into the intent tree.

  Rationale: leads imperative, keeps the two `docs/` scope signals, spells out all seven elements + three record acronyms for keyword recall, lists six paraphrases from real user vocabulary, references workflow terms so process asks trigger too, and closes with a negative clause. Slightly long but well under 1024 and pushes hardest on Lens B goals.

  **Alt 2 — Terse (~500 chars). Use if length is a concern across many skills.**

  > Use when producing or maintaining traceable intent documents under `docs/product/` or `docs/engineering/` — jobs, outcomes, risks, requirements, architecture, components, and their Change Records (CR), Product Decision Records (PDR), and Architectural Decision Records (ADR). Also triggers on paraphrases: "write a product spec", "capture this decision", "log why we changed X", "draft an ADR", "trace this requirement", "define the customer's job". Skip for code, READMEs, or standalone notes.

  Rationale: shortest option that still meets all baseline principles. Trades exhaustive paraphrase coverage for compactness; likely close to Alt 1 on trigger rate but with a smaller context footprint.

  **Alt 3 — Structured, explicit categories (~940 chars). Use if evals show category-level misses.**

  > Use when working with the intent framework's traceable documentation tree in `docs/product/` or `docs/engineering/`. Covers seven element types — product, jobs (JTBD), outcomes, risks, requirements, architecture, components — and three record types — Change Records (CR), Product Decision Records (PDR), Architectural Decision Records (ADR). Trigger on explicit asks ("add a requirement", "draft an ADR", "file a CR") and paraphrased ones ("capture why we chose Postgres", "write down the customer's job", "log this spec change", "trace outcome to component", "review these docs for coherence", "where does this doc live / how do I name it"). Also trigger when following the draft/judge/review workflow or running an alignment check across the trace. Skip for code changes, product marketing, standalone READMEs, or notes that aren't part of the intent tree.

  Rationale: makes the taxonomy legible so the agent can match on category rather than exact keyword; slightly redundant with Alt 1 but useful if evals show the agent fails to connect a paraphrase to the correct element type.

- Effort: small (single-line frontmatter change)
- Depends on: none
