---
title: "Writing statements — the right shape per field"
summary: "Each prose field has a fixed shape: story (JTBD), statement (declarative), responsibility (one-line charter), summary (elevator pitch)."
---
# Writing statements

**The question.** What shape should this element's prose take? The field name
tells you — it's deliberate semantic typing, not decoration.

**The tests, by field.**

- **`story`** (jobs) — a jobs-to-be-done narrative. Shape:
  *"When \<situation\>, I want \<motivation\> so I can \<expected outcome\>."*
  It describes a situation and goal, never a solution.
- **`statement`** (outcomes, risks, requirements) — one declarative sentence.
  - outcome: names a measurable improvement ("Minimize the time to ...").
  - risk: names a failure that could happen ("A reference points at ...").
  - requirement: says what must be true ("Every declared reference must
    resolve.").
- **`responsibility`** (components) — a one-line charter: the single thing this
  component is on the hook for. If it needs "and", consider splitting.
- **`summary`** (roots, decision records) — an elevator pitch: the core in one
  or two plain sentences.

**Good.** "When I read an element, I want to see why it exists so I can trust
it." (a story — situation, motivation, outcome)

**Bad.** "Add a rationale panel to the element view." (a solution masquerading
as a story)

**Failure mode.** Prose that fights its field — solutions in `story`, paragraphs
in `responsibility`, features in a risk `statement`. The shape is the author's
prompt toward the right altitude; ignore it and the element lands at the wrong
one.
